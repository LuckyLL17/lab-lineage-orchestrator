package jobs

import (
	"context"
	"time"
)

func (r *Runner) retentionLoop(ctx context.Context) {
	defer r.wg.Done()
	ticker := time.NewTicker(12 * r.tick)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			cutoff := time.Now().UTC().Add(-7 * 24 * time.Hour)
			removed := r.service.RetainSince(cutoff)
			r.log.Info("retention pass completed", "removed", removed)
		}
	}
}
