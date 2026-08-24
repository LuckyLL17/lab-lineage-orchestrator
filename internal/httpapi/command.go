package httpapi

import (
	"net/http"

	"github.com/local/lab-lineage-orchestrator/internal/app"
	"github.com/local/lab-lineage-orchestrator/internal/domain"
)

func (r *Router) command(w http.ResponseWriter, req *http.Request) {
	var command app.Command
	if err := decodeJSON(req, &command); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	command.Subject = domain.ID("http:" + string(command.Subject))
	event, err := r.service.Apply(command)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusAccepted, event)
}
