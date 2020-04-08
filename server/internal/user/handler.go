package user

import (
	"alexandria/internal/common"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

type loginHandler struct {
	service Service
}

func MakeLoginHandler(mr *mux.Router, service Service) http.Handler {
	h := &loginHandler{
		service: service,
	}
	mr.HandleFunc("/auth/", h.Login).Methods("GET")

	return mr
}

func (h *loginHandler) Login(w http.ResponseWriter, r *http.Request) {
	logrus.Info("here")
	ctx := r.Context()

	username, password, ok := r.BasicAuth()
	if !ok {
		logrus.Warn("missing auth header")
		common.MakeError(w, http.StatusUnauthorized, "login", "missing auth header", "login")
		return
	}

	token, err := h.service.Authenticate(ctx, username, password)
	if err != nil {
		logrus.WithError(err).Error("failed to login")
		common.MakeError(w, http.StatusUnauthorized, "login", "invalid login", "login")
		return
	}

	common.EncodeResponse(r.Context(), w, token)
}

