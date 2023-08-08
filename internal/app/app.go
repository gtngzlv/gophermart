package app

import (
	"github.com/go-chi/chi"
	"github.com/gtngzlv/gophermart/internal/config"
	"github.com/gtngzlv/gophermart/internal/handlers"
	"github.com/gtngzlv/gophermart/internal/logger"
	"github.com/gtngzlv/gophermart/internal/storage"
	"net/http"
)

func Run() error {
	router := chi.NewRouter()
	cfg := config.LoadConfig()
	log := logger.NewLogger()
	s := storage.Init(cfg, log)
	h := handlers.NewHandler(cfg, log, router, s)
	return http.ListenAndServe(cfg.RunAddress, h.Router)
}
