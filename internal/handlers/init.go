package handlers

import (
	"github.com/go-chi/chi"
	"go.uber.org/zap"

	"github.com/go-chi/chi/v5/middleware"

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
	h.Router.Use(middleware.Compress(5, "text/html",
		"application/x-gzip",
		"text/plain",
		"application/json"))
	h.Router.Post("/api/user/register", h.Register)
	h.Router.Post("/api/user/login", h.Login)

	h.Router.Group(func(r chi.Router) {
		r.Use(auth.Authorization)
		r.Post("/api/user/orders", h.LoadOrders)
		r.Post("/api/user/balance/withdraw", h.DeductPoints)

		r.Get("/api/user/orders", h.GetOrders)
		r.Get("/api/user/balance", h.GetBalance)
		r.Get("/api/user/withdrawals", h.GetWithdrawals)
	})

}
