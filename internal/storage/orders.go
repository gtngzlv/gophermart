package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"

	customErr "github.com/gtngzlv/gophermart/internal/errors"
	"github.com/gtngzlv/gophermart/internal/model"
)

func (p PostgresDB) GetOrderByNumber(orderNumber string) (model.GetOrdersResponse, error) {
	p.log.Info("GetOrderByNumber: provided order num is ", orderNumber)
	var (
		order model.GetOrdersResponse
	)
	query := `SELECT ID, NUMBER, USER_ID, UPLOADED_AT FROM ORDERS WHERE NUMBER=$1`
	res := p.db.QueryRow(query, orderNumber)
	err := res.Scan(&order.ID, &order.Number, &order.UserID, &order.UploadedAt)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			{
				p.log.Info("GetOrderByNumber err is", err)
				return order, customErr.ErrNoDBResult
			}
		default:
			{
				p.log.Info("GetOrderByNumber err is", err)
				return order, err
			}
		}
	}
	return order, nil
}

func (p PostgresDB) LoadOrder(orderNumber string, user model.User) error {
	queryOrders := `INSERT INTO ORDERS(NUMBER, USER_ID, UPLOADED_AT) 
			  VALUES($1, $2, $3)`
	queryAccruals := `INSERT INTO ACCRUALS(ORDER_NUMBER, USER_ID, UPLOADED_AT) VALUES($1, $2, $3)`
	queryWithdrawn := `INSERT INTO WITHDRAWALS(ORDER_NUMBER, USER_ID) VALUES($1, $2)`

	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			txErr := tx.Rollback()
			if txErr != nil {
				err = fmt.Errorf("LoadOrder: failed to rollback %s", txErr.Error())
			}
		}
	}()
	_, err = tx.ExecContext(context.Background(), queryOrders, orderNumber, user.ID, time.Now())
	if err != nil {
		if pgerrcode.IsIntegrityConstraintViolation(string(err.(*pq.Error).Code)) {
			return customErr.ErrDuplicateValue
		}
		p.log.Errorf("DB LoadOrder: failed to exec query insert into orders, %s", err)
		return err
	}

	_, err = tx.ExecContext(context.Background(), queryAccruals, orderNumber, user.ID, time.Now())
	if err != nil {
		p.log.Errorf("DB LoadOrder: failed to exec query insert into accruals, %s", err)
		return err
	}

	_, err = tx.ExecContext(context.Background(), queryWithdrawn, orderNumber, user.ID)
	if err != nil {
		p.log.Errorf("DB LoadOrder: failed to exec query insert into withdrawals, %s", err)
		return err
	}

	return tx.Commit()
}

func (p PostgresDB) GetOrdersByUserID(userID int) ([]model.GetOrdersResponse, error) {
	var (
		order  model.GetOrdersResponse
		orders []model.GetOrdersResponse
		err    error
	)
	query := "SELECT order_number, status, amount, uploaded_at from accruals where user_id=$1"
	rows, err := p.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	for rows.Next() {
		err = rows.Scan(&order.Number, &order.Status, &order.Accrual, &order.UploadedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}

func (p PostgresDB) GetOrdersForProcessing(poolSize int) ([]string, error) {
	var orders []string
	rows, err := p.db.Query(
		"SELECT order_number FROM accruals WHERE status IN ($1, $2) ORDER BY uploaded_at LIMIT $3", "NEW", "PROCESSING", poolSize,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
	}()

	for rows.Next() {
		var orderID string
		if err = rows.Scan(&orderID); err != nil {
			return orders, err
		}
		orders = append(orders, orderID)
	}
	err = rows.Err()
	return orders, err
}

func (p PostgresDB) UpdateOrderState(order *model.GetOrderAccrual) error {
	res, err := p.db.Exec(
		"UPDATE accruals SET status=$1, amount=$2 WHERE order_number = $3",
		order.Status, order.Accrual, order.Order,
	)
	if err != nil {
		return err
	}
	rAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	p.log.Info("UpdateOrderState affected rows count", rAffected)
	return err
}
