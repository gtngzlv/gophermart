package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gtngzlv/gophermart/internal/auth"
	"github.com/gtngzlv/gophermart/internal/errors"
	"github.com/gtngzlv/gophermart/internal/model"
	"github.com/gtngzlv/gophermart/internal/utils"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var user model.User
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.Errorf("UserID: failed to read from body %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		h.log.Errorf("UserID: failed to unmarshal body %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userDB, err := h.repo.GetUserByLogin(user.Login)
	if !utils.CheckHashAndPassword(userDB.Password, user.Password) {
		h.log.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil && err != errors.ErrNoDBResult {
		h.log.Errorf("failed to login, %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = auth.GenerateCookie(w, user.ID)
	if err != nil {
		h.log.Errorf("Failed to generate cookie, %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var user model.User
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.Errorf("Register: failed to read from body %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		h.log.Errorf("Register: failed to unmarshal body %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Посмотрим, что юзера с таким логином нет
	userInDB, err := h.repo.GetUserByLogin(user.Login)
	switch {
	case err == nil && userInDB.Login != "":
		h.log.Infof("Register: user with provided login %s exists", user.Login)
		w.WriteHeader(http.StatusConflict)
		return
	case err == errors.ErrNoDBResult:
		cryptedPsw, err := utils.HashString(user.Password)
		if err != nil {
			h.log.Errorf("Register: failed to encrypt password")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		userID, err := h.repo.Register(user.Login, cryptedPsw)
		if err != nil {
			h.log.Errorf("Register: failed while registering in storage")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = auth.GenerateCookie(w, userID)
		if err != nil {
			h.log.Errorf("Failed to generate cookie, %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}
