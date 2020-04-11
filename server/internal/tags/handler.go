package tags

import (
	"alexandria/internal/common"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type tagHandler struct {
	repo Repository
}

func MakeLinksHandler(mr *mux.Router, repo Repository) http.Handler {
	r := mr.PathPrefix("/tags").Subrouter()
	h := &tagHandler{
		repo: repo,
	}
	r.HandleFunc("/", h.FindAll).Methods("GET")
	r.HandleFunc("/", h.Create).Methods("POST")

	return r
}

func (h *tagHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	entities, err := h.repo.FindAllTags()
	if err != nil {
		common.MakeError(w, http.StatusBadRequest, "tags", "Server error", "findall")
		return
	}

	common.EncodeResponse(ctx, w, entities)
}

func (h *tagHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	b, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	entity := Tag{}
	if err := json.Unmarshal(b, &entity); err != nil {
		logrus.WithError(err).Error("unable to unmarshal tag")
		common.MakeError(w, http.StatusBadRequest, "tags", "Bad Request", "create")
		return
	}

	entity, err := h.repo.CreateTag(entity)
	if err != nil {
		common.MakeError(w, http.StatusInternalServerError, "tags", "Server error", "create")
		return
	}

	common.EncodeResponse(ctx, w, entity)
}