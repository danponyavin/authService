package service

import (
	"AuthService/pkg/models"
	"AuthService/pkg/storage"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

const (
	signingKey               = "s9sf,-fS7jdY6l"
	RefreshTTL time.Duration = time.Hour * 720
	AccessTTL  time.Duration = time.Hour * 1
)

type UserStorage interface {
	CreateUser(userID uuid.UUID) error
	GetUserData(userID uuid.UUID) (models.UserData, error)
	CreateSession(auth models.AuthModel) error
}

type UserService struct {
	storage UserStorage
}

func NewUserService(userStorage UserStorage) UserService {
	return UserService{userStorage}
}

func (u *UserService) GetTokens(auth models.AuthModel) (models.TokensResponse, error) {
	resp := models.TokensResponse{}
	accessToken, err := createJWT(auth.UserID, auth.ClientIP, AccessTTL)
	if err != nil {
		log.Println("createJWT:", err)
		return resp, err
	}

	refreshToken, err := generateRefreshToken()
	if err != nil {
		log.Println("generateRefreshToken:", err)
		return resp, err
	}

	hashedRefreshToken, err := hashRefreshToken(refreshToken)
	if err != nil {
		log.Println("hashRefreshToken:", err)
		return resp, err
	}
	auth.RefreshTokenHash = hashedRefreshToken
	auth.RefreshTokenTTL = RefreshTTL
	resp.AccessToken = accessToken
	resp.RefreshToken = hashedRefreshToken

	//Если пользователь еще не существует, то создаем его
	_, err = u.storage.GetUserData(auth.UserID)
	if err != nil {
		if errors.Is(err, storage.UserDoesNotExist) {
			err = u.storage.CreateUser(auth.UserID)
			if err != nil {
				log.Println("CreateUser:", err)
				return resp, err
			}
		} else {
			log.Println("GetUserData:", err)
			return resp, err
		}
	}

	err = u.storage.CreateSession(auth)
	if err != nil {
		log.Println("CreateSession:", err)
		return resp, err
	}

	return resp, nil
}

func createJWT(userID uuid.UUID, ip string, TTL time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"ip":      ip,
		"exp":     time.Now().Add(TTL).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenString, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

func hashRefreshToken(refreshToken string) (string, error) {
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	hashedTokenString := string(hashedToken)

	return hashedTokenString, nil
}
