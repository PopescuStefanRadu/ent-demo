package pkg_test

import (
	"database/sql"
	"fmt"
	"github.com/PopescuStefanRadu/ent-demo/pkg"
	"github.com/PopescuStefanRadu/ent-demo/pkg/user"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

var (
	SqlDB *sql.DB
)

func TestMain(m *testing.M) {
	var err error
	SqlDB, err = sql.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		panic(err.Error())
	}

	if code := m.Run(); code != 0 {
		panic(fmt.Errorf("test exited with non-zero code: %d", code))
	}
}

func TestCreateUser(t *testing.T) {
	r, _, ctx, app, mocks := pkg.InitTest(t, SqlDB)

	mocks.DogClient.EXPECT().GetRandomDogUrl(gomock.Any()).Return("https://example.org", nil).Times(1)

	createdUser, err := app.CreateUser(ctx, &user.CreateUserParams{
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

func TestGetUserById(t *testing.T) {
	r, _, ctx, app, mocks := pkg.InitTest(t, SqlDB)

	mocks.DogClient.EXPECT().GetRandomDogUrl(gomock.Any()).Return("https://example.org", nil).Times(2)

	createdUser, err := app.CreateUser(ctx, &user.CreateUserParams{
		Username: "testUser",
		Email:    "testUser@mail.example",
	})
	r.NoError(err)

	userById, err := app.GetUserById(ctx, createdUser.Id)
	r.NoError(err)

	r.Equal(createdUser, userById)
	r.Equal(&user.User{
		Id:          createdUser.Id,
		Username:    "testUser",
		Email:       "testUser@mail.example",
		DogPhotoUrl: "https://example.org",
		CreatedAt:   createdUser.CreatedAt,
		UpdatedAt:   createdUser.UpdatedAt,
	}, createdUser)
}

func TestUpdateUser(t *testing.T) {
	r, _, ctx, app, mocks := pkg.InitTest(t, SqlDB)

	mocks.DogClient.EXPECT().GetRandomDogUrl(gomock.Any()).Return("https://example.org", nil).Times(2)

	createdUser, err := app.CreateUser(ctx, &user.CreateUserParams{
		Username: "testUser",
		Email:    "testUser@mail.example",
	})
	r.NoError(err)
	r.NotNil(createdUser)

	updatedUser, err := app.UpdateUser(ctx, &user.UpdateUserParams{
		Id:       createdUser.Id,
		Username: "testUser2",
		Email:    "testUser2@mail.example",
	})
	r.NoError(err)

	r.Equal(createdUser.Id, updatedUser.Id)
	r.Equal("testUser2", updatedUser.Username)
	r.Equal("testUser2@mail.example", updatedUser.Email)
	r.Equal(createdUser.CreatedAt, updatedUser.CreatedAt)
	r.Equal(createdUser.DogPhotoUrl, updatedUser.DogPhotoUrl)
	r.LessOrEqual(createdUser.UpdatedAt, updatedUser.UpdatedAt)
}

func TestDeleteUser(t *testing.T) {
	r, _, ctx, app, mocks := pkg.InitTest(t, SqlDB)

	mocks.DogClient.EXPECT().GetRandomDogUrl(gomock.Any()).Return("https://example.org", nil).Times(3)

	createdUser, err := app.CreateUser(ctx, &user.CreateUserParams{
		Username: "testUser",
		Email:    "testUser@mail.example",
	})
	r.NoError(err)
	r.NotNil(createdUser)

	createdUser2, err := app.CreateUser(ctx, &user.CreateUserParams{
		Username: "testUser2",
		Email:    "testUser2@mail.example",
	})
	r.NoError(err)
	r.NotNil(createdUser2)

	r.NoError(app.DeleteUserById(ctx, createdUser.Id))

	allUsers, err := app.FindAllUsersByFilter(ctx, nil)
	r.NoError(err)

	require.Equal(t, []user.User{{
		Id:          createdUser2.Id,
		Username:    "testUser2",
		Email:       "testUser2@mail.example",
		DogPhotoUrl: "https://example.org",
		CreatedAt:   createdUser2.CreatedAt,
		UpdatedAt:   createdUser2.UpdatedAt,
	}}, allUsers)
}

func TestGetUsersByIds(t *testing.T) {
	r, _, ctx, app, mocks := pkg.InitTest(t, SqlDB)

	mocks.DogClient.EXPECT().GetRandomDogUrl(gomock.Any()).Return("https://example.org", nil).Times(5)

	usersToCreate := []*user.CreateUserParams{
		ToPtr(user.CreateUserParams{
			Username: "testUser",
			Email:    "testUser@mail.example",
		}),
		ToPtr(user.CreateUserParams{
			Username: "testUser2",
			Email:    "testUser2@mail.example",
		}),
		ToPtr(user.CreateUserParams{
			Username: "testUser3",
			Email:    "testUser3@mail.example",
		}),
	}

	expectedUsers := make([]user.User, len(usersToCreate))
	for i, u := range usersToCreate {
		createdUser, err := app.CreateUser(ctx, u)
		r.NoError(err)
		r.NotNil(createdUser)
		expectedUsers[i] = *createdUser
	}

	actualUsers, err := app.FindAllUsersByFilter(ctx, &user.FindAllFilter{IdsIn: []int{expectedUsers[0].Id, expectedUsers[1].Id}})

	r.NoError(err)
	r.Equal(expectedUsers[0:2], actualUsers)
}

func ToPtr[T any](t T) *T {
	p := &t
	return p
}
