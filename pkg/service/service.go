package service

import "AuthService/pkg/storage"

type Service struct {
	UserService
}

func NewService(storage *storage.PostgresStorage) *Service {
	return &Service{
		UserService: NewUserService(storage),
	}
}
