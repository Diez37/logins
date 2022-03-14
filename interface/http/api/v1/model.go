package v1

import (
	"github.com/google/uuid"
	"time"
)

type Page struct {
	Meta    *Meta    `json:"meta"`
	Records []*Login `json:"records"`
}

type Meta struct {
	Count int64 `json:"count"`
	Page  uint  `json:"page"`
	Limit uint  `json:"limit"`
}

type Login struct {
	Uuid      uuid.UUID  `json:"uuid"`
	Login     string     `json:"login"`
	Ban       bool       `json:"ban"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdateAt  *time.Time `json:"updateAt"`
}
