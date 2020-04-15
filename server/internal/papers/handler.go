package papers

import (
	"alexandria/internal/common"
	"alexandria/internal/documents"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func MakePaperHandler(mr *mux.Router, service PaperService) http.Handler {
	r := mr.PathPrefix("/papers").Subrouter()

	h := &paperHandler{
		service: service,
	}

	r.HandleFunc("/", h.FindAll).Methods("GET")
	r.HandleFunc("/{id}", h.FindByID).Methods("GET")
	r.HandleFunc("/", h.Create).Methods("POST")
	r.HandleFunc("/{id}/tags/", h.AddTag).Methods("POST")
	r.HandleFunc("/{id}/tags/", h.RemoveTag).Methods("DELETE")

	return r
}

type paperHandler struct {
	service PaperService
}

func (h *paperHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		common.MakeError(w, http.StatusBadRequest, "paper", "Unable to parse form", "create")
		return
	}
	if file == nil {
		common.MakeError(w, http.StatusBadRequest, "paper", "File missing from form", "create")
		return
	}
	defer file.Close()

	displayName, ok := r.MultipartForm.Value["name"]
	if !ok {
		common.MakeError(w, http.StatusBadRequest, "paper", "Name missing from form", "create")
		return
	}

	book := &documents.Document{
		DisplayName: displayName[0],
		Name:        fileHeader.Filename,
		Type:        "paper",
	}

	if err := h.service.Add(ctx, file, book); err != nil {
		common.MakeError(w, http.StatusInternalServerError, "book", err.Error(), "add")
		return
	}

	w.WriteHeader(http.StatusCreated)
	common.EncodeResponse(r.Context(), w, book)
}

func (h *paperHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	entity, err := h.service.FindAll(ctx)

	if err != nil {
		common.MakeError(w, http.StatusInternalServerError, "paper", "Server Error", "findall")
		return
	}

	common.EncodeResponse(r.Context(), w, entity)
}

func (h *paperHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, ok := vars["id"]

	if !ok {
		common.MakeError(w, http.StatusBadRequest, "document", "Missing Id", "findbyid")
		return
	}

	entity, err := h.service.FindByID(ctx, id)

	if err != nil {
		common.MakeError(w, http.StatusInternalServerError, "document", "Server Error", "findbyid")
		return
	}

	common.EncodeResponse(r.Context(), w, entity)
}

func (h *paperHandler) AddTag(w http.ResponseWriter, r *http.Request) {
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

	entity, err := h.service.FindByID(ctx, id)
	if err != nil {
		common.MakeError(w, http.StatusInternalServerError, "links", "Server error", "addTag")
		return
	}

	common.EncodeResponse(ctx, w, entity)
}

func (h *paperHandler) RemoveTag(w http.ResponseWriter, r *http.Request) {
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

	entity, err := h.service.FindByID(ctx, id)
	if err != nil {
		common.MakeError(w, http.StatusInternalServerError, "links", "Server error", "removeTag")
		return
	}

	common.EncodeResponse(ctx, w, entity)
}

type tagRequest struct {
	Tag string `json"tag"`
}
