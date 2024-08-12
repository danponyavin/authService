package handler

import (
	"AuthService/pkg/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

type Error struct {
	Message string `json:"message"`
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	base := router.Group("/api/v1")
	base.POST("auth/tokens", h.Auth)
	base.POST("auth/refresh", h.RefreshTokens)

	return router
}

type AuthRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
	Email  string    `json:"email" binding:"required,email"`
}

func (h *Handler) Auth(c *gin.Context) {
	var authRequest AuthRequest
	if err := c.ShouldBind(&authRequest); err != nil {
		c.JSON(http.StatusBadRequest, Error{Message: err.Error()})
	}
	fmt.Println(authRequest.UserID, authRequest.Email)
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) RefreshTokens(c *gin.Context) {
	var refreshTokenRequest RefreshTokenRequest
	if err := c.ShouldBind(&refreshTokenRequest); err != nil {
		c.JSON(http.StatusBadRequest, Error{Message: err.Error()})
	}
	fmt.Println(refreshTokenRequest.RefreshToken)
}
