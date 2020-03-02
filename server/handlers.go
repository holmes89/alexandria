package main

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func MakeDocumentHandler(mr *mux.Router, service DocumentService) http.Handler {
	r := mr.PathPrefix("/documents").Subrouter()

	h := &documentHandler{
		service: service,
	}

	r.HandleFunc("/{id}", h.FindByID).Methods("GET")
	r.HandleFunc("/{id}", h.Delete).Methods("DELETE")
	r.HandleFunc("/scan", h.Scan).Methods("PUT")
	r.HandleFunc("/", h.FindAll).Methods("GET")

	return r
}

type documentHandler struct {
	service DocumentService
}

func MakeBookHandler(mr *mux.Router, service BookService) http.Handler {
	r := mr.PathPrefix("/books").Subrouter()

	h := &bookHandler{
		service: service,
	}

	r.HandleFunc("/", h.FindAll).Methods("GET")
	r.HandleFunc("/", h.Create).Methods("POST")

	return r
}

type bookHandler struct {
	service BookService
}

func MakePaperHandler(mr *mux.Router, service PaperService) http.Handler {
	r := mr.PathPrefix("/papers").Subrouter()

	h := &paperHandler{
		service: service,
	}

	r.HandleFunc("/", h.FindAll).Methods("GET")
	r.HandleFunc("/", h.Create).Methods("POST")

	return r
}

type paperHandler struct {
	service PaperService
}

func (h *documentHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, ok := vars["id"]

	if !ok {
		makeError(w, http.StatusBadRequest, "document", "Missing Id", "findbyid")
		return
	}

	entity, err := h.service.GetByID(ctx, id)

	if err != nil {
		makeError(w, http.StatusInternalServerError, "document", "Server Error", "findbyid")
		return
	}

	encodeResponse(r.Context(), w, entity)
}

func (h *documentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, ok := vars["id"]

	if !ok {
		makeError(w, http.StatusBadRequest, "document", "Missing Id", "delete")
		return
	}

	if err := h.service.Delete(ctx, id); err != nil {
		makeError(w, http.StatusInternalServerError, "document", "Server Error", "delete")
		return
	}

	encodeResponse(r.Context(), w, map[string]string{"status": "success"})
}

func (h *documentHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	entity, err := h.service.GetAll(ctx, nil)

	if err != nil {
		makeError(w, http.StatusInternalServerError, "document", "Server Error", "findall")
		return
	}

	encodeResponse(r.Context(), w, entity)
}

func (h *documentHandler) Scan(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := h.service.Scan(ctx)

	if err != nil {
		makeError(w, http.StatusInternalServerError, "document", "Server Error", "scan")
		return
	}

	encodeResponse(r.Context(), w, map[string]string{"status": "success"})
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

	book := &Document{
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

func (h *bookHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	entity, err := h.service.GetAll(ctx)

	if err != nil {
		makeError(w, http.StatusInternalServerError, "book", "Server Error", "findall")
		return
	}

	encodeResponse(r.Context(), w, entity)
}

func (h *paperHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		makeError(w, http.StatusBadRequest, "paper", "Unable to parse form", "create")
		return
	}
	if file == nil {
		makeError(w, http.StatusBadRequest, "paper", "File missing from form", "create")
		return
	}
	defer file.Close()

	displayName, ok := r.MultipartForm.Value["name"]
	if !ok {
		makeError(w, http.StatusBadRequest, "paper", "Name missing from form", "create")
		return
	}

	book := &Document{
		DisplayName: displayName[0],
		Name:        fileHeader.Filename,
		Type:        "paper",
	}

	if err := h.service.Add(ctx, file, book); err != nil {
		makeError(w, http.StatusInternalServerError, "book", err.Error(), "add")
		return
	}

	w.WriteHeader(http.StatusCreated)
	encodeResponse(r.Context(), w, book)
}

func (h *paperHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	entity, err := h.service.GetAll(ctx)

	if err != nil {
		makeError(w, http.StatusInternalServerError, "paper", "Server Error", "findall")
		return
	}

	encodeResponse(r.Context(), w, entity)
}

type loginHandler struct {
	service UserService
}

func MakeLoginHandler(mr *mux.Router, service UserService) http.Handler {
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
		makeError(w, http.StatusUnauthorized, "login", "missing auth header", "login")
		return
	}

	token, err := h.service.Authenticate(ctx, username, password)
	if err != nil {
		logrus.WithError(err).Error("failed to login")
		makeError(w, http.StatusUnauthorized, "login", "invalid login", "login")
		return
	}

	encodeResponse(r.Context(), w, token)
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
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return enc.Encode(response)
}
