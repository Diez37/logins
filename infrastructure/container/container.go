package container

import (
	"github.com/Diez37/logins/infrastructure/repository"
	"github.com/diez37/go-packages/container"
)

func AddProvide(container container.Container) error {
	return container.Provides(
		repository.NewSql,
	)
}
