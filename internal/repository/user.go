package repository

import (
	"context"
	"database/sql"

	"go.uber.org/zap"

	customErr "github.com/gtngzlv/gophermart/internal/errors"
	"github.com/gtngzlv/gophermart/internal/model"
	"github.com/gtngzlv/gophermart/internal/utils"
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

func (u *UserPostgres) Login(user model.User) error {
	userInDB, err := u.GetUserByLogin(user.Login)
	if err != nil {
		u.log.Errorf("DB Login: failed to get user by login")
		return err
	}
	if !utils.CheckHashAndPassword(userInDB.Password, user.Password) {
		return err
	}
	return nil
}

func (u *UserPostgres) Register(login, password string) error {
	query :=
		`INSERT INTO USERS(login, password) 
		VALUES($1, $2);`
	_, err := u.db.ExecContext(u.ctx, query, login, password)
	if err != nil {
		u.log.Errorf("DB Register: failed to exec query, %s", err)
		return err
	}
	return nil
}