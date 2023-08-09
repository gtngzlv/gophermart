package handlers

import (
	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/gtngzlv/gophermart/internal/errors"
	"io"
	"net/http"
)

func (h *Handler) LoadOrders(w http.ResponseWriter, r *http.Request) {
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

	if err = goluhn.Validate(orderNum); err != nil {
		h.log.Info("Provided order num invalid")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	userInfo, err := h.getUserInfoByToken(w, r)
	if err != nil {
		h.log.Errorf("getUserInfoByToken: failed, %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user, err := h.storage.GetUserByLogin(userInfo.Login)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Errorf("LoadOrders: failed to get user by login, %s", err)
		return
	}

	existingOrder, err := h.storage.GetOrderByNumber(orderNum)
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
	if existingOrder.Number > 0 && existingOrder.UserID != userInfo.ID {
		h.log.Infof("Provided order num %s already exist", orderNum)
		w.WriteHeader(http.StatusConflict)
		return
	}
	if err = h.storage.LoadOrder(orderNum, user); err != nil {
		switch err {
		case errors.ErrDuplicateValue:
			{
				w.WriteHeader(http.StatusOK)
				return
			}
		default:
			{
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}
	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {

}
