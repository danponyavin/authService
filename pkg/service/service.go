package service

import (
	"AuthService/pkg/models"
	"AuthService/pkg/storage"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	GetTokens(auth models.AuthModel) (models.Tokens, error)
	RefreshTokens(inp models.RefreshModel) (models.Tokens, error)
}

type Service struct {
	Authorization
}

func NewService(storage *storage.PostgresStorage) *Service {
	return &Service{
		Authorization: NewAuthService(storage),
	}
}
