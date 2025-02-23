package resolver

import "github.com/jackc/pgx/v5"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// This file will not be regenerated automatically.
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Conn *pgx.Conn
}

// NewResolver initializes the resolver with a database connection
func NewResolver(conn *pgx.Conn) *Resolver {
	return &Resolver{
		Conn: conn,
	}
}
