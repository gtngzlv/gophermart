package repository

import (
	"context"
	"database/sql"
	"time"

	"go.uber.org/zap"

	"github.com/gtngzlv/gophermart/internal/enums"
	customErr "github.com/gtngzlv/gophermart/internal/errors"
	"github.com/gtngzlv/gophermart/internal/model"
)

type LoyaltyPostgres struct {
	ctx context.Context
	db  *sql.DB
	log zap.SugaredLogger
}

func NewLoyaltyPostgres(ctx context.Context, db *sql.DB, log zap.SugaredLogger) *LoyaltyPostgres {
	return &LoyaltyPostgres{
		ctx: ctx,
		db:  db,
		log: log,
	}
}

func (p LoyaltyPostgres) DeductPoints(w model.WithdrawBalanceRequest, userID int, orderNumber string) error {
	queryUpdateCurrentBalance := `UPDATE users SET balance=balance-$1 WHERE ID=$2`
	queryUpdateWithdrawal := `UPDATE orders SET AMOUNT=$1 WHERE user_id=$2 AND NUMBER=$3 AND OPERATION_TYPE=$4`

	tx, err := p.db.Begin()
	if err != nil {
		tx.Rollback()
		p.log.Error("Error while begin tx")
		return err
	}
	_, err = tx.ExecContext(p.ctx, queryUpdateCurrentBalance,
		w.Sum,
		userID)
	if err != nil {
		p.log.Error("Failed to update current balance")
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(p.ctx, queryUpdateWithdrawal,
		w.Sum,
		userID,
		orderNumber,
		enums.Withdrawal)
	if err != nil {
		p.log.Error("Failed to update withdrawal")
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (p LoyaltyPostgres) GetWithdrawals(userID int) ([]*model.GetWithdrawalsResponse, error) {
	var (
		withdrawal model.GetWithdrawalsResponse
		response   []*model.GetWithdrawalsResponse
		err        error
	)
	queryGet := `SELECT NUMBER, AMOUNT, UPLOADED_AT
			  FROM ORDERS
			  WHERE USER_ID=$1 AND OPERATION_TYPE=$2 ORDER BY UPLOADED_AT`

	res, err := p.db.Query(queryGet, userID, enums.Withdrawal)
	if res.Err() != nil {
		return nil, res.Err()
	}
	if err != nil {
		return nil, err
	}
	for res.Next() {
		err = res.Scan(&withdrawal.Order, &withdrawal.Sum, &withdrawal.ProcessedAt)
		if err != nil {
			return nil, err
		}
		withdrawal.ProcessedAt.Format(time.RFC3339)
		if withdrawal.Sum == 0 {
			continue
		}
		response = append(response, &withdrawal)
	}
	return response, nil
}

func (p LoyaltyPostgres) GetBalance(userID int) (*model.GetBalanceResponse, error) {
	var balance model.GetBalanceResponse
	current := `SELECT balance FROM users WHERE ID=$1`
	withdrawn := `SELECT sum(AMOUNT) FROM orders where user_id=$1 and operation_type=$2`

	resAcc := p.db.QueryRow(current, userID)
	if resAcc.Err() != nil {
		return nil, resAcc.Err()
	}
	err := resAcc.Scan(&balance.Current)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, customErr.ErrNoDBResult
		default:
			return nil, err
		}
	}

	resWith := p.db.QueryRow(withdrawn, userID, enums.Withdrawal)
	if resWith.Err() != nil {
		return nil, resWith.Err()
	}
	err = resWith.Scan(&balance.Withdrawn)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, customErr.ErrNoDBResult
		default:
			return nil, err
		}
	}
	return &balance, nil
}
