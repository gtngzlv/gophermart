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
	var u model.User
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.Errorf("Login: failed to read from body %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &u)
	if err != nil {
		h.log.Errorf("Login: failed to unmarshal body %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userDB, err := h.repo.GetUserByLogin(u.Login)
	if !utils.CheckHashAndPassword(userDB.Password, u.Password) {
		h.log.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil && err != errors.ErrNoDBResult {
		h.log.Errorf("failed to login, %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = auth.GenerateCookie(w, u.Login)
	if err != nil {
		h.log.Errorf("Failed to generate cookie, %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var u model.User
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.Errorf("Register: failed to read from body %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &u)
	if err != nil {
		h.log.Errorf("Register: failed to unmarshal body %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Посмотрим, что юзера с таким логином нет
	userInDB, err := h.repo.GetUserByLogin(u.Login)
	switch {
	case err == nil && userInDB.Login != "":
		h.log.Infof("Register: user with provided login %s exists", u.Login)
		w.WriteHeader(http.StatusConflict)
		return
	case err == errors.ErrNoDBResult:
		cryptedPsw, err := utils.HashString(u.Password)
		if err != nil {
			h.log.Errorf("Register: failed to encrypt password")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = h.repo.Register(u.Login, cryptedPsw)
		if err != nil {
			h.log.Errorf("Register: failed while registering in storage")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = auth.GenerateCookie(w, u.Login)
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

func (h *Handler) getUserInfoByToken(w http.ResponseWriter, r *http.Request) (model.User, error) {
	login := auth.GetUserLoginFromToken(w, r)
	user, err := h.repo.GetUserByLogin(login)
	if err != nil {
		return model.User{}, err
	}
	return *user, nil
}
