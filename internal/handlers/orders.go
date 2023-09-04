package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ShiraazMoollatjie/goluhn"

	"github.com/gtngzlv/gophermart/internal/auth"
	"github.com/gtngzlv/gophermart/internal/errors"
)

func (h *Handler) LoadOrders(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromToken(w, r)
	contentType := r.Header["Content-Type"]
	if contentType[0] != "text/plain" {
		h.log.Infof("Received non text/plain")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.Errorf("PostURL: error: %s while reading body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	orderNum := string(body)
	h.log.Info("LoadOrders: order num in body", orderNum)

	if err = goluhn.Validate(orderNum); err != nil {
		h.log.Info("Provided order num invalid")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	existingOrder, err := h.repo.GetOrderByNumber(orderNum)
	if err != nil {
		switch err {
		case errors.ErrNoDBResult:
			{
				break
			}
		default:
			{
				h.log.Errorf("LoadOrder: failed to check if exist, %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}
	if existingOrder != nil && existingOrder.UserID == userID {
		h.log.Info("Provided order num %s already loaded with by user", orderNum)
		w.WriteHeader(200)
		return
	}
	if existingOrder != nil && existingOrder.UserID != userID {
		h.log.Infof("Provided order num %s already exist", orderNum)
		w.WriteHeader(http.StatusConflict)
		return
	}

	// загрузили заказ
	if err = h.repo.LoadOrder(orderNum, userID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.log.Info("LoadOrders: saved order with number", orderNum)
	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	userID := auth.GetUserIDFromToken(w, r)

	orders, err := h.repo.GetOrdersByUserID(userID)
	if err != nil {
		switch err {
		case errors.ErrNoDBResult:
			{
				h.log.Infof("GetOrdersByUserID: no orders for user with id %s", userID)
				w.WriteHeader(http.StatusNoContent)
				return
			}
		default:
			{
				h.log.Errorf("getUserInfoByToken: failed, %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}
	resp, err := json.Marshal(orders)
	if err != nil {
		h.log.Errorf("GetOrders: failed to marshal resp %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
