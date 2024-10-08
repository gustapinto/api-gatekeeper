package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gustapinto/api-gatekeeper/internal/model"
	"github.com/gustapinto/api-gatekeeper/internal/service"
	httputil "github.com/gustapinto/api-gatekeeper/pkg/http_util"
)

type User struct {
	Service service.User
}

func (u User) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateUserParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteBadRequest(w, errors.New("failed to parse request body"))
		return
	}

	user, err := u.Service.Create(req)
	if err != nil {
		if strings.Contains(err.Error(), "badparams:") {
			httputil.WriteBadRequest(w, err)
			return
		}

		httputil.WriteUnprocessableEntity(w, err)
		return
	}

	httputil.WriteCreated(w, user)
}

func (u User) Update(w http.ResponseWriter, r *http.Request) {
	var req model.UpdateUserParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteBadRequest(w, errors.New("failed to parse request body"))
		return
	}

	req.ID = r.PathValue("userId")

	user, err := u.Service.Update(req)
	if err != nil {
		if strings.Contains(err.Error(), "badparams:") {
			httputil.WriteBadRequest(w, err)
			return
		}

		httputil.WriteUnprocessableEntity(w, err)
		return
	}

	httputil.WriteOk(w, user)
}

func (u User) Delete(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("userId")

	if err := u.Service.Delete(userId); err != nil {
		if strings.Contains(err.Error(), "badparams:") {
			httputil.WriteBadRequest(w, err)
			return
		}

		httputil.WriteUnprocessableEntity(w, err)
		return
	}

	httputil.WriteNoContent(w)
}

func (u User) GetByID(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("userId")

	user, err := u.Service.GetByID(userId)
	if err != nil {
		if strings.Contains(err.Error(), "badparams:") {
			httputil.WriteBadRequest(w, err)
			return
		}

		httputil.WriteUnprocessableEntity(w, err)
		return
	}

	httputil.WriteOk(w, user)
}

func (u User) GetAll(w http.ResponseWriter, r *http.Request) {
	user, err := u.Service.GetAll()
	if err != nil {
		httputil.WriteUnprocessableEntity(w, err)
		return
	}

	httputil.WriteOk(w, user)
}
