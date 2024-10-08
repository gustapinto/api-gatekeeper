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
