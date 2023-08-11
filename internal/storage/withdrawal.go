package storage

import (
	"context"
	"time"

	"github.com/gtngzlv/gophermart/internal/model"
)

func (p PostgresDB) WithdrawLoyalty(w model.WithdrawBalanceRequest, userID int, orderNumber string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	queryUpdateCurrentBalance := `UPDATE ACCRUALS SET AMOUNT=AMOUNT-$1 WHERE USER_ID=$2 AND ORDER_NUMBER=$3`
	queryUpdateWithdrawal := `UPDATE WITHDRAWALS SET AMOUNT=AMOUNT+$1 WHERE USER_ID=$2 AND ORDER_NUMBER=$3`

	tx, err := p.db.Begin()
	if err != nil {
		tx.Rollback()
		p.log.Error("Error while begin tx")
		return err
	}
	_, err = tx.ExecContext(ctx, queryUpdateCurrentBalance, w.Sum, userID, orderNumber)
	if err != nil {
		p.log.Error("Failed to update current balance")
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, queryUpdateWithdrawal, w.Sum, userID, orderNumber)
	if err != nil {
		p.log.Error("Failed to update withdrawal")
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (p PostgresDB) GetWithdrawals(userID int) ([]model.GetWithdrawalsResponse, error) {
	var (
		withdrawal model.GetWithdrawalsResponse
		response   []model.GetWithdrawalsResponse
		err        error
	)
	queryGet := `SELECT ORDER_NUMBER, AMOUNT, PROCESSED_AT 
			  FROM WITHDRAWALS 
			  WHERE USER_ID=$1 ORDER BY PROCESSED_AT`

	res, err := p.db.Query(queryGet, userID)
	if res.Err() != nil {
		return response, res.Err()
	}
	if err != nil {
		return response, err
	}
	for res.Next() {
		err = res.Scan(&withdrawal.Order, &withdrawal.Sum, &withdrawal.ProcessedAt)
		if err != nil {
			return []model.GetWithdrawalsResponse{}, err
		}
		withdrawal.ProcessedAt.Format(time.RFC3339)
		if withdrawal.Sum == 0 {
			continue
		}
		response = append(response, withdrawal)
	}
	return response, nil

}
