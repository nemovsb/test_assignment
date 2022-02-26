package storage

import (
	"time"

	"gorm.io/gorm"
)

type Sites struct {
	gorm.Model
	Name        string
	LoadingTime time.Time
}
