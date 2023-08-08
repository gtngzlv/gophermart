package storage

import (
	"database/sql"
	"github.com/gtngzlv/gophermart/internal/config"
	"github.com/gtngzlv/gophermart/internal/model"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"go.uber.org/zap"
)

type Storage interface {
	GetUserByLogin(login string) (model.User, error)
	Login(user model.User) error
	Register(login, password string) error
	GetBalance()
	LoadOrders()
	GetOrders()
	WithdrawLoyalty()
	GetWithdrawals()
}

type PostgresDB struct {
	log zap.SugaredLogger
	db  *sql.DB
}

func Init(cfg *config.AppConfig, log zap.SugaredLogger) Storage {
	var d Storage
	db, err := sql.Open("postgres", cfg.DatabaseAddress)
	if err != nil {
		return nil
	}
	if err = goose.SetDialect("postgres"); err != nil {
		log.Errorf("Init DB: failed while goose set dialect, %s", err)
		return nil
	}
	if err = goose.Up(db, "migrations"); err != nil {
		log.Errorf("Init DB: failed while goose up, %s", err)
		return nil
	}
	pg := PostgresDB{
		log: log,
		db:  db,
	}
	d = pg
	return d
}
