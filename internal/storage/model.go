package storage

import (
	"time"

	"gorm.io/gorm"
)

type Sites struct {
	gorm.Model
	Name        string
	LoadingTime time.Duration
}

type Report struct {
	Name        string        `json:"name"`
	AvgDuration time.Duration `json:"avg_duration"`
}
