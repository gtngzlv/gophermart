package storage

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"go.uber.org/zap"

	"github.com/gtngzlv/gophermart/internal/config"
	"github.com/gtngzlv/gophermart/internal/model"
)

type Storage interface {
	GetUserByLogin(login string) (model.User, error)
	Login(user model.User) error
	Register(login, password string) error

	GetBalance(userID int) (model.GetBalanceResponse, error)
	GetOrderByNumber(orderNumber string) (model.GetOrdersResponse, error)
	GetOrdersByUserID(userID int) ([]model.GetOrdersResponse, error)
	LoadOrder(orderNumber string, user model.User) error
	WithdrawLoyalty(withdrawal model.WithdrawBalanceRequest, userID int, orderNumber string) error
	GetWithdrawals(userID int) ([]model.GetWithdrawalsResponse, error)
	GetOrdersForProcessing(poolSize int) ([]string, error)
	UpdateOrderState(order *model.GetOrderAccrual) error
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
