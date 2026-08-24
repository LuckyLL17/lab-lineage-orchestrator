package platform

import (
	"fmt"
	"sync/atomic"
	"time"
)

var sequence uint64

func NewID(prefix string) string {
	value := atomic.AddUint64(&sequence, 1)
	return fmt.Sprintf("%s-%d-%d", prefix, time.Now().UnixNano(), value)
}

func CorrelationID() string {
	return NewID("corr")
}
