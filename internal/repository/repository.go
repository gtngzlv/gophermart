package repository

import (
	"context"
	"database/sql"

	"go.uber.org/zap"

	"github.com/gtngzlv/gophermart/internal/model"
)

type User interface {
	GetUserByLogin(login string) (*model.User, error)
	Register(login, password string) error
}

type Order interface {
	GetOrderByNumber(orderNumber string) (*model.GetOrdersResponse, error)
	LoadOrder(orderNumber string, user model.User) error
	GetOrdersByUserID(userID int) ([]*model.GetOrdersResponse, error)
	GetOrdersForProcessing(poolSize int) ([]string, error)
	UpdateOrderState(order *model.GetOrderAccrual) error
}

type Loyalty interface {
	DeductPoints(w model.WithdrawBalanceRequest, userID int, orderNumber string) error
	GetWithdrawals(userID int) ([]*model.GetWithdrawalsResponse, error)
	GetBalance(userID int) (*model.GetBalanceResponse, error)
}

type Repository struct {
	User
	Order
	Loyalty
}

func NewRepository(ctx context.Context, db *sql.DB, log zap.SugaredLogger) *Repository {
	return &Repository{
		User:    NewUserPostgres(ctx, db, log),
		Order:   NewOrderPostgres(ctx, db, log),
		Loyalty: NewLoyaltyPostgres(ctx, db, log),
	}
}
