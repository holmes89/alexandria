package journal

import (
	"alexandria/internal/common"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type journalHandler struct {
	repo Repository
}

func MakeJournalHandler(mr *mux.Router, repo Repository) http.Handler {
	r := mr.PathPrefix("/journal").Subrouter()
	h := &journalHandler{
		repo: repo,
	}
	r.HandleFunc("/entry/", h.FindAll).Methods("GET")
	r.HandleFunc("/entry/", h.Create).Methods("POST")

	return r
}

func (h *journalHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	entities, err := h.repo.FindAllEntries()
	if err != nil {
		common.MakeError(w, http.StatusBadRequest, "journal", "Server error", "findall")
		return
	}

	common.EncodeResponse(ctx, w, entities)
}

func (h *journalHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	b, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	entity := Entry{}
	if err := json.Unmarshal(b, &entity); err != nil {
		logrus.WithError(err).Error("unable to unmarshal journal entry")
		common.MakeError(w, http.StatusBadRequest, "journal", "Bad Request", "create")
		return
	}

	entity, err := h.repo.CreateEntry(entity)
	if err != nil {
		common.MakeError(w, http.StatusInternalServerError, "journal", "Server error", "create")
		return
	}

	common.EncodeResponse(ctx, w, entity)
}
