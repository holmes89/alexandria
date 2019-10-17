package main

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func MakeBookHandler(mr *mux.Router, service BookService) http.Handler {
	r := mr.PathPrefix("books").Subrouter()

	h := &bookHandler{
		service: service,
	}

	r.HandleFunc("/{id}", h.FindByID).Methods("GET")
	r.HandleFunc("/", h.FindAll).Methods("GET")
	r.HandleFunc("/", h.Create).Methods("POST")

	return r
}

type bookHandler struct {
	service BookService
}

func (h *bookHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, ok := vars["id"]

	if !ok {
		makeError(w, http.StatusBadRequest, "book", "Missing Id", "findbyid")
		return
	}

	entity, err := h.service.GetByID(ctx, id)

	if err != nil {
		makeError(w, http.StatusInternalServerError, "book", "Server Error", "findbyid")
		return
	}

	encodeResponse(r.Context(), w, entity)
}

func (h *bookHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	entity, err := h.service.GetAll(ctx)

	if err != nil {
		makeError(w, http.StatusInternalServerError, "book", "Server Error", "findall")
		return
	}

	encodeResponse(r.Context(), w, entity)
}

func (h *bookHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		makeError(w, http.StatusBadRequest, "book", "Unable to parse form", "create")
		return
	}
	if file == nil {
		makeError(w, http.StatusBadRequest, "book", "File missing from form", "create")
		return
	}
	defer file.Close()

	displayName, ok := r.MultipartForm.Value["name"]
	if !ok {
		makeError(w, http.StatusBadRequest, "book", "Name missing from form", "create")
		return
	}

	book := &Book{
		DisplayName: displayName[0],
		Name:        fileHeader.Filename,
		Type:        "book",
	}

	if err := h.service.Add(ctx, file, book); err != nil {
		makeError(w, http.StatusInternalServerError, "book", err.Error(), "add")
		return
	}

	w.WriteHeader(http.StatusCreated)
	encodeResponse(r.Context(), w, book)
}

func makeError(w http.ResponseWriter, code int, domain string, message string, method string) {
	logrus.WithFields(
		logrus.Fields{
			"type":   code,
			"domain": domain,
			"method": method,
		}).Error(strings.ToLower(message))
	http.Error(w, message, code)
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
