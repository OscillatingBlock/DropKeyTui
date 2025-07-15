package api

import "time"

type ErrorResponse struct {
	Message string `json:"message"`
}

type RegisterUserRequest struct {
	PublicKey string `json:"public_key"`
}

type RegisterUserResponse struct {
	ID string `json:"id"`
}

type User struct {
	ID        string `json:"id"`
	PublicKey string `json:"public_key"`
}

type AuthRequest struct {
	ID        string `json:"id"`
	Signature string `json:"signature"`
	Challenge string `json:"challenge"`
	PublicKey string `json:"public_key"`
}

type AuthResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type Paste struct {
	ID         string    `json:"ID"`
	Ciphertext string    `json:"ciphertext"`
	Signature  string    `json:"signature"`
	PublicKey  string    `json:"public_key"`
	ExpiresAt  time.Time `json:"expires_in"`
}

type PasteRequest struct {
	Ciphertext string `json:"ciphertext"`
	Signature  string `json:"signature"`
	PublicKey  string `json:"public_key"`
	ExpiresIn  int    `json:"expires_in"`
}

type CreatePasteResponse struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}
