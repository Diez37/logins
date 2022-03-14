package repository

import (
	"context"
	"github.com/google/uuid"
)

type Finder interface {
	FindByUuid(ctx context.Context, uuid uuid.UUID) (*Login, error)
	FindByLogin(ctx context.Context, login string) (*Login, error)
}

type Saver interface {
	Insert(ctx context.Context, login *Login) (*Login, error)
	Update(ctx context.Context, login *Login) (*Login, error)
}

type Blocker interface {
	BanByUuid(ctx context.Context, uuid uuid.UUID) (bool, error)
}

type Getter interface {
	Count(ctx context.Context) (int64, error)
	Page(ctx context.Context, page uint, limit uint) ([]*Login, error)
}

type Repository interface {
	Finder
	Saver
	Blocker
	Getter
}
