package http

import (
	"errors"
	"net/http"
	"time"

	"github.com/Ilyusha206420/first-day-of-internship/backend/internal/services"
	"github.com/gin-gonic/gin"
)

const (
	accessCookieName  = "access_token"
	refreshCookieName = "refresh_token"
)

type AuthHandler struct {
	authService    *services.AuthService
	sessionService *services.SessionService
	accessTTL      time.Duration
	refreshTTL     time.Duration
	secureCookie   bool
}

func NewAuthHandler(
	authService *services.AuthService,
	sessionService *services.SessionService,
	accessTTL time.Duration,
	refreshTTL time.Duration,
	secureCookie bool,
) *AuthHandler {
	return &AuthHandler{
		authService:    authService,
		sessionService: sessionService,
		accessTTL:      accessTTL,
		refreshTTL:     refreshTTL,
		secureCookie:   secureCookie,
	}
}

type loginRequest struct {
	Username string `json:"username" binding:"requiered"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var r loginRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Ivalid request body",
		})
		return
	}

	ctx := c.Request.Context()

	user, err := h.authService.AuthWithLoginPassword(ctx, r.Username, r.Password)
	if err != nil {
		if errors.Is(err, services.InvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid username or password",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unexcepted Internal error",
		})
		return
	}

	tokens, err := h.sessionService.CreateSession(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unexcepted internal error",
		})
		return
	}

	h.setTokenCookies(c, tokens.AccessToken, tokens.RefreshToken)
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (h *AuthHandler) setTokenCookies(c *gin.Context, accessToken string, refreshToken string) {
	h.setAccessCookie(c, accessToken)

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(refreshCookieName, refreshToken, int(h.refreshTTL.Seconds()), "/", "", h.secureCookie, true)

}

func (h *AuthHandler) setAccessCookie(c *gin.Context, accessToken string) {
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(accessCookieName, accessToken, int(h.accessTTL.Seconds()), "/", "", h.secureCookie, true)
}
