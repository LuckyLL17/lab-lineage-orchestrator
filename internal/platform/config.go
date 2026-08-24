package platform

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Addr     string
	Tick     time.Duration
	DataDir  string
	MaxBatch int
	Shutdown time.Duration
}

func LoadConfig() Config {
	tick := 15 * time.Second
	if raw := os.Getenv("APP_TICK_SECONDS"); raw != "" {
		if value, err := strconv.Atoi(raw); err == nil && value > 0 {
			tick = time.Duration(value) * time.Second
		}
	}
	addr := os.Getenv("APP_ADDR")
	if addr == "" {
		addr = ":8090"
	}
	dataDir := os.Getenv("APP_DATA_DIR")
	if dataDir == "" {
		dataDir = "./data"
	}
	return Config{
		Addr:     addr,
		Tick:     tick,
		DataDir:  dataDir,
		MaxBatch: 250,
		Shutdown: 8 * time.Second,
	}
}
