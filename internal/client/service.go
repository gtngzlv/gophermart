package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"golang.org/x/time/rate"

	"github.com/gtngzlv/gophermart/internal/model"
)

func (a *accrualClient) worker() {
	for orderID := range a.OrderQueue {
		if a.limiter != nil && !a.limiter.Allow() {
			err := a.limiter.Wait(context.Background())
			if err != nil {
				log.Println(err)
				wg.Done()
				return
			}
		}
		receivedOrder, err := a.GetOrderByNumber(orderID)
		if err != nil {
			log.Println(err)
		}
		if receivedOrder != nil {
			log.Println("received order status ", *receivedOrder)
			if err = a.db.UpdateOrderState(receivedOrder); err != nil {
				log.Println(err)
			}
		}
		wg.Done()
	}
}

func (a *accrualClient) GetOrderByNumber(orderNum string) (*model.GetOrderAccrual, error) {
	res, err := http.Get(a.host + a.endpoint + orderNum)
	log.Println("getorder endpoint ", a.host+a.endpoint+orderNum)
	if err != nil {
		log.Println("GetOrderByNumber err http get", err)
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusTooManyRequests {
		var resBody []byte
		resBody, err = io.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		var rl int
		_, err = fmt.Sscanf(tooManyRequestTemplate, string(resBody), &rl)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		a.setLimit(rl)
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

func (a *accrualClient) setLimit(n int) {
	if n <= 0 {
		a.limiter = nil
		return
	}
	a.limiter = rate.NewLimiter(rate.Every(time.Minute/time.Duration(n)), n)
}