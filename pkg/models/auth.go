package models

import (
	"github.com/google/uuid"
	"time"
)

type AuthModel struct {
	UserID           uuid.UUID `db:"user_id"`
	ClientIP         string    `db:"ip"`
	RefreshTokenHash string    `db:"refresh_token"`
	RefreshTokenTTL  time.Duration
}

type UserData struct {
	UserID uuid.UUID
	Email  string
}

type TokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}
