package main

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strings"
	"time"
)

var (
	key             = os.Getenv("JWT_SECRET")
	ErrInvalidLogin = errors.New("invalid login")
	defaultuser     = os.Getenv("DEFAULT_USER")
	defaultpassword = os.Getenv("DEFAULT_PASSWORD")
)

type User struct {
	ID       string `firestore:"id"`
	Username string `firestore:"username"`
	Password string `firestore:"password"`
}

type Token struct {
	Type        string    `json:"type"`
	AccessToken string    `json:"token"`
	Expires     time.Time `json:"expires"`
}

type UserService interface {
	Authenticate(ctx context.Context, username, password string) (*Token, error)
}

type UserRepository interface {
	FindUserByUsername(ctx context.Context, username string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
}

func NewUserFirestoreRepository(database *UserFirestoreDatabase) UserRepository {
	return database
}

func NewUserPostgresRepository(database *PostgresDatabase) UserRepository {
	return database
}

type userService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	if key == "" {
		logrus.Fatal("jwt not set")
	}

	s := &userService{
		repo: repo,
	}

	logrus.Info("checking it see if user exists")
	if defaultuser == "" {
		logrus.Fatal("default username not set")
	}
	ctx := context.Background()
	user, err := repo.FindUserByUsername(ctx, defaultuser)
	if err != nil {
		logrus.WithError(err).Fatal("cannot contact firebase")
	}
	if user == nil {
		logrus.Info("default user does not exist creating")
		if defaultpassword == "" {
			logrus.Fatal("default username not set")
		}
		if err := s.createUser(ctx, defaultuser, defaultpassword); err != nil {
			logrus.WithError(err).Fatal("unable to create default user")
		}
		logrus.Info("default user created")
	}

	return s
}

func (s *userService) createUser(ctx context.Context, username, password string) error {
	username = strings.ToLower(username)
	ePwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logrus.WithError(err).Error("unable to encrypt password")
		return errors.New("unable to encrypt password")
	}

	user := &User{
		ID:       uuid.New().String(),
		Username: username,
		Password: string(ePwd),
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		logrus.WithError(err).Error("unable to create user")
		return errors.New("unable to create user")
	}

	return nil
}

func (s *userService) Authenticate(ctx context.Context, username, password string) (*Token, error) {
	username = strings.ToLower(username)
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

	tokenExpiration := time.Now().Add(time.Hour * 72)
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
	return token.SignedString([]byte(key))
}
