package pkg

import (
	"context"
	"database/sql"
	entsql "entgo.io/ent/dialect/sql"
	"fmt"
	"github.com/PopescuStefanRadu/ent-demo/pkg/ent"
	"github.com/PopescuStefanRadu/ent-demo/pkg/entwrap"
	"github.com/PopescuStefanRadu/ent-demo/pkg/user"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"testing"
)

var DBDriverName = "sqlite3"

type Config struct {
	DBUrl            string
	DebugPersistence bool
}

type App struct {
	Logger   zerolog.Logger
	Migrator Migrator
	*user.Service
}

type Migrator interface {
	Migrate(ctx context.Context) error
}

func NewAppFromConfig(l zerolog.Logger, cfg *Config) (*App, error) {
	sqlDB, err := sql.Open(DBDriverName, cfg.DBUrl)
	if err != nil {
		return nil, fmt.Errorf("could not open db: %w", err)
	}

	drv := entsql.OpenDB(DBDriverName, sqlDB)

	opts := []ent.Option{ent.Driver(drv)}
	if cfg.DebugPersistence {
		opts = append(opts, ent.Log(func(a ...any) {
			l.Info().Msgf("ent: %s", fmt.Sprint(a))
		}))
	}

	EntClient := ent.NewClient(opts...)

	userRepository := &entwrap.UserRepository{Client: EntClient.User}
	userService := &user.Service{UserRepository: userRepository}

	return &App{
		Logger:   l,
		Migrator: entwrap.Migrator{Ent: EntClient, Logger: l},
		Service:  userService,
	}, nil
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

func InitTest(t *testing.T, db *sql.DB) (*require.Assertions, zerolog.Logger, context.Context, *App) {
	r := require.New(t)
	l := zerolog.New(zerolog.NewTestWriter(t))
	ctx := l.WithContext(context.Background())
	app := initApp(ctx, t, l, db)

	return r, l, ctx, app
}

func initApp(ctx context.Context, t *testing.T, l zerolog.Logger, db *sql.DB) *App {
	drv := entsql.OpenDB("sqlite3", db)
	EntClient := ent.NewClient(ent.Driver(drv), ent.Log(func(a ...any) {
		l.Info().Msgf("ent: %s", fmt.Sprint(a))
	}), ent.Debug())

	userRepository := &entwrap.UserRepository{Client: EntClient.User}
	userService := &user.Service{UserRepository: userRepository}

	app := &App{
		Logger:   l,
		Migrator: entwrap.Migrator{Ent: EntClient, Logger: l},
		Service:  userService,
	}

	require.NoError(t, app.Init(ctx))

	t.Cleanup(func() {
		if err := app.Cleanup(ctx); err != nil {
			l.Err(err).Msg("Could not clean up application state")
		}
	})

	return app
}
