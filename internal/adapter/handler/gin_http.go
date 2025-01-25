package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mar-cial/space-auth/internal/core/port"
)

type authHandler struct {
	authService port.AuthService
}

func (a *authHandler) Register(c *gin.Context) {
	var req struct {
		Phonenumber string `json:"phonenumber" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := a.authService.Register(c.Request.Context(), req.Phonenumber)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (a *authHandler) Login(c *gin.Context) {
	var req struct {
		Phonenumber string `json:"phonenumber" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := a.authService.Login(c.Request.Context(), req.Phonenumber)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (a *authHandler) Logout(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := a.authService.Logout(c.Request.Context(), req.Token)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func NewAuthHandler(as port.AuthService) port.AuthHandler {
	return &authHandler{authService: as}
}
