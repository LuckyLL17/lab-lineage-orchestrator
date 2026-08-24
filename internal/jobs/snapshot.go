package jobs

import (
	"context"
	"time"
)

func (r *Runner) snapshotLoop(ctx context.Context) {
	defer r.wg.Done()
	ticker := time.NewTicker(2 * r.tick)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			snapshot := r.service.Snapshot()
			r.log.Info("snapshot captured", "digest", snapshot.Digest, "events", snapshot.Events)
		}
	}
}
