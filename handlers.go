package main

import (
	"github.com/gorilla/mux"
	"net/http"
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

}

func (h *bookHandler) FindAll(w http.ResponseWriter, r *http.Request) {

}

func (h *bookHandler) Create(w http.ResponseWriter, r *http.Request) {

}
