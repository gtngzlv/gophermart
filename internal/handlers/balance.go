package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gtngzlv/gophermart/internal/auth"
	"github.com/gtngzlv/gophermart/internal/errors"
)

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromToken(w, r)

	balance, err := h.repo.GetBalance(userID)
	if err != nil {
		switch err {
		case errors.ErrNoDBResult:
			{
				h.log.Info("GetBalance: no balance for provided userID", userID)
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
