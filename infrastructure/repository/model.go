package repository

import (
	"github.com/google/uuid"
	"time"
)

type Login struct {
	Id        int64      `db:"-"`
	Uuid      uuid.UUID  `db:"uuid"`
	Login     string     `db:"login"`
	Banned    bool       `db:"banned"`
	CreatedAt *time.Time `db:"created_at"`
	UpdateAt  *time.Time `db:"update_at"`
}
