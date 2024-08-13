package handler

import (
	"AuthService/pkg/models"
	"AuthService/pkg/service"
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
	base.GET("auth/tokens/:user_id", h.GetTokensHandler)
	base.POST("auth/refresh", h.RefreshTokens)

	return router
}

func (h *Handler) GetTokensHandler(c *gin.Context) {
	var userID uuid.UUID
	userIDParam := c.Param("user_id")
	if userIDParam != "" {
		parsedUserID, err := uuid.Parse(userIDParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, Error{Message: "Invalid user_id"})
			return
		}
		userID = parsedUserID
	} else {
		c.JSON(http.StatusBadRequest, Error{Message: "Empty user_id"})
		return
	}

	authModel := models.AuthModel{
		UserID:   userID,
		ClientIP: c.ClientIP(),
	}

	response, err := h.services.UserService.GetTokens(authModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Error{Message: "Service error"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) RefreshTokens(c *gin.Context) {
	var refreshTokenRequest models.RefreshTokenRequest
	if err := c.ShouldBind(&refreshTokenRequest); err != nil {
		c.JSON(http.StatusBadRequest, Error{Message: err.Error()})
		return
	}
}
