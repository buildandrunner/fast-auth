package handler

import (
	"bytes"
	"errors"
	"io"
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

	c.HTML(http.StatusOK, "user_registered.html", gin.H{"user": user})
}

func (a *authHandler) Login(c *gin.Context) {
	// Debug: Read raw request body
	body, _ := io.ReadAll(c.Request.Body)
	log.Println("Raw Request Body:", string(body))

	// Reset the request body so Gin can bind it
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	var creds domain.Credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		log.Println("JSON Unmarshal Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Validate user credentials
	valid, err := a.authService.ValidateUser(c.Request.Context(), creds)
	if err != nil || !valid {
		log.Println("Invalid login attempt:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Create session after successful validation
	session, err := a.authService.CreateSession(c.Request.Context(), creds.Phonenumber)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to sign in"})
		return
	}

	// Set session cookie
	c.SetCookie("session_id", session.Token, int(session.ExpiresAt.Sub(time.Now()).Seconds()), "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Welcome!"})
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
