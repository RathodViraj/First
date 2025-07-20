package service

import (
	"First/model"
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserSrv *UserService
}

func NewAuthService(userService *UserService) *AuthService {
	return &AuthService{userService}
}

func (s *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *AuthService) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *AuthService) GenerateToken(user *model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Id,
		"exp": time.Now().Add(time.Minute * 15).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func (s *AuthService) ValidateToken(tokenStr string) (int, string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method.")
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return -1, "", nil
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id := int(claims["sub"].(float64))
		role := claims["role"].(string)
		return id, role, nil
	}

	return -1, "", errors.New("Invalid or expired token.")
}

func (s *AuthService) Authenticate(login *model.LoginRequest) (*model.User, error) {
	// Log the email being searched
	log.Printf("Attempting login for email: %s", login.Email)

	user, err := s.UserSrv.GetUserByEmail(login.Email)
	if err != nil {
		log.Printf("Authentication failed for %s: %v", login.Email, err)
		return nil, errors.New("invalid credentials")
	}

	log.Printf("User found: ID %d, checking password...", user.Id)

	if !s.CheckPassword(login.Password, user.Password) {
		log.Printf("Password mismatch for user %d", user.Id)
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
