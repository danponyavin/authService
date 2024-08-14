package handler

import (
	"AuthService/docs"
	"AuthService/pkg/models"
	"AuthService/pkg/service"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

const BasePath = "/api/v1/"

// @BasePath /api/v1
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	docs.SwaggerInfo.BasePath = BasePath
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	base := router.Group("/api/v1")
	base.GET("auth/tokens/:user_id", h.GetTokensHandler)
	base.POST("auth/refresh", h.RefreshTokens)

	return router
}

// GetTokens godoc
// @Summary Получение токенов
// @Schemes
// @Description Получение пары токенов по userID
// @Accept json
// @Produce json
// @Param user_id path string true "User ID" Example(5fd3b119-408e-451e-8bd3-641b38fa8cde)
// @Success 200 {object} models.Tokens
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /auth/tokens/{user_id} [get]
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

// RefreshTokens godoc
// @Summary Обновление токенов
// @Schemes
// @Description Обновление Access и Refresh токенов по Refresh токену
// @Accept json
// @Produce json
// @Param data body models.RefreshTokenRequest true "Входные параметры"
// @Success 200 {object} models.Tokens
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /auth/refresh [post]
func (h *Handler) RefreshTokens(c *gin.Context) {
	var refreshTokenRequest models.RefreshTokenRequest
	if err := c.ShouldBind(&refreshTokenRequest); err != nil {
		c.JSON(http.StatusBadRequest, Error{Message: err.Error()})
		return
	}

	ip := c.ClientIP()

	response, err := h.services.UserService.RefreshTokens(refreshTokenRequest.RefreshToken, ip)
	if err != nil {
		switch {
		case errors.Is(err, service.InvalidRefreshTokenError), errors.Is(err, service.TokenExpiredError):
			c.JSON(http.StatusBadRequest, Error{Message: err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, response)
}
