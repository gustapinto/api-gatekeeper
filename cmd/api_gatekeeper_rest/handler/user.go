package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gustapinto/api-gatekeeper/cmd/api_gatekeeper_rest/dto/response"
	"github.com/gustapinto/api-gatekeeper/internal/model"
	"github.com/gustapinto/api-gatekeeper/internal/service"
	httputil "github.com/gustapinto/api-gatekeeper/pkg/http_util"
)

type User struct {
	userService *service.User
	jwtService  *service.JWT
}

func NewUser(userService *service.User, jwtService *service.JWT) User {
	return User{
		userService: userService,
		jwtService:  jwtService,
	}
}

func (u User) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateUserParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteBadRequest(w, errors.New("failed to parse request body"))
		return
	}

	user, err := u.userService.Create(req)
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

	user, err := u.userService.Update(req)
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

	if err := u.userService.Delete(userId); err != nil {
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

	user, err := u.userService.GetByID(userId)
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
	user, err := u.userService.GetAll()
	if err != nil {
		httputil.WriteUnprocessableEntity(w, err)
		return
	}

	httputil.WriteOk(w, user)
}

func (u User) Login(w http.ResponseWriter, r *http.Request) {
	username, password, err := httputil.ParseBasicAuthorizationToken(r.Header.Get("Authorization"))
	if err != nil {
		if strings.Contains(err.Error(), "badparams:") {
			httputil.WriteBadRequest(w, err)
			return
		}

		httputil.WriteUnauthorized(w)
		return
	}

	user, err := u.userService.Login(username, password)
	if err != nil {
		httputil.WriteBadRequest(w, err)
		return
	}

	tokenType := r.Header.Get("X-Token-Type")
	if strings.ToLower(strings.TrimSpace(tokenType)) == "jwt" {
		token, err := u.jwtService.GenerateToken(user)
		if err != nil {
			httputil.WriteBadRequest(w, err)
			return
		}

		httputil.WriteOk(w, response.JWTTokenresponse{
			Token: token,
		})
		return
	}
}
