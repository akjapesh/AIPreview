package resolver

import (
	"github.com/jackc/pgx/v5"
	"backend/pkg/mlServices"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// This file will not be regenerated automatically.
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Conn *pgx.Conn
	MLClient *mlServices.MLClient
}

// NewResolver initializes the resolver with a database connection and ml client setup
func NewResolver(conn *pgx.Conn,mlClient *mlServices.MLClient) *Resolver {
	return &Resolver{
		Conn: conn,
		MLClient: mlClient,
	}
}
