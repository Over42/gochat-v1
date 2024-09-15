package user

import (
	"gochatv1/config"

	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repository Repository
	config     *config.Config
	validate   *validator.Validate
}

func NewService(repo Repository, cfg *config.Config, val *validator.Validate) Service {
	return &service{
		repo,
		cfg,
		val,
	}
}

type CreateUserReq struct {
	Username string `json:"username" validate:"required,min=3"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type CreateUserRes struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (s *service) CreateUser(ctx context.Context, req *CreateUserReq) (*CreateUserRes, error) {
	err := s.validate.Struct(req)
	if err != nil {
		return nil, err
	}

	context, cancel := context.WithTimeout(ctx, s.config.DBTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("password hashing failed: %w", err)
	}

	newUser := &User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	user, err := s.repository.CreateUser(context, newUser)
	if err != nil {
		return nil, err
	}

	res := &CreateUserRes{
		Username: user.Username,
		Email:    user.Email,
	}

	return res, nil
}

type LoginUserReq struct {
	Email    string `json:"email"    validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginUserRes struct {
	accessToken string
	ID          string `json:"id"`
	Username    string `json:"username"`
}

type JWTClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (s *service) Login(ctx context.Context, req *LoginUserReq) (*LoginUserRes, error) {
	err := s.validate.Struct(req)
	if err != nil {
		return nil, err
	}

	context, cancel := context.WithTimeout(ctx, s.config.DBTimeout)
	defer cancel()

	user, err := s.repository.GetUserByEmail(context, req.Email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		ID:       strconv.FormatInt(user.ID, 10),
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    strconv.FormatInt(user.ID, 10),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	signedToken, err := token.SignedString([]byte(s.config.JWTKey))
	if err != nil {
		return nil, err
	}

	res := &LoginUserRes{
		accessToken: signedToken,
		ID:          strconv.FormatInt(user.ID, 10),
		Username:    user.Username,
	}

	return res, nil
}
