package ideas

import (
	"alexandria/internal/common"
	"alexandria/internal/tags"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type ideaHandler struct {
	repo     Repository
	tagsRepo tags.Repository
}

func MakeIdeaHandler(mr *mux.Router, repo Repository, tagsRepo tags.Repository) http.Handler {
	r := mr.PathPrefix("/idea").Subrouter()
	h := &ideaHandler{
		repo:     repo,
		tagsRepo: tagsRepo,
	}
	r.HandleFunc("/", h.FindAll).Methods("GET")
	r.HandleFunc("/{id}", h.FindByID).Methods("GET")
	r.HandleFunc("/", h.Create).Methods("POST")
	r.HandleFunc("/{id}/resource/", h.AddResource).Methods("POST")
	r.HandleFunc("/{id}/resource/", h.RemoveResource).Methods("DELETE")
	r.HandleFunc("/{id}/tags/", h.AddTag).Methods("POST")
	r.HandleFunc("/{id}/tags/", h.RemoveTag).Methods("DELETE")

	return r
}

func (h *ideaHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	entities, err := h.repo.GetIdeas()
	if err != nil {
		common.MakeError(w, http.StatusBadRequest, "ideas", "Server error", "findall")
		return
	}

	common.EncodeResponse(ctx, w, entities)
}

func (h *ideaHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	id := vars["id"]

	entity, err := h.repo.GetIdeaByID(id)
	if err != nil {
		common.MakeError(w, http.StatusBadRequest, "ideas", "Server error", "find")
		return
	}

	common.EncodeResponse(ctx, w, entity)
}

func (h *ideaHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	b, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	entity := Idea{}
	if err := json.Unmarshal(b, &entity); err != nil {
		logrus.WithError(err).Error("unable to unmarshal idea")
		common.MakeError(w, http.StatusBadRequest, "idea", "Bad Request", "create")
		return
	}

	entity, err := h.repo.CreateIdea(entity)
	if err != nil {
		common.MakeError(w, http.StatusInternalServerError, "idea", "Server error", "create")
		return
	}

	common.EncodeResponse(ctx, w, entity)
}

func (h *ideaHandler) AddResource(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	id := vars["id"]

	b, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	entity := IdeaResource{}
	if err := json.Unmarshal(b, &entity); err != nil {
		logrus.WithError(err).Error("unable to unmarshal idea")
		common.MakeError(w, http.StatusBadRequest, "idea", "Bad Request", "addresource")
		return
	}

	if err := h.repo.AddIdeaResource(IdeaResource{
		ID:         id,
		ResourceID: entity.ResourceID,
	}); err != nil {
		logrus.WithError(err).Error("unable to unmarshal idea")
		common.MakeError(w, http.StatusBadRequest, "idea", "Bad Request", "addresource")
		return
	}

	i, err := h.repo.GetIdeaByID(id)
	if err != nil {
		common.MakeError(w, http.StatusBadRequest, "ideas", "Server error", "find")
		return
	}
	common.EncodeResponse(ctx, w, i)
}

func (h *ideaHandler) RemoveResource(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	id := vars["id"]

	b, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	entity := IdeaResource{}
	if err := json.Unmarshal(b, &entity); err != nil {
		logrus.WithError(err).Error("unable to unmarshal idea")
		common.MakeError(w, http.StatusBadRequest, "idea", "Bad Request", "removeresource")
		return
	}

	if err := h.repo.RemoveIdeaResource(IdeaResource{
		ID:         id,
		ResourceID: entity.ResourceID,
	}); err != nil {
		logrus.WithError(err).Error("unable to unmarshal idea")
		common.MakeError(w, http.StatusBadRequest, "idea", "Bad Request", "removeresource")
		return
	}

	i, err := h.repo.GetIdeaByID(id)
	if err != nil {
		common.MakeError(w, http.StatusBadRequest, "ideas", "Server error", "find")
		return
	}
	common.EncodeResponse(ctx, w, i)
}

func (h *ideaHandler) AddTag(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	b, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	req := tagRequest{}
	if err := json.Unmarshal(b, &req); err != nil {
		logrus.WithError(err).Error("unable to unmarshal link tag")
		common.MakeError(w, http.StatusBadRequest, "ideas", "Bad Request", "addTag")
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.tagsRepo.AddResourceTag(id, common.IdeaResource, req.Tag); err != nil {
		common.MakeError(w, http.StatusInternalServerError, "ideas", "Server error", "addTag")
		return
	}

	entity, err := h.repo.GetIdeaByID(id)
	if err != nil {
		common.MakeError(w, http.StatusInternalServerError, "ideas", "Server error", "addTag")
		return
	}

	common.EncodeResponse(ctx, w, entity)
}

func (h *ideaHandler) RemoveTag(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	b, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	req := tagRequest{}
	if err := json.Unmarshal(b, &req); err != nil {
		logrus.WithError(err).Error("unable to unmarshal link tag")
		common.MakeError(w, http.StatusBadRequest, "ideas", "Bad Request", "removeTag")
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.tagsRepo.RemoveResourceTag(id, req.Tag); err != nil {
		common.MakeError(w, http.StatusInternalServerError, "ideas", "Server error", "removeTag")
		return
	}

	entity, err := h.repo.GetIdeaByID(id)
	if err != nil {
		common.MakeError(w, http.StatusInternalServerError, "ideas", "Server error", "removeTag")
		return
	}

	common.EncodeResponse(ctx, w, entity)
}

type tagRequest struct {
	Tag string `json:"tag"`
}
