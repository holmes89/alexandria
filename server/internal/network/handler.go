package network

import (
	"alexandria/internal/common"
	"github.com/gorilla/mux"
	"net/http"
)

type networkHandler struct {
	service Service
}

func MakeNetworkHandler(mr *mux.Router, service Service) http.Handler {
	r := mr.PathPrefix("/network").Subrouter()
	h := &networkHandler{
		service: service,
	}

	r.HandleFunc("/", h.GetNetwork).Methods("GET")

	return r
}

func (h *networkHandler) GetNetwork(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	entity, err := h.service.GetNetwork()
	if err != nil {
		common.MakeError(w, http.StatusBadRequest, "network", "Server error", "get")
		return
	}

	common.EncodeResponse(ctx, w, entity)
}
