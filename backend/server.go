package main

import (
	"backend/graph"
	"log"
	"net/http"
	"os"
	"fmt"
	"context"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
	"backend/graph/resolver"
    "backend/pkg/auth"
    "github.com/go-chi/chi"
    "github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
)

const defaultPort = "8080"

func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		if r.Header.Get("X-Auth-Optional") == "true" {
			next.ServeHTTP(w, r)
			return
		}

        cookie,err := r.Cookie("auth_token")

		if err != nil {
			// If no auth_token, just pass request without auth
			fmt.Println("🔴 Missing Authorization Token")
            http.Error(w, "Missing Authorization Token", http.StatusUnauthorized)
            return
		}
        

        tokenStr := cookie.Value
        claims, err := auth.ValidateToken(tokenStr)
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        // Add user ID to context
        ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()

	// ✅ CORS middleware (Allows frontend at http://localhost:3000)
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Change this in production
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true, // ✅ Allows sending cookies (JWT in cookies)
	})

    router.Use(middleware.Logger)
    router.Use(middleware.Recoverer)
	router.Use(corsMiddleware.Handler)
	router.Use(authMiddleware)   

	resolver := resolver.NewResolver()

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
    router.Handle("/graphql", srv)

    log.Printf("connect to http://localhost:8080/ for GraphQL playground")
    log.Fatal(http.ListenAndServe(":8080", router))

}
