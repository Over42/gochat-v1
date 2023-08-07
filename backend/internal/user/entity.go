package user

import (
	"gochat/config"

	"context"

	"github.com/go-playground/validator/v10"
)

type User struct {
	ID       int64
	Username string
	Email    string
	Password string
}

type Service interface {
	CreateUser(ctx context.Context, req *CreateUserReq) (*CreateUserRes, error)
	Login(ctx context.Context, req *LoginUserReq) (*LoginUserRes, error)
}

type Repository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}

func Init(cfg *config.Config, val *validator.Validate, db DBTx) *Handler {
	userRep := NewRepository(db)
	userSvc := NewService(userRep, cfg, val)
	userHdl := NewHandler(userSvc)
	return userHdl
}
