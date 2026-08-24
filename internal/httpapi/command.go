package httpapi

import (
	"net/http"

	"github.com/local/lab-lineage-orchestrator/internal/app"
)

func (r *Router) command(w http.ResponseWriter, req *http.Request) {
	var command app.Command
	headerActor := req.Header.Get("X-Actor")
	if headerActor != "" {
		command.Actor = headerActor
	}
	if err := decodeJSON(req, &command); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	event, err := r.service.Apply(command)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusAccepted, event)
}
