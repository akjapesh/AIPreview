package resolver

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
    "backend/internal/model"
)

// This file will not be regenerated automatically.
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{
    users map[string]*model.InternalUser
}

func NewResolver() *Resolver {
    return &Resolver{
        users: make(map[string]*model.InternalUser),
    }
}
