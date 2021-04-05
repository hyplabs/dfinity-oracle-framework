package models

import "time"

// Config is the configuration for the oracle to be made
type Config struct {
	CanisterName   string
	UpdateInterval time.Duration
}
