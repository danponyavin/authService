package storage

import (
	"AuthService/pkg/models"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"time"
)

const defaultEmail = "test@test.com"

var (
	UserDoesNotExist = errors.New("User does not exist")
)

type UserStorage interface {
	CreateUser(userID uuid.UUID) error
	GetUserData(userID uuid.UUID) (models.UserData, error)
	CreateSession(auth models.AuthModel) error
	GetSession(refreshTokenHash string) (models.Session, error)
	DeleteSession(refreshTokenHash string) error
}

func (p *PostgresStorage) CreateUser(userID uuid.UUID) error {

	//Заполню моковыми данными email пользователя
	_, err := p.db.Exec("INSERT INTO users (id, email) VALUES ($1, $2)", userID.String(), defaultEmail)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresStorage) GetUserData(userID uuid.UUID) (models.UserData, error) {
	var user models.UserData
	err := p.db.QueryRow("SELECT * FROM users WHERE id = $1", userID).Scan(&user.UserID, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.UserData{}, UserDoesNotExist
		}
		return models.UserData{}, err
	}

	return user, nil
}

func (p *PostgresStorage) CreateSession(auth models.AuthModel) error {

	_, err := p.db.Exec("INSERT INTO refresh_sessions (user_id, ip, refresh_token, issued_at, expires_in) "+
		"VALUES ($1, $2, $3, $4, $5)",
		auth.UserID.String(), auth.ClientIP, auth.RefreshTokenHash, time.Now(), time.Now().Add(auth.RefreshTokenTTL))
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresStorage) GetSession(refreshTokenHash string) (models.Session, error) {

	var s models.Session
	err := p.db.QueryRow("SELECT user_id, refresh_token, ip, issued_at, expires_in FROM refresh_sessions "+
		"WHERE refresh_token = $1", refreshTokenHash).Scan(&s.UserID, &s.RefreshTokenHash, &s.ClientIP,
		&s.IssuedAt, &s.ExpiresIn)
	if err != nil {
		return s, err
	}

	return s, nil
}

func (p *PostgresStorage) DeleteSession(refreshTokenHash string) error {

	_, err := p.db.Exec("DELETE FROM refresh_sessions WHERE refresh_token = $1", refreshTokenHash)
	if err != nil {
		return err
	}

	return nil
}
