package user_test

import (
	"gochat/db"
	"gochat/internal/user"

	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	_ "github.com/lib/pq"
)

func TestRepositoryCreateUser(t *testing.T) {
	conn, tx, err := db.OpenTestDB()
	if err != nil {
		t.Fatalf("Failed to open test DB connection: %s", err)
	}
	defer db.CloseTestDB(tx, conn)

	userRep := user.NewRepository(tx)

	tests := []struct {
		name  string
		input *user.User
		want  *user.User
	}{
		{
			"Should create user",
			&user.User{Username: "test_user", Password: "password", Email: "test_user@gmail.com"},
			&user.User{ID: 1, Username: "test_user", Password: "password", Email: "test_user@gmail.com"},
		},
		{
			"User email already exists",
			&user.User{Username: "user", Email: "user@gmail.com", Password: "password"},
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ans, _ := userRep.CreateUser(context.Background(), test.input)
			// ID always incremented after transaction
			if !cmp.Equal(ans, test.want, cmpopts.IgnoreFields(user.User{}, "ID")) {
				t.Errorf("got %#v, want %#v", ans, test.want)
			}
		})
	}
}

func TestRepositoryGetUserByEmail(t *testing.T) {
	conn, tx, err := db.OpenTestDB()
	if err != nil {
		t.Fatalf("Failed to open test DB connection: %s", err)
	}
	defer db.CloseTestDB(tx, conn)

	userRep := user.NewRepository(tx)

	tests := []struct {
		name  string
		input string
		want  *user.User
	}{
		{
			"Should get user",
			"user@gmail.com",
			&user.User{ID: 1, Username: "user", Password: "password", Email: "user@gmail.com"},
		},
		{
			"User does not exist",
			"user_nonexist@gmail.com",
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ans, _ := userRep.GetUserByEmail(context.Background(), test.input)
			// Don't compare with hashed password
			if !cmp.Equal(ans, test.want, cmpopts.IgnoreFields(user.User{}, "Password")) {
				t.Errorf("got %#v, want %#v", ans, test.want)
			}
		})
	}
}
