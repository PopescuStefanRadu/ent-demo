package app

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/PopescuStefanRadu/ent-demo/pkg/ent"
	"github.com/PopescuStefanRadu/ent-demo/pkg/entwrap"
	"github.com/PopescuStefanRadu/ent-demo/pkg/external/dog"
	"github.com/PopescuStefanRadu/ent-demo/pkg/user"
	mockUser "github.com/PopescuStefanRadu/ent-demo/pkg/user/mock"
	_ "github.com/mattn/go-sqlite3" // importing drivers because it's easier to have a shared initialization for the application based on the config
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const DBDriverName = "sqlite3"

type Config struct {
	DBUrl            string
	DebugPersistence bool
	DogClientConfig  dog.ClientConfig
}

type App struct {
	Logger   zerolog.Logger
	Migrator Migrator
	*user.Service
}

type Mocks struct {
	DogClient *mockUser.MockDog
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
	dogClient := dog.NewClient(cfg.DogClientConfig)
	userService := &user.Service{UserRepository: userRepository, DogClient: dogClient}

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

func InitTest(t *testing.T, db *sql.DB) (*require.Assertions, zerolog.Logger, context.Context, *App, Mocks) {
	t.Helper()
	r := require.New(t)
	l := zerolog.New(zerolog.NewTestWriter(t))
	ctx := l.WithContext(context.Background())

	ctl := gomock.NewController(t)

	mockDog := mockUser.NewMockDog(ctl)
	mocks := Mocks{DogClient: mockDog}

	app := initApp(ctx, t, l, db, mocks)

	return r, l, ctx, app, mocks
}

func initApp(ctx context.Context, t *testing.T, l zerolog.Logger, db *sql.DB, mocks Mocks) *App {
	t.Helper()

	drv := entsql.OpenDB("sqlite3", db)
	EntClient := ent.NewClient(ent.Driver(drv), ent.Log(func(a ...any) {
		l.Info().Msgf("ent: %s", fmt.Sprint(a))
	}), ent.Debug())

	userRepository := &entwrap.UserRepository{Client: EntClient.User}
	userService := &user.Service{UserRepository: userRepository, DogClient: mocks.DogClient}

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
