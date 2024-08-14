package service

import "AuthService/pkg/storage"

type Service struct {
	Authorization
}

func NewService(storage *storage.PostgresStorage) *Service {
	return &Service{
		Authorization: NewAuthService(storage),
	}
}
