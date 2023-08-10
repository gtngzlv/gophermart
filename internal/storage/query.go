package storage

import (
	"database/sql"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"

	customErr "github.com/gtngzlv/gophermart/internal/errors"
	"github.com/gtngzlv/gophermart/internal/model"
	"github.com/gtngzlv/gophermart/internal/utils"
)

func (p PostgresDB) GetUserByLogin(login string) (model.User, error) {
	var user model.User
	query := `SELECT id, login, password FROM USERS WHERE LOGIN=$1`
	res := p.db.QueryRow(query, login)
	err := res.Scan(&user.ID, &user.Login, &user.Password)
	switch {
	case err == sql.ErrNoRows:
		return user, customErr.ErrNoDBResult
	case err != nil:
		return user, err
	default:
		return user, nil
	}
}

func (p PostgresDB) Login(user model.User) error {
	userInDB, err := p.GetUserByLogin(user.Login)
	if err != nil {
		p.log.Errorf("DB Login: failed to get user by login")
		return err
	}
	if !utils.CheckHashAndPassword(userInDB.Password, user.Password) {
		return err
	}
	return nil
}

func (p PostgresDB) Register(login, password string) error {
	query :=
		`INSERT INTO USERS(login, password) 
		VALUES($1, $2);`
	_, err := p.db.Exec(query, login, password)
	if err != nil {
		p.log.Errorf("DB Register: failed to exec query, %s", err)
		return err
	}
	return nil
}

func (p PostgresDB) GetBalance(userID int) (model.GetBalanceResponse, error) {
	var balance model.GetBalanceResponse
	query := `SELECT CURRENT_BALANCE, WITHDRAWN FROM BALANCE WHERE USER_ID=$1`
	res := p.db.QueryRow(query, userID)
	if res.Err() != nil {
		return balance, res.Err()
	}
	err := res.Scan(&balance.Current, &balance.Withdrawn)
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

func (p PostgresDB) GetOrderByNumber(orderNumber string) (model.GetOrdersResponse, error) {
	var (
		order model.GetOrdersResponse
	)
	query := `SELECT NUMBER, USER_ID, STATUS, ACCRUAL, UPLOADED_AT FROM ORDERS WHERE NUMBER=$1`
	res := p.db.QueryRow(query, orderNumber)
	err := res.Scan(&order.Number, &order.UserID, &order.Status, &order.Accrual, &order.UploadedAt)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			{
				return order, customErr.ErrNoDBResult
			}
		default:
			{
				return order, err
			}
		}
	}
	return order, nil
}

func (p PostgresDB) LoadOrder(orderNumber string, user model.User) error {
	query := `INSERT INTO ORDERS(NUMBER, USER_ID, STATUS, ACCRUAL, UPLOADED_AT) 
			  VALUES($1, $2, $3, $4, $5)`
	_, err := p.db.Exec(query, orderNumber, user.ID, StatusNew, 0, time.Now())
	if err != nil {
		if pgerrcode.IsIntegrityConstraintViolation(string(err.(*pq.Error).Code)) {
			return customErr.ErrDuplicateValue
		}
		p.log.Errorf("DB LoadOrder: failed to exec query, %s", err)
		return err
	}
	return nil
}

func (p PostgresDB) GetOrdersByUserID(userID int) ([]model.GetOrdersResponse, error) {
	var (
		order  model.GetOrdersResponse
		orders []model.GetOrdersResponse
		err    error
	)
	query := `SELECT NUMBER, STATUS, ACCRUAL, UPLOADED_AT FROM ORDERS WHERE USER_ID=$1`
	res, err := p.db.Query(query, userID)
	if res.Err() != nil {
		return orders, res.Err()
	}
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			{
				return orders, customErr.ErrNoDBResult
			}
		default:
			return orders, err
		}
	}
	for res.Next() {
		err = res.Scan(&order.Number, &order.Status, &order.Accrual, &order.UploadedAt)
		if err != nil {
			return []model.GetOrdersResponse{}, err
		}
		order.UploadedAt.Format(time.RFC3339)
		orders = append(orders, order)
	}
	return orders, nil
}

func (p PostgresDB) WithdrawLoyalty() {
	//TODO implement me
	panic("implement me")
}

func (p PostgresDB) GetWithdrawals() {
	//TODO implement me
	panic("implement me")
}
