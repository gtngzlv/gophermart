package storage

import (
	"database/sql"
	"github.com/gtngzlv/gophermart/internal/errors"
	"github.com/gtngzlv/gophermart/internal/model"
	"github.com/gtngzlv/gophermart/internal/utils"
)

func (p PostgresDB) GetUserByLogin(login string) (model.User, error) {
	var user model.User
	query := `SELECT login, password FROM USERS WHERE LOGIN=$1`
	res := p.db.QueryRow(query, login)
	err := res.Scan(&user.Login, &user.Password)
	switch {
	case err == sql.ErrNoRows:
		return user, errors.ErrNoDBResult
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

func (p PostgresDB) GetBalance() {
	//TODO implement me
	panic("implement me")
}

func (p PostgresDB) LoadOrders() {
	//TODO implement me
	panic("implement me")
}

func (p PostgresDB) GetOrders() {
	//TODO implement me
	panic("implement me")
}

func (p PostgresDB) WithdrawLoyalty() {
	//TODO implement me
	panic("implement me")
}

func (p PostgresDB) GetWithdrawals() {
	//TODO implement me
	panic("implement me")
}
