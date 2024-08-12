package service

type UserStorage interface {
}

type UserService struct {
	storage UserStorage
}

func NewUserService(userStorage UserStorage) UserService {
	return UserService{userStorage}
}
