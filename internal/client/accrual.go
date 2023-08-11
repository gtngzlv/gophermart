package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gtngzlv/gophermart/internal/model"
	"github.com/gtngzlv/gophermart/internal/storage"

	"golang.org/x/time/rate"
)

const tooManyRequestTemplate = "No more than %d requests per minute allowed"

type AccrualClient struct {
	db       storage.Storage
	host     string
	endpoint string
	poolSize int
	limiter  *rate.Limiter

	OrderQueue chan string
}

var wg sync.WaitGroup

func NewAccrualProcessing(db storage.Storage, host string, poolSize int) *AccrualClient {
	proc := &AccrualClient{
		db:         db,
		host:       host,
		endpoint:   "/api/orders/",
		poolSize:   poolSize,
		OrderQueue: make(chan string, poolSize),
	}
	for i := 0; i < poolSize; i++ {
		go proc.worker()
	}
	return proc
}

func (a *AccrualClient) worker() {
	for orderID := range a.OrderQueue {
		if a.limiter != nil && !a.limiter.Allow() {
			err := a.limiter.Wait(context.Background())
			if err != nil {
				log.Println(err)
				wg.Done()
				return
			}
		}
		receivedOrder, err := a.GetOrder(orderID)
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

func (a *AccrualClient) Run() {
	for {
		orderList, err := a.db.GetOrdersForProcessing(a.poolSize)
		if err != nil || len(orderList) == 0 {
			time.Sleep(5 * time.Second)
			continue
		}
		wg.Add(len(orderList))
		for _, orderID := range orderList {
			a.OrderQueue <- orderID
		}
		wg.Wait()
	}
}

func (a *AccrualClient) GetOrder(orderNum string) (*model.GetOrderAccrual, error) {
	res, err := http.Get(a.host + a.endpoint + orderNum)
	log.Println("getorder endpoint ", a.host+a.endpoint+orderNum)
	if err != nil {
		log.Println("GetOrder err http get", err)
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
		log.Print("GetOrder err", err)
	}
	return &orders, nil
}

func (a *AccrualClient) setLimit(n int) {
	if n <= 0 {
		a.limiter = nil
		return
	}
	a.limiter = rate.NewLimiter(rate.Every(time.Minute/time.Duration(n)), n)
}
