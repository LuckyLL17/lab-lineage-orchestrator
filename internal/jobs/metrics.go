package jobs

import (
	"context"
	"time"
)

func (r *Runner) metricsLoop(ctx context.Context) {
	defer r.wg.Done()
	ticker := time.NewTicker(4 * r.tick)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			metrics := r.service.Metrics()
			r.log.Info("metrics sampled", "recent", metrics.Recent, "counters", len(metrics.Counters))
		}
	}
}
