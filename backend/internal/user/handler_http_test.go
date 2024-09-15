package user_test

import (
	"gochatv1/config"
	"gochatv1/db"
	"gochatv1/internal/user"
	"gochatv1/router"

	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	_ "github.com/lib/pq"
)

func TestHandlerCreateUser(t *testing.T) {
	conn, tx, err := db.OpenTestDB()
	if err != nil {
		t.Fatalf("Failed to open test DB connection: %s", err)
	}
	defer db.CloseTestDB(tx, conn)

	cfg := config.New()
	val := validator.New()
	userHdl := user.Init(cfg, val, tx)
	rtr := router.InitRouter(cfg, userHdl, nil)

	tests := []struct {
		name  string
		input *user.CreateUserReq
		want  *user.CreateUserRes
		code  int
	}{
		{
			"Should create user",
			&user.CreateUserReq{Username: "test_user", Email: "test_user@gmail.com", Password: "password"},
			&user.CreateUserRes{Username: "test_user", Email: "test_user@gmail.com"},
			http.StatusOK,
		},
		{
			"Bad request no password",
			&user.CreateUserReq{Username: "test_user2", Email: "test_user2@gmail.com", Password: ""},
			&user.CreateUserRes{Username: "", Email: ""},
			http.StatusBadRequest,
		},
		{
			"User email already exists",
			&user.CreateUserReq{Username: "user", Email: "user@gmail.com", Password: "password"},
			&user.CreateUserRes{Username: "", Email: ""},
			http.StatusInternalServerError,
		},
		{
			"Empty username",
			&user.CreateUserReq{Username: "", Email: "user@gmail.com", Password: "password"},
			&user.CreateUserRes{Username: "", Email: ""},
			http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			c := gin.CreateTestContextOnly(recorder, rtr)

			reqBody, _ := json.Marshal(test.input)
			req := httptest.NewRequest("POST", "/signup", bytes.NewReader(reqBody))
			c.Request = req

			userHdl.CreateUser(c)

			resBody, _ := io.ReadAll(recorder.Body)
			res := &user.CreateUserRes{}
			_ = json.Unmarshal(resBody, res)

			if !cmp.Equal(res, test.want) {
				t.Errorf("got %#v, want %#v", res, test.want)
			}

			if recorder.Code != test.code {
				t.Errorf("got %d, want %d", recorder.Code, test.code)
			}
		})
	}
}

func TestHandlerLogin(t *testing.T) {
	conn, tx, err := db.OpenTestDB()
	if err != nil {
		t.Fatalf("Failed to open test DB connection: %s", err)
	}
	defer db.CloseTestDB(tx, conn)

	cfg := config.New()
	val := validator.New()
	userHdl := user.Init(cfg, val, tx)
	rtr := router.InitRouter(cfg, userHdl, nil)

	tests := []struct {
		name       string
		input      *user.LoginUserReq
		want       *user.LoginUserRes
		code       int
		wantCookie bool
		cookie     string
	}{
		{
			"Should login user",
			&user.LoginUserReq{Email: "user@gmail.com", Password: "password"},
			&user.LoginUserRes{ID: "1", Username: "user"},
			http.StatusOK,
			true,
			"jwt",
		},
		{
			"User does not exist",
			&user.LoginUserReq{Email: "user_notexist@gmail.com", Password: "password"},
			&user.LoginUserRes{ID: "", Username: ""},
			http.StatusInternalServerError,
			false,
			"",
		},
		{
			"Bad request no password",
			&user.LoginUserReq{Email: "user@gmail.com", Password: ""},
			&user.LoginUserRes{ID: "", Username: ""},
			http.StatusBadRequest,
			false,
			"",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			c := gin.CreateTestContextOnly(recorder, rtr)

			reqBody, _ := json.Marshal(test.input)
			req := httptest.NewRequest("POST", "/login", bytes.NewReader(reqBody))
			c.Request = req

			userHdl.Login(c)

			resBody, _ := io.ReadAll(recorder.Body)
			res := &user.LoginUserRes{}
			_ = json.Unmarshal(resBody, res)

			if !cmp.Equal(res, test.want, cmpopts.IgnoreFields(user.LoginUserRes{}, "accessToken")) {
				t.Errorf("got %#v, want %#v", res, test.want)
			}

			if recorder.Code != test.code {
				t.Errorf("got %d, want %d", recorder.Code, test.code)
			}

			if (test.wantCookie) && (recorder.Result().Cookies()[0].Name != test.cookie) {
				t.Error("auth cookie was not created")
			}
		})
	}
}

func TestHandlerLogout(t *testing.T) {
	conn, tx, err := db.OpenTestDB()
	if err != nil {
		t.Fatalf("Failed to open test DB connection: %s", err)
	}
	defer db.CloseTestDB(tx, conn)

	cfg := config.New()
	val := validator.New()
	userHdl := user.Init(cfg, val, tx)
	rtr := router.InitRouter(cfg, userHdl, nil)

	tests := []struct {
		name string
		want string
		code int
	}{
		{
			"Should logout user",
			"logout successful",
			http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			c := gin.CreateTestContextOnly(recorder, rtr)

			req := httptest.NewRequest("GET", "/logout", nil)
			c.Request = req

			userHdl.Logout(c)

			resBody, _ := io.ReadAll(recorder.Body)
			res := make(map[string]string)
			_ = json.Unmarshal(resBody, &res)

			if !cmp.Equal(res["message"], test.want) {
				t.Errorf("got %s, want %s", res, test.want)
			}

			if recorder.Code != test.code {
				t.Errorf("got %d, want %d", recorder.Code, test.code)
			}

			if recorder.Result().Cookies()[0].Value != "" {
				t.Error("auth cookie was not deleted")
			}
		})
	}
}
