package handler

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mar-cial/space-auth/internal/core/domain"
	"github.com/mar-cial/space-auth/internal/core/port"
)

var ErrInternalServer = errors.New("Internal server error")

type authHandler struct {
	authService port.AuthService
}

func (a *authHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()

	var creds domain.Credentials
	if err := c.ShouldBind(&creds); err != nil {
		log.Println(err)
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": ErrInternalServer})
		return
	}

	user, err := a.authService.CreateUser(ctx, creds)
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": ErrInternalServer})
		return
	}

	session, err := a.authService.CreateSession(ctx, user.ID)
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": ErrInternalServer})
		return
	}

	// dev only. I know this is bad.
	// Store the session token as a cookie
	c.SetCookie(
		"session_id",
		session.Token,
		int(session.ExpiresAt.Sub(time.Now()).Seconds()),
		"/",
		"localhost",
		false,
		true,
	)

	c.HTML(http.StatusOK, "user_registered.html", gin.H{"user": ErrInternalServer})
}

func (a *authHandler) Login(c *gin.Context) {
	// credentials
	var creds domain.Credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		log.Println(err)
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": ErrInternalServer})
		return
	}

	session, err := a.authService.CreateSession(c.Request.Context(), creds.Phonenumber)
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Unable to sign in"})
		return
	}

	// Store the session token as a cookie
	c.SetCookie(
		"session_id",  // Cookie name
		session.Token, // Cookie value (session token)
		int(session.ExpiresAt.Sub(time.Now()).Seconds()), // Max age in seconds
		"/",         // Path
		"localhost", // Domain (adjust for your environment)
		false,       // Secure (set true for HTTPS)
		true,        // HttpOnly (prevent JavaScript access)
	)

	// Render an HTML snippet for HTMX
	c.HTML(http.StatusOK, "login_success.html", gin.H{"message": "welcome!"})
}

func (a *authHandler) Logout(c *gin.Context) {
	// Retrieve the session token from the cookie
	sessionToken, err := c.Cookie("session_id")
	if err != nil {
		log.Println("Missing session cookie:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Session not found"})
		return
	}

	// Delete the session from the storage
	if err := a.authService.DeleteSession(c.Request.Context(), sessionToken); err != nil {
		log.Println("Error deleting session:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log out"})
		return
	}

	// Clear the session cookie
	c.SetCookie(
		"session_id", // Cookie name
		"",           // Clear the cookie value
		-1,           // Expire the cookie immediately
		"/",          // Path
		"localhost",  // Domain (adjust for production)
		false,        // Secure (set true for HTTPS in production)
		true,         // HttpOnly
	)

	c.HTML(http.StatusOK, "logout_success.html", gin.H{"message": "Logged out successfully"})
}

func NewAuthHandler(srv port.AuthService) port.AuthHandler {
	return &authHandler{authService: srv}
}
