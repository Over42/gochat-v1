package user

import (
	"context"
	"database/sql"
)

// Makes possible to inject DB connection (in prod) or Tx transaction (in tests)
type DBTx interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type repository struct {
	db DBTx
}

func NewRepository(db DBTx) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(ctx context.Context, user *User) (*User, error) {
	query := "INSERT INTO users(username, password, email) VALUES ($1, $2, $3)"
	_, err := r.db.ExecContext(ctx, query, user.Username, user.Password, user.Email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	user := User{}
	query := "SELECT id, email, username, password FROM users WHERE email = $1"
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
