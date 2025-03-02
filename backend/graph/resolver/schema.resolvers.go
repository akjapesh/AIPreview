package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.66

import (
	"backend/database"
	"backend/graph"
	"backend/graph/model"
	"backend/pkg/auth"
	"backend/pkg/mlServices"
	"context"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// SignUp is the resolver for the signUp field.
func (r *mutationResolver) SignUp(ctx context.Context, input model.SignUpInput) (bool, error) {
	exists, err := database.UserExists(r.Conn, input.Email, input.Name)
	if err != nil {
		return false, fmt.Errorf("error checking existing user: %w", err)
	}
	if exists {
		return false, fmt.Errorf("user with given email or username already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return false, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create the user struct
	newUser := database.User{
		Username:     strings.ToLower(input.Name),
		Email:        strings.ToLower(input.Email),
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		IsActive:     true,
	}

	// Save user to the database
	_, err = database.SaveUser(r.Conn, newUser)
	if err != nil {
		return false, fmt.Errorf("failed to create user: %w", err)
	}

	return true, nil
}

// SignIn is the resolver for the signIn field.
func (r *mutationResolver) SignIn(ctx context.Context, input model.SignInInput) (*model.AuthResponse, error) {
	// Find user
	user, err := database.FetchUserByEmail(r.Conn, strings.ToLower(input.Email))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	// Verify the password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	// Generate the token
	token, err := auth.GenerateToken(user.ID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Convert the internal user model to the GraphQL model
	graphQLUser := &model.User{
		ID:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
		IsActive: &user.IsActive,
		// CreatedAt: user.CreatedAt.Format(time.RFC3339),
		// UpdatedAt: user.CreatedAt.Format(time.RFC3339),
	}

	return &model.AuthResponse{
		Token: token,
		User:  graphQLUser,
	}, nil
}

// GetUserByEmailID is the resolver for the getUserByEmailId field.
func (r *queryResolver) GetUserByEmailID(ctx context.Context, emailID string) (*model.User, error) {
	user, err := database.FetchUserByEmail(r.Conn, strings.ToLower(emailID))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}
	createdAt := user.CreatedAt.Format(time.RFC3339)
	updatedAt := user.UpdatedAt.Format(time.RFC3339)
	return &model.User{
		ID:        user.ID.String(),
		Username:  user.Username,
		Email:     user.Email,
		IsActive:  &user.IsActive,
		CreatedAt: &createdAt,
		UpdatedAt: &updatedAt,
	}, nil
}

// GenerateResponse is the resolver for the generateResponse field.
func (r *subscriptionResolver) GenerateResponse(ctx context.Context, userID string, discussionID string, prompt string) (<-chan *model.GeneratedCode, error) {
	// Create output channel for streaming responses
	out := make(chan *model.GeneratedCode)

	// Call ML service to get streaming channel
	mlStream, err := r.MLClient.GenerateComponent(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to start generation: %w", err)
	}
	fmt.Println("mlStream==", mlStream)
	go func() {
		defer close(out)
		var fullResponse strings.Builder
		var finalFiles []*model.CodeFile

		// Stream chunks from ML service
		for chunk := range mlStream {
			fmt.Println("chunk==", chunk)
			select {
			case <-ctx.Done():
				// Handle client disconnection
				fmt.Println("-------last chunk---------")
				return
			default:
				fullResponse.WriteString(chunk)

				// Send incremental chunk to client
				out <- &model.GeneratedCode{
					Chunk:      &chunk,
					IsComplete: false,
				}
			}
		}

		// After stream completes, process final output
		sanitized := mlServices.SanitizeCode(fullResponse.String())
		finalFiles = []*model.CodeFile{
			{
				Name:    "App.tsx",
				Content: sanitized,
			},
		}

		// Send final completion message
		out <- &model.GeneratedCode{
			Files:      finalFiles,
			IsComplete: true,
		}
	}()

	return out, nil
}

// Mutation returns graph.MutationResolver implementation.
func (r *Resolver) Mutation() graph.MutationResolver { return &mutationResolver{r} }

// Query returns graph.QueryResolver implementation.
func (r *Resolver) Query() graph.QueryResolver { return &queryResolver{r} }

// Subscription returns graph.SubscriptionResolver implementation.
func (r *Resolver) Subscription() graph.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
