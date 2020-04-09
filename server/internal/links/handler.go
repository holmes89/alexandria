package links

import (
	"alexandria/internal/common"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type linkHandler struct {
	service Service
}

func MakeLinksHandler(mr *mux.Router, service Service) http.Handler {
	r := mr.PathPrefix("/links").Subrouter()
	h := &linkHandler{
		service: service,
	}
	r.HandleFunc("/", h.FindAll).Methods("GET")
	r.HandleFunc("/", h.Create).Methods("POST")

	return r
}

func (h *linkHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	entities, err := h.service.FindAll()
	if err != nil {
		common.MakeError(w, http.StatusBadRequest, "links", "Server error", "findall")
		return
	}

	common.EncodeResponse(ctx, w, entities)
}

func (h *linkHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	b, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	entity := Link{}
	if err := json.Unmarshal(b, &entity); err != nil {
		logrus.WithError(err).Error("unable to unmarshal link")
		common.MakeError(w, http.StatusBadRequest, "links", "Bad Request", "create")
		return
	}

	entity, err := h.service.Create(entity)
	if err != nil {
		common.MakeError(w, http.StatusInternalServerError, "links", "Server error", "create")
		return
	}

	common.EncodeResponse(ctx, w, entity)
}
