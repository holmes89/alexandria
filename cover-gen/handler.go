package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

func MakeThumnbnailHandler(mr *mux.Router, service *Service) http.Handler {
	r := mr.PathPrefix("/thumbnail").Subrouter()

	h := &thumbnailHandler{
		service: service,
	}
	r.HandleFunc("/", h.Post).Methods("POST")

	return r
}

type thumbnailHandler struct {
	service *Service
}

func (h *thumbnailHandler) Post(w http.ResponseWriter, r *http.Request) {

	dec := json.NewDecoder(r.Body)
	e := &request{}

	if err := dec.Decode(e); err != nil {
		logrus.WithError(err).Error("unable to decode message")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.CreateCover(e.ID, e.Path); err != nil {
		logrus.WithError(err).Error("unable to decode message")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("success"))
}

type request struct {
	ID string `json:"id"`
	Path string `json:"path"`
}
