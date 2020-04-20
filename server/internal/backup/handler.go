package backup

import (
	"alexandria/internal/common"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

type backupHandler struct {
	service Service
}

func MakeBackupHandler(r *mux.Router, service Service) http.Handler {

	h := &backupHandler{
		service: service,
	}

	r.HandleFunc("/backup/", h.Backup).Methods("POST")
	r.HandleFunc("/restore/{id}", h.Restore).Methods("POST")

	return r
}

func (h *backupHandler) Backup(w http.ResponseWriter, r *http.Request) {

	if err := h.service.Backup(); err != nil {
		common.MakeError(w, http.StatusInternalServerError, "backup", "unable to backup", "backup")
		return
	}

	common.EncodeResponse(r.Context(), w, map[string]string{"status": "success"})
}

func (h *backupHandler) Restore(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, ok := vars["id"]

	if !ok {
		common.MakeError(w, http.StatusBadRequest, "backup", "Missing Id", "restore")
		return
	}

	v := r.URL.Query()
	restoreType := ParseRestore(v.Get("type"))

	if restoreType == RestoreUnknown {
		common.MakeError(w, http.StatusBadRequest, "backup", "unsupported backup", "restore")
		return
	}

	logrus.WithFields(logrus.Fields{"id": id, "type": restoreType}).Info("attempting to restore database")

	if err := h.service.Restore(id, restoreType); err != nil {
		common.MakeError(w, http.StatusInternalServerError, "backup", "unable to restore", "restore")
		return
	}

	common.EncodeResponse(r.Context(), w, map[string]string{"status": "success"})
}
