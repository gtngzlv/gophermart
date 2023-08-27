package client

import (
	"sync"
	"time"

	"golang.org/x/time/rate"

	"github.com/gtngzlv/gophermart/internal/repository"
)

const tooManyRequestTemplate = "No more than %d requests per minute allowed"

type accrualClient struct {
	db       *repository.Repository
	host     string
	endpoint string
	poolSize int
	limiter  *rate.Limiter

	OrderQueue chan string
}

var wg sync.WaitGroup

func NewAccrualProcessing(db *repository.Repository, host string, poolSize int) *accrualClient {
	proc := &accrualClient{
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

func (a *accrualClient) Run() {
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
