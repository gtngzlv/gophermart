package app

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/gtngzlv/gophermart/internal/client"
	"github.com/gtngzlv/gophermart/internal/config"
	"github.com/gtngzlv/gophermart/internal/handlers"
	"github.com/gtngzlv/gophermart/internal/logger"
	"github.com/gtngzlv/gophermart/internal/storage"
)

func Run() error {
	router := chi.NewRouter()
	cfg := config.LoadConfig()
	log := logger.NewLogger()
	s := storage.Init(cfg, log)
	accrualClient := client.NewAccrualProcessing(s, cfg.AccrualSystemAddress, 10)
	go accrualClient.Run()

	h := handlers.NewHandler(cfg, log, router, s)
	return http.ListenAndServe(cfg.RunAddress, h.Router)
}
