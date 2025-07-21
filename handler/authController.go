package handler

import (
	"First/model"
	"First/service"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var loginReq model.LoginRequest

	if err := ctx.BindJSON(&loginReq); err != nil {
		log.Printf("JSON decode error: %v", err)
		JSONError(ctx, http.StatusBadRequest, "Invalid request body")
		return
	}

	loginReq.Email = strings.TrimSpace(strings.ToLower(loginReq.Email))
	loginReq.Password = strings.TrimSpace(loginReq.Password)

	if loginReq.Email == "" || loginReq.Password == "" {
		JSONError(ctx, http.StatusBadRequest, "Email and password cannot be empty")
		return
	}

	user, err := h.authService.Authenticate(&loginReq)
	if err != nil {
		log.Printf("Login failed for %s: %v", loginReq.Email, err)
		JSONError(ctx, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	token, err := h.authService.GenerateToken(user)
	if err != nil {
		log.Printf("Token generation failed: %v", err)
		JSONError(ctx, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"token": token})
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	var user model.User

	if err := ctx.BindJSON(&user); err != nil {
		JSONError(ctx, http.StatusBadRequest, "Invalid request body")
		return
	}

	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	user.Name = strings.TrimSpace(user.Name)
	user.Password = strings.TrimSpace(user.Password)

	if user.Password == "" {
		JSONError(ctx, http.StatusBadRequest, "Password cannot be empty")
		return
	}

	hashedPassword, err := h.authService.HashPassword(user.Password)
	if err != nil {
		log.Printf("Password hashing error: %v", err)
		JSONError(ctx, http.StatusInternalServerError, "Failed to process password")
		return
	}
	user.Password = hashedPassword

	if err := h.authService.UserSrv.RegisterUser(&user); err != nil {
		log.Printf("Registration failed: %v", err)
		JSONError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	user.Password = ""
	ctx.IndentedJSON(http.StatusCreated, user)
}
