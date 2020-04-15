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
	r.HandleFunc("/{id}", h.FindByID).Methods("GET")
	r.HandleFunc("/", h.Create).Methods("POST")
	r.HandleFunc("/{id}/tags/", h.AddTag).Methods("POST")
	r.HandleFunc("/{id}/tags/", h.RemoveTag).Methods("DELETE")

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

func (h *linkHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	id := vars["id"]

	entity, err := h.service.FindByID(id)
	if err != nil {
		common.MakeError(w, http.StatusBadRequest, "links", "Server error", "find")
		return
	}

	common.EncodeResponse(ctx, w, entity)
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

func (h *linkHandler) AddTag(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	b, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	req := tagRequest{}
	if err := json.Unmarshal(b, &req); err != nil {
		logrus.WithError(err).Error("unable to unmarshal link tag")
		common.MakeError(w, http.StatusBadRequest, "links", "Bad Request", "addTag")
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.service.AddTag(id, req.Tag); err != nil {
		common.MakeError(w, http.StatusInternalServerError, "links", "Server error", "addTag")
		return
	}

	entity, err := h.service.FindByID(id)
	if err != nil {
		common.MakeError(w, http.StatusInternalServerError, "links", "Server error", "addTag")
		return
	}

	common.EncodeResponse(ctx, w, entity)
}

func (h *linkHandler) RemoveTag(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	b, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	req := tagRequest{}
	if err := json.Unmarshal(b, &req); err != nil {
		logrus.WithError(err).Error("unable to unmarshal link tag")
		common.MakeError(w, http.StatusBadRequest, "links", "Bad Request", "removeTag")
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.service.RemoveTag(id, req.Tag); err != nil {
		common.MakeError(w, http.StatusInternalServerError, "links", "Server error", "removeTag")
		return
	}

	entity, err := h.service.FindByID(id)
	if err != nil {
		common.MakeError(w, http.StatusInternalServerError, "links", "Server error", "removeTag")
		return
	}

	common.EncodeResponse(ctx, w, entity)
}

type tagRequest struct {
	Tag string `json"tag"`
}
