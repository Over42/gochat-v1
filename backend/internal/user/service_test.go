package user_test

import (
	"gochatv1/config"
	"gochatv1/db"
	"gochatv1/internal/user"

	"context"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	_ "github.com/lib/pq"
)

func TestServiceCreateUser(t *testing.T) {
	conn, tx, err := db.OpenTestDB()
	if err != nil {
		t.Fatalf("Failed to open test DB connection: %s", err)
	}
	defer db.CloseTestDB(tx, conn)

	cfg := config.New()
	val := validator.New()
	userRep := user.NewRepository(tx)
	userSvc := user.NewService(userRep, cfg, val)

	tests := []struct {
		name  string
		input *user.CreateUserReq
		want  *user.CreateUserRes
	}{
		{
			"Should create user",
			&user.CreateUserReq{Username: "test_user", Email: "test_user@gmail.com", Password: "password"},
			&user.CreateUserRes{Username: "test_user", Email: "test_user@gmail.com"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ans, _ := userSvc.CreateUser(context.Background(), test.input)
			if !cmp.Equal(ans, test.want) {
				t.Errorf("got %#v, want %#v", ans, test.want)
			}
		})
	}
}

func TestServiceLogin(t *testing.T) {
	conn, tx, err := db.OpenTestDB()
	if err != nil {
		t.Fatalf("Failed to open test DB connection: %s", err)
	}
	defer db.CloseTestDB(tx, conn)

	cfg := config.New()
	val := validator.New()
	userRep := user.NewRepository(tx)
	userSvc := user.NewService(userRep, cfg, val)

	tests := []struct {
		name  string
		input *user.LoginUserReq
		want  *user.LoginUserRes
	}{
		{
			"Should login user",
			&user.LoginUserReq{Email: "user@gmail.com", Password: "password"},
			&user.LoginUserRes{ID: "1", Username: "user"},
		},
		{
			"User does not exist",
			&user.LoginUserReq{Email: "user_notexist@gmail.com", Password: "password"},
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ans, _ := userSvc.Login(context.Background(), test.input)
			if !cmp.Equal(ans, test.want, cmpopts.IgnoreFields(user.LoginUserRes{}, "accessToken")) {
				t.Errorf("got %#v, want %#v", ans, test.want)
			}
		})
	}
}
