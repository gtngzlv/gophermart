package client

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gtngzlv/gophermart/internal/model"
	"github.com/gtngzlv/gophermart/internal/repository"
)

type accrualClient struct {
	db       *repository.Repository
	host     string
	endpoint string
	poolSize int
	rw       sync.RWMutex
}

func NewAccrualClient(db *repository.Repository, host string, poolSize int) *accrualClient {
	proc := &accrualClient{
		db:       db,
		host:     host,
		endpoint: "/api/orders/",
		poolSize: poolSize,
	}
	return proc
}

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
			duration := returnRetryDuration(res)
			time.Sleep(time.Second * duration)
			a.rw.Unlock()
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

func returnRetryDuration(r *http.Response) time.Duration {
	headerParam := "Retry-After"
	duration := r.Header.Get(headerParam)
	convertedDuration, err := strconv.Atoi(duration)
	if err != nil {
		return 0
	}
	return time.Duration(convertedDuration)
}
