package model

import (
	"backend/graph/model" // Import generated models
)

type InternalUser struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

func (u *InternalUser) ToGraphQL() *model.User {
	return &model.User{
		ID:       u.ID,
		Username: u.Name,
		Email:    u.Email,
	}
}
