package app

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/gtngzlv/gophermart/internal/client"
	"github.com/gtngzlv/gophermart/internal/config"
	"github.com/gtngzlv/gophermart/internal/handlers"
	"github.com/gtngzlv/gophermart/internal/logger"
	"github.com/gtngzlv/gophermart/internal/repository"
)

func Run() error {
	router := chi.NewRouter()
	cfg := config.LoadConfig()
	log := logger.NewLogger()
	db, err := repository.InitPG(cfg.DatabaseAddress, log)
	if err != nil {
		log.Fatal(err)
	}

	repos := repository.NewRepository(context.Background(), db, log)
	accrualClient := client.NewAccrualProcessing(repos, cfg.AccrualSystemAddress, 10)
	go accrualClient.Run()

	h := handlers.NewHandler(cfg, log, router, repos)
	return http.ListenAndServe(cfg.RunAddress, h.Router)
}
