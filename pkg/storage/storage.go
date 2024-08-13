package storage

import (
	"AuthService/pkg/models"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"time"
)

const defaultEmail = "test@test.com"

var UserDoesNotExist = errors.New("User does not exist")

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
