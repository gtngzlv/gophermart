package repository

import (
	"context"
	"database/sql"

	"go.uber.org/zap"

	customErr "github.com/gtngzlv/gophermart/internal/errors"
	"github.com/gtngzlv/gophermart/internal/model"
)

type UserPostgres struct {
	ctx context.Context
	db  *sql.DB
	log zap.SugaredLogger
}

func NewUserPostgres(ctx context.Context, db *sql.DB, log zap.SugaredLogger) *UserPostgres {
	return &UserPostgres{
		ctx: ctx,
		db:  db,
		log: log,
	}
}

func (u *UserPostgres) GetUserByLogin(login string) (*model.User, error) {
	var user model.User
	query := `SELECT id, login, password FROM USERS WHERE LOGIN=$1`
	res := u.db.QueryRowContext(u.ctx, query, login)
	err := res.Scan(&user.ID, &user.Login, &user.Password)
	switch {
	case err == sql.ErrNoRows:
		return nil, customErr.ErrNoDBResult
	case err != nil:
		return nil, err
	default:
		return &user, nil
	}
}

func (u *UserPostgres) Register(login, password string) (int, error) {
	var id int
	query :=
		`INSERT INTO USERS(login, password) VALUES($1, $2) RETURNING Users.id;`
	res := u.db.QueryRowContext(u.ctx, query, login, password)
	err := res.Scan(&id)
	if err != nil {
		u.log.Errorf("DB Register: failed to exec query, %s", err)
		return 0, err
	}
	return id, nil
}
