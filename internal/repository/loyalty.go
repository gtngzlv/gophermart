package repository

import (
	"context"
	"database/sql"
	"time"

	"go.uber.org/zap"

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

func (p LoyaltyPostgres) WithdrawLoyalty(w model.WithdrawBalanceRequest, userID int, orderNumber string) error {
	queryUpdateCurrentBalance := `UPDATE ACCRUALS SET AMOUNT=AMOUNT-$1 WHERE USER_ID=$2 AND ORDER_NUMBER=$3`
	queryUpdateWithdrawal := `UPDATE WITHDRAWALS SET AMOUNT=AMOUNT+$1 WHERE USER_ID=$2 AND ORDER_NUMBER=$3`

	tx, err := p.db.Begin()
	if err != nil {
		tx.Rollback()
		p.log.Error("Error while begin tx")
		return err
	}
	_, err = tx.ExecContext(p.ctx, queryUpdateCurrentBalance, w.Sum, userID, orderNumber)
	if err != nil {
		p.log.Error("Failed to update current balance")
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(p.ctx, queryUpdateWithdrawal, w.Sum, userID, orderNumber)
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
	queryGet := `SELECT ORDER_NUMBER, AMOUNT, PROCESSED_AT 
			  FROM WITHDRAWALS 
			  WHERE USER_ID=$1 ORDER BY PROCESSED_AT`

	res, err := p.db.Query(queryGet, userID)
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
	queryAccruals := `SELECT sum(AMOUNT) FROM ACCRUALS WHERE USER_ID=$1`
	queryWithdrawn := `SELECT sum(AMOUNT) FROM WITHDRAWALS where user_id=$1`

	resAcc := p.db.QueryRow(queryAccruals, userID)
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

	resWith := p.db.QueryRow(queryWithdrawn, userID)
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
