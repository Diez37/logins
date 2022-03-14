package domain

import (
	"github.com/google/uuid"
	"time"
)

type Login struct {
	Uuid      uuid.UUID
	Login     string
	Banned    bool
	CreatedAt *time.Time
	UpdateAt  *time.Time
}
