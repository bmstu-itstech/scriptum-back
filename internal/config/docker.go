package config

import "time"

type Docker struct {
	ImagePrefix   string
	RunnerTimeout time.Duration
}
