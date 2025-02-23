package main

import (
	"backend/database"
	"backend/graph"
	"backend/graph/resolver"
	"backend/pkg/auth"
	"context"
	"fmt"

	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jackc/pgx"
	"github.com/rs/cors"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

func authMiddleware(next http.Handler) http.Handler { //This middleware extracts the Authorization token from HTTP headers.
	fmt.Println("authMiddleware")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("X-Auth-Optional") == "true" {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("auth_token")

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

func server(conn *pgx.Conn) {
	database.LoadEnvFile(".env")
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

	// Pass the database connection to the resolver
	resolver := resolver.NewResolver(conn)
	//srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

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

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))

}
