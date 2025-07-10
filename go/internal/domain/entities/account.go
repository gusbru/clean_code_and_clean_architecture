package entities

import "github.com/google/uuid"

type Account struct {
	Account_id uuid.UUID `json:"account_id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Document   string    `json:"document"`
	Password   string    `json:"password"`
}
