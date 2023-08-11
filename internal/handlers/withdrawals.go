package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/ShiraazMoollatjie/goluhn"

	customErr "github.com/gtngzlv/gophermart/internal/errors"
	"github.com/gtngzlv/gophermart/internal/model"
)

func (h *Handler) WithdrawLoyalty(w http.ResponseWriter, r *http.Request) {
	var withdrawRequest model.WithdrawBalanceRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.Error("WithdrawLoyalty: failed while read body", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	buf := bytes.NewBuffer(body)
	err = json.NewDecoder(buf).Decode(&withdrawRequest)
	if err != nil {
		h.log.Error("WithdrawLoyalty: failed while decode", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	if err = goluhn.Validate(withdrawRequest.Order); err != nil {
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
	balance, err := h.storage.GetBalance(userInfo.ID)
	if err != nil {
		h.log.Errorf("GetBalance: failed, %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if balance.Current < withdrawRequest.Sum {
		h.log.Info("Withdraw: balance not enough")
		w.WriteHeader(http.StatusPaymentRequired)
		return
	}
	// если нет такого заказа, мы его создаем
	_, err = h.storage.GetOrderByNumber(withdrawRequest.Order)
	if err == customErr.ErrNoDBResult {
		h.log.Infof("WithdrawLoyalty: provided order with num %s not exist, creating", withdrawRequest.Order)
		h.storage.LoadOrder(withdrawRequest.Order, userInfo)
	}

	// cписываем
	err = h.storage.WithdrawLoyalty(withdrawRequest, userInfo.ID, withdrawRequest.Order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.log.Info("Order received to withdraw loyalty", withdrawRequest.Order)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userInfo, err := h.getUserInfoByToken(w, r)
	if err != nil {
		h.log.Errorf("getUserInfoByToken: failed, %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	withdrawals, err := h.storage.GetWithdrawals(userInfo.ID)
	if err != nil {
		switch err {
		case customErr.ErrNoDBResult:
			{
				h.log.Info("No withdrawals for provided user", userInfo.ID)
				w.WriteHeader(http.StatusNoContent)
				return
			}
		default:
			{
				h.log.Error("Failed to get withdrawals", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}

	h.log.Infof("%v user withdrawals are %v", userInfo.ID, withdrawals)
	resp, err := json.Marshal(withdrawals)
	if err != nil {
		h.log.Error("Failed to marshal get withdrawals", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(resp)
	w.WriteHeader(http.StatusOK)
}
