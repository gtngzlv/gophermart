package client

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gtngzlv/gophermart/internal/model"
)

func (a *accrualClient) GetOrderByNumber(orderNum string) (*model.GetOrderAccrual, error) {
	a.rw.RLock()
	res, err := http.Get(a.host + a.endpoint + orderNum)
	a.rw.RUnlock()
	log.Println("getorder endpoint ", a.host+a.endpoint+orderNum)
	if err != nil {
		log.Println("GetOrderByNumber err http get", err)
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusTooManyRequests {
		b := a.rw.TryLock()
		if b {
			time.Sleep(time.Second * 60)
			a.rw.Unlock()
			// поправить на кастомный эррор
			return nil, errors.New("429 error code")
		}
	}

	if res.StatusCode != http.StatusOK {
		return nil, err
	}

	var orders model.GetOrderAccrual
	if err = json.NewDecoder(res.Body).Decode(&orders); err != nil {
		log.Print("GetOrderByNumber err", err)
	}
	return &orders, nil
}
