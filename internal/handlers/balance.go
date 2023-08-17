package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gtngzlv/gophermart/internal/errors"
)

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	userInfo, err := h.getUserInfoByToken(w, r)
	if err != nil {
		h.log.Errorf("getUserInfoByToken: failed, %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	balance, err := h.repo.GetBalance(userInfo.ID)
	if err != nil {
		switch err {
		case errors.ErrNoDBResult:
			{
				h.log.Info("GetBalance: no balance for provided userID", userInfo.ID)
				w.WriteHeader(http.StatusNoContent)
				return
			}
		default:
			h.log.Error("GetBalance: error while select to db", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	resp, err := json.Marshal(balance)
	if err != nil {
		h.log.Errorf("GetBalance: failed to marshal resp %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}
