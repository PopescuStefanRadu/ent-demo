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
	"time"
)

var (
	SqlDB     *sql.DB
	EntClient *ent.Client
)

func TestMain(m *testing.M) {
	SqlDB, err := sql.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		panic(err.Error())
	}

	drv := entsql.OpenDB("sqlite3", SqlDB)
	EntClient = ent.NewClient(ent.Driver(drv))

	if code := m.Run(); code != 0 {
		panic(fmt.Errorf("test exited with non-zero code: %d", code))
	}
}

func TestCreateUser(t *testing.T) {
	r := require.New(t)

	l := zerolog.New(zerolog.NewTestWriter(t))
	ctx := l.WithContext(context.Background())

	userRepository := &entwrap.UserRepository{Client: EntClient.User}
	userService := &user.Service{UserRepository: userRepository}

	app := &App{
		Logger:   l,
		Migrator: entwrap.Migrator{Ent: EntClient},
		Service:  userService,
	}
	r.NoError(app.Init(ctx))
	t.Cleanup(func() {
		if err := app.Cleanup(ctx); err != nil {
			l.Err(err).Msg("Could not clean up application state")
		}
	})
	createdUser, err := app.CreateUser(ctx, &user.User{
		Username: "testUser",
		Email:    "testUser@mail.example",
	})

	r.NoError(err)
	r.NotEqual(0, createdUser.Id)
	r.NotEqual(time.Time{}, createdUser.UpdatedAt)
	r.NotEqual(time.Time{}, createdUser.CreatedAt)
	r.Equal("testUser", createdUser.Username)
	r.Equal("testUser@mail.example", createdUser.Email)
}
