package domain

import (
	"github.com/google/uuid"
	"time"
)

type Login struct {
	Uuid      uuid.UUID
	Login     string
	Ban       bool
	CreatedAt *time.Time
	UpdateAt  *time.Time
}
