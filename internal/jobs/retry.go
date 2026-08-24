package jobs

import (
	"context"
	"time"
)

func (r *Runner) retryLoop(ctx context.Context) {
	defer r.wg.Done()
	ticker := time.NewTicker(r.tick + 3*time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			event := r.service.EmitOperationalEvent("retry-scan", "system", "retry queue inspected")
			r.log.Info("retry scan completed", "event", event.ID)
		}
	}
}
