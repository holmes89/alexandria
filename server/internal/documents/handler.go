package documents

import (
	"alexandria/internal/common"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func MakeDocumentHandler(mr *mux.Router, service DocumentService) http.Handler {
	r := mr.PathPrefix("/documents").Subrouter()

	h := &documentHandler{
		service: service,
	}

	r.HandleFunc("/{id}", h.FindByID).Methods("GET")
	r.HandleFunc("/{id}", h.UpdateFields).Methods("PATCH")
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
		common.MakeError(w, http.StatusBadRequest, "document", "Missing Id", "findbyid")
		return
	}

	entity, err := h.service.GetByID(ctx, id)

	if err != nil {
		common.MakeError(w, http.StatusInternalServerError, "document", "Server Error", "findbyid")
		return
	}

	common.EncodeResponse(r.Context(), w, entity)
}

func (h *documentHandler) UpdateFields(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	b, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	req := Document{}
	if err := json.Unmarshal(b, &req); err != nil {
		logrus.WithError(err).Error("unable to unmarshal link tag")
		common.MakeError(w, http.StatusBadRequest, "document", "Bad Request", "updateFields")
		return
	}

	vars := mux.Vars(r)

	id, ok := vars["id"]

	if !ok {
		common.MakeError(w, http.StatusBadRequest, "document", "Missing Id", "updateFields")
		return
	}

	entity, err := h.service.UpdateFields(ctx, id, req)

	if err != nil {
		common.MakeError(w, http.StatusInternalServerError, "document", "Server Error", "updateFields")
		return
	}

	common.EncodeResponse(r.Context(), w, entity)
}

func (h *documentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, ok := vars["id"]

	if !ok {
		common.MakeError(w, http.StatusBadRequest, "document", "Missing Id", "delete")
		return
	}

	if err := h.service.Delete(ctx, id); err != nil {
		common.MakeError(w, http.StatusInternalServerError, "document", "Server Error", "delete")
		return
	}

	common.EncodeResponse(r.Context(), w, map[string]string{"status": "success"})
}

func (h *documentHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	entity, err := h.service.GetAll(ctx, nil)

	if err != nil {
		common.MakeError(w, http.StatusInternalServerError, "document", "Server Error", "findall")
		return
	}

	common.EncodeResponse(r.Context(), w, entity)
}

func (h *documentHandler) Scan(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := h.service.Scan(ctx)

	if err != nil {
		common.MakeError(w, http.StatusInternalServerError, "document", "Server Error", "scan")
		return
	}

	common.EncodeResponse(r.Context(), w, map[string]string{"status": "success"})
}

func (h *bookHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		common.MakeError(w, http.StatusBadRequest, "book", "Unable to parse form", "create")
		return
	}
	if file == nil {
		common.MakeError(w, http.StatusBadRequest, "book", "File missing from form", "create")
		return
	}
	defer file.Close()

	displayName, ok := r.MultipartForm.Value["name"]
	if !ok {
		common.MakeError(w, http.StatusBadRequest, "book", "Name missing from form", "create")
		return
	}

	book := &Document{
		DisplayName: displayName[0],
		Name:        fileHeader.Filename,
		Type:        "book",
	}

	if err := h.service.Add(ctx, file, book); err != nil {
		common.MakeError(w, http.StatusInternalServerError, "book", err.Error(), "add")
		return
	}

	w.WriteHeader(http.StatusCreated)
	common.EncodeResponse(r.Context(), w, book)
}

func (h *bookHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	entity, err := h.service.GetAll(ctx)

	if err != nil {
		common.MakeError(w, http.StatusInternalServerError, "book", "Server Error", "findall")
		return
	}

	common.EncodeResponse(r.Context(), w, entity)
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

	book := &Document{
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

	entity, err := h.service.GetAll(ctx)

	if err != nil {
		common.MakeError(w, http.StatusInternalServerError, "paper", "Server Error", "findall")
		return
	}

	common.EncodeResponse(r.Context(), w, entity)
}
