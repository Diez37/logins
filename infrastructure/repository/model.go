package repository

import (
	"github.com/google/uuid"
	"time"
)

type Login struct {
	Uuid      uuid.UUID  `db:"uuid"`
	Login     string     `db:"login"`
	Ban       bool       `db:"ban"`
	CreatedAt *time.Time `db:"created_at"`
	UpdateAt  *time.Time `db:"update_at"`
}