package client

import (
	"log"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/gtngzlv/gophermart/internal/enums"
	"github.com/gtngzlv/gophermart/internal/repository"
)

const tooManyRequestTemplate = "No more than %d requests per minute allowed"

//

type accrualClient struct {
	db       *repository.Repository
	host     string
	endpoint string
	poolSize int
	rw       sync.RWMutex
}

func NewAccrualProcessing(db *repository.Repository, host string, poolSize int) *accrualClient {
	proc := &accrualClient{
		db:       db,
		host:     host,
		endpoint: "/api/orders/",
		poolSize: poolSize,
	}
	return proc
}

func (a *accrualClient) Run() {
	var errgr errgroup.Group
	errgr.SetLimit(a.poolSize)
	for {
		orderList, err := a.db.GetOrdersForProcessing(a.poolSize)
		if err != nil || len(orderList) == 0 {
			time.Sleep(5 * time.Second)
			continue
		}
		for _, orderID := range orderList {
			orderID := orderID
			errgr.Go(func() error {
				receivedOrder, err := a.GetOrderByNumber(orderID)
				if err != nil {
					return err
				}
				if receivedOrder != nil {
					log.Println("received order status ", *receivedOrder)
					if receivedOrder.Status == enums.StatusProcessed {
						if err = a.db.UpdateOrderStateProcessed(receivedOrder); err != nil {
							return err
						}
					} else if receivedOrder.Status == enums.StatusInvalid {
						if err = a.db.UpdateOrderStateInvalid(receivedOrder); err != nil {
							return err
						}
					}
				}
				return nil
			})
		}
	}
}
