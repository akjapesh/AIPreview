// services/ml_client.go
package mlServices

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"golang.org/x/time/rate"
	"context"
	"bufio"
	"strings"
)

type MLClient struct {
	BaseURL    string
	HTTPClient *http.Client
	rateLimiter *rate.Limiter
}

func NewMLClient(baseURL string,requestsPerSecond int,timeout time.Duration) *MLClient {
	return &MLClient{
		BaseURL:    baseURL,
		HTTPClient:  &http.Client{
			Timeout: timeout,
		},
		rateLimiter: rate.NewLimiter(rate.Every(time.Second/time.Duration(requestsPerSecond)), requestsPerSecond),
	}
}

type GenerationRequest struct {
	Prompt string `json:"prompt"`
}

type GenerationResponse struct {
	Code    string   `json:"code"`
	Context []string `json:"context"`
}

func (c *MLClient) GenerateComponent(ctx context.Context, prompt string) (<-chan string, error) {
    ch := make(chan string)
    
    go func() {
        defer close(ch)
        const maxRetries = 3
        var lastError error
        var retryCount int

        for retryCount = 0; retryCount < maxRetries; retryCount++ {
            select {
            case <-ctx.Done():
                return
            default:
                // Rate limiting
                if err := c.rateLimiter.Wait(ctx); err != nil {
                    ch <- fmt.Sprintf("Rate limiter error: %v", err)
                    return
                }

                // Create request
                reqBody := GenerationRequest{Prompt: prompt}
                jsonBody, err := json.Marshal(reqBody)
                if err != nil {
                    ch <- fmt.Sprintf("Marshal error: %v", err)
                    return
                }

                req, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/generate", bytes.NewBuffer(jsonBody))
                if err != nil {
                    ch <- fmt.Sprintf("Request creation error: %v", err)
                    return
                }
                req.Header.Set("Content-Type", "application/json")
                req.Header.Set("Accept", "text/event-stream")

                // Execute request
                resp, err := c.HTTPClient.Do(req)
                if err != nil {
                    lastError = fmt.Errorf("request failed: %w", err)
                    time.Sleep(time.Duration(retryCount) * 500 * time.Millisecond)
                    continue
                }

                // Handle response
                if resp.StatusCode != http.StatusOK {
                    body, _ := io.ReadAll(resp.Body)
                    resp.Body.Close()
                    lastError = fmt.Errorf("status %d: %s", resp.StatusCode, body)
                    time.Sleep(time.Duration(retryCount) * 500 * time.Millisecond)
                    continue
                }

                // Successful connection - monitor for data
                fmt.Println("Connected to ML service, waiting for stream...")
                dataReceived := false
                timeout := time.After(30 * time.Second)

                reader := bufio.NewReader(resp.Body)
                defer resp.Body.Close()

            readLoop:
                for {
                    select {
                    case <-timeout:
                        lastError = fmt.Errorf("no data received within 30 seconds")
                        break readLoop
                    case <-ctx.Done():
						fmt.Println("--------all out--------")
                        return
                    default:
                        line, err := reader.ReadString('\n')
                        if err != nil {
                            if err == io.EOF {
                                if !dataReceived {
                                    lastError = fmt.Errorf("connection closed without data")
                                }
                                break readLoop
                            }
                            lastError = fmt.Errorf("read error: %w", err)
                            break readLoop
                        }

                        line = strings.TrimSpace(line)
                        if line == "" {
                            continue // Skip empty lines
                        }

                        if !dataReceived {
                            dataReceived = true
                            fmt.Println("Started receiving data from stream")
                        }

                        if strings.HasPrefix(line, "data: ") {
                            content := strings.TrimPrefix(line, "data: ")
                            select {
                            case ch <- content:
                            case <-ctx.Done():
                                return
                            }
                        }
                    }
                }

                if dataReceived {
                    return // Successfully received data
                }

                // If we got here, retry
                fmt.Printf("Retrying (attempt %d/%d)", retryCount+1, maxRetries)
            }
        }

        if lastError != nil {
            ch <- fmt.Sprintf("Failed after %d attempts: %v", retryCount, lastError)
        }
    }()

    return ch, nil
}