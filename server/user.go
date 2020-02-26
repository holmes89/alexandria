package main

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

var (
	key = os.Getenv("JWT_SECRET")
	ErrInvalidLogin = errors.New("invalid login")
)


type User struct {
	ID string `firestore:"id"`
	Username string `firestore:"username"`
	Password string `firestore:"password"`
}

type Token struct {
	Type         string    `json:"type"`
	AccessToken  string    `json:"token"`
	Expires      time.Time `json:"expires"`
}

type UserService interface {
	Authenticate(ctx context.Context, username, password string) (*Token, error)
}

type UserRepository interface {
	FindUserByUsername(ctx context.Context, username string) (*User, error)
}

func NewUserFirestoreRepository(database *UserFirestoreDatabase) UserRepository {
	return database
}

type userService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	if key == "" {
		logrus.Fatal("jwt not set")
	}
	return &userService{
		repo: repo,
	}
}

func (s *userService) Authenticate(ctx context.Context, username, password string) (*Token, error) {
	user, err := s.repo.FindUserByUsername(ctx, username)
	if err != nil {
		logrus.WithError(err).Error("unable to find user")
		return nil, errors.New("unable to find user")
	}
	if user == nil {
		return nil, ErrInvalidLogin
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		logrus.WithField("username", username).WithError(err).Error("invalid login")
		return nil, ErrInvalidLogin
	}

	tokenExpiration := time.Now().Add(time.Hour*72)
	tokenString, err := s.getToken(user, tokenExpiration)

	if err != nil {
		logrus.WithField("username", username).WithError(err).Error("unable to generate token")
		return nil, err
	}

	token := &Token{
		AccessToken: tokenString,
		Expires:     tokenExpiration,
		Type:        "bearer",
	}

	return token, nil
}

// getToken is an internal method used to generate JWT Token
func (s *userService) getToken(user *User, expiration time.Time) (string, error) {

	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)

	// Create a map to store our claims
	claims := token.Claims.(jwt.MapClaims)


	// Set token claims
	claims["sub"] = user.ID
	claims["iss"] = "http://think.jholmestech.com"
	claims["exp"] = expiration.Unix()

	// Sign the token with our secret
	return token.SignedString(key)
}