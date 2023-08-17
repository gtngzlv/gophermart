package handlers

import (
	"github.com/go-chi/chi"
	"go.uber.org/zap"

	"github.com/gtngzlv/gophermart/internal/auth"
	"github.com/gtngzlv/gophermart/internal/config"
	"github.com/gtngzlv/gophermart/internal/logger"
	"github.com/gtngzlv/gophermart/internal/repository"
)

type Handler struct {
	Router *chi.Mux
	log    zap.SugaredLogger
	cfg    *config.AppConfig
	repo   *repository.Repository
}

func NewHandler(cfg *config.AppConfig, log zap.SugaredLogger, m *chi.Mux, r *repository.Repository) *Handler {
	h := &Handler{
		Router: m,
		log:    log,
		cfg:    cfg,
		repo:   r,
	}
	h.init()
	return h
}

func (h *Handler) init() {
	h.Router.Use(logger.WithLogging)
	h.Router.Use(auth.Authorization)
	h.Router.Post("/api/user/register", h.Register)
	h.Router.Post("/api/user/login", h.Login)
	h.Router.Post("/api/user/orders", h.LoadOrders)
	h.Router.Post("/api/user/balance/withdraw", h.WithdrawLoyalty)

	h.Router.Get("/api/user/orders", h.GetOrders)
	h.Router.Get("/api/user/balance", h.GetBalance)
	h.Router.Get("/api/user/withdrawals", h.GetWithdrawals)
}
