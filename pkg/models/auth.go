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

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type Session struct {
	UserID           uuid.UUID `db:"user_id"`
	ClientIP         string    `db:"ip"`
	RefreshTokenHash string    `db:"refresh_token"`
	IssuedAt         time.Time `db:"issued_at"`
	ExpiresIn        time.Time `db:"expires_in"`
}
