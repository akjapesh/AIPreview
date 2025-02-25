package resolver

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// This file will not be regenerated automatically.
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Conn       *pgxpool.Pool
	RedisCache *redis.Client
}

// NewResolver initializes the resolver with a database connection
func NewResolver(conn *pgxpool.Pool, redisClient *redis.Client) *Resolver {
	return &Resolver{
		Conn:       conn,
		RedisCache: redisClient,
	}
}
