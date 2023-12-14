package pkg

import (
	"context"
	"github.com/PopescuStefanRadu/ent-demo/pkg/user"
	"github.com/rs/zerolog"
)

type App struct {
	Logger   zerolog.Logger
	Migrator Migrator
	*user.Service
}

type Migrator interface {
	Migrate(ctx context.Context) error
}

func (a App) Init(ctx context.Context) error {
	a.Logger.Info().Msg("Migrating")
	err := a.Migrator.Migrate(ctx)
	if err == nil {
		a.Logger.Info().Msg("Migration complete")
	}
	return err
}

func (a App) Cleanup(ctx context.Context) error {
	a.Logger.Info().Msg("Cleaning up application state")
	_, err := a.UserRepository.DeleteAll(ctx)
	if err != nil {
		return err
	}
	a.Logger.Info().Msg("Finished cleaning up application state")
	return nil
}
