package storage

import (
	"database/sql"

	customErr "github.com/gtngzlv/gophermart/internal/errors"
	"github.com/gtngzlv/gophermart/internal/model"
)

func (p PostgresDB) GetBalance(userID int) (model.GetBalanceResponse, error) {
	var balance model.GetBalanceResponse
	queryAccruals := `SELECT sum(AMOUNT) FROM ACCRUALS WHERE USER_ID=$1`
	queryWithdrawn := `SELECT sum(AMOUNT) FROM WITHDRAWALS where user_id=$1`

	resAcc := p.db.QueryRow(queryAccruals, userID)
	if resAcc.Err() != nil {
		return balance, resAcc.Err()
	}
	err := resAcc.Scan(&balance.Current)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return balance, customErr.ErrNoDBResult
		default:
			return balance, err
		}
	}

	resWith := p.db.QueryRow(queryWithdrawn, userID)
	if resWith.Err() != nil {
		return balance, resWith.Err()
	}
	err = resWith.Scan(&balance.Withdrawn)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return balance, customErr.ErrNoDBResult
		default:
			return balance, err
		}
	}
	return balance, nil
}
