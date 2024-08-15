package service

import (
	"AuthService/pkg/models"
	"AuthService/pkg/storage"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"
	"gopkg.in/gomail.v2"
	"log"
	"time"
)

const (
	salt                     = "Hk23nS9dG9-;D"
	signingKey               = "s9sf,-fS7jdY6l"
	RefreshTTL time.Duration = time.Hour * 720
	AccessTTL  time.Duration = time.Hour * 1
)

var (
	TokenExpiredError        = errors.New("Token is expired")
	InvalidRefreshTokenError = errors.New("Invalid Refresh Token")
)

type AuthService struct {
	storage storage.UserStorage
}

func NewAuthService(userStorage storage.UserStorage) *AuthService {
	return &AuthService{userStorage}
}

func (a *AuthService) GetTokens(auth models.AuthModel) (models.Tokens, error) {
	resp, err := CreateTokens(auth.UserID, auth.ClientIP)
	if err != nil {
		log.Println("CreateTokens:", err)
		return resp, err
	}

	hashedRefreshToken := hashRefreshToken(resp.RefreshToken)

	auth.RefreshTokenHash = hashedRefreshToken
	auth.RefreshTokenTTL = RefreshTTL

	//Если пользователь еще не существует, то создаем его
	_, err = a.storage.GetUserData(auth.UserID)
	if err != nil {
		if errors.Is(err, storage.UserDoesNotExist) {
			err = a.storage.CreateUser(auth.UserID)
			if err != nil {
				log.Println("CreateUser:", err)
				return resp, err
			}
		} else {
			log.Println("GetUserData:", err)
			return resp, err
		}
	}

	err = a.storage.CreateSession(auth)
	if err != nil {
		log.Println("CreateSession:", err)
		return resp, err
	}

	return resp, nil
}

func createJWT(userID uuid.UUID, ip string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"ip":      ip,
		"exp":     time.Now().Add(AccessTTL).Unix(),
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

func hashRefreshToken(refreshToken string) string {
	combinedToken := refreshToken + salt
	tokenBytes := []byte(combinedToken)

	hash := sha3.New512()
	hash.Write(tokenBytes)
	hashedToken := hash.Sum(nil)

	hashedTokenHex := hex.EncodeToString(hashedToken)

	return hashedTokenHex
}

func CreateTokens(userID uuid.UUID, ip string) (models.Tokens, error) {
	tokens := models.Tokens{}
	accessToken, err := createJWT(userID, ip)
	if err != nil {
		log.Println("createJWT:", err)
		return tokens, err
	}
	tokens.AccessToken = accessToken

	refreshToken, err := generateRefreshToken()
	if err != nil {
		log.Println("generateRefreshToken:", err)
		return tokens, err
	}
	tokens.RefreshToken = refreshToken

	return tokens, nil
}

func (a *AuthService) RefreshTokens(inp models.RefreshModel) (models.Tokens, error) {
	var resp models.Tokens
	oldRefreshTokenHash := hashRefreshToken(inp.RefreshToken)

	oldSession, err := a.storage.GetSession(oldRefreshTokenHash)
	if err != nil {
		log.Println("GetSession:", InvalidRefreshTokenError)
		return resp, InvalidRefreshTokenError
	}

	if oldSession.ExpiresIn.Before(time.Now()) {
		return resp, TokenExpiredError
	}

	if oldSession.ClientIP != inp.ClientIP {
		userData, err := a.storage.GetUserData(oldSession.UserID)
		if err != nil {
			log.Println("GetUserData:", err)
			return resp, err
		}

		m := gomail.NewMessage()
		m.SetHeader("From", "auth_service@test.ru")
		m.SetHeader("To", userData.Email)
		m.SetHeader("Subject", "Подозрительный вход")

		warningMessage := fmt.Sprintf("На Ваш аккаунт был осуществлен вход с другого ip адреса: %s", inp.ClientIP)
		log.Println(warningMessage)

		m.SetBody("text/html", fmt.Sprintf("<p>%s</p>", warningMessage))

		d := gomail.NewDialer("smtp.mail.ru", 465, "test", "test")
		if err = d.DialAndSend(m); err != nil {
			// Так как данные моковые, то не будем возвращать ошибку отправки письма
			err = nil
		}
	}

	resp, err = CreateTokens(oldSession.UserID, oldSession.ClientIP)
	if err != nil {
		log.Println("CreateTokens:", err)
		return resp, err
	}

	newHashedRefreshToken := hashRefreshToken(resp.RefreshToken)

	err = a.storage.DeleteSession(oldRefreshTokenHash)
	if err != nil {
		log.Println("DeleteSession:", err)
		return resp, err
	}

	err = a.storage.CreateSession(models.AuthModel{UserID: oldSession.UserID, ClientIP: oldSession.ClientIP,
		RefreshTokenHash: newHashedRefreshToken, RefreshTokenTTL: RefreshTTL})
	if err != nil {
		log.Println("CreateSession:", err)
		return resp, err
	}

	return resp, nil
}
