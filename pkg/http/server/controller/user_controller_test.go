package controller_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	application "github.com/PopescuStefanRadu/ent-demo/pkg/app"
	"github.com/PopescuStefanRadu/ent-demo/pkg/http/server"
	"github.com/PopescuStefanRadu/ent-demo/pkg/http/server/request"
	"github.com/PopescuStefanRadu/ent-demo/pkg/http/server/response"
	"github.com/PopescuStefanRadu/ent-demo/pkg/user"
	"go.uber.org/mock/gomock"
)

var SqlDB *sql.DB //nolint:golint,revive,stylecheck

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

func TestCreate(t *testing.T) {
	r, _, ctx, app, mocks := application.InitTest(t, SqlDB)

	mocks.DogClient.EXPECT().GetRandomDogURL(gomock.Any()).Return("https://example.org", nil).Times(1)

	gin := server.NewRouter(app)

	body, err := json.Marshal(request.CreateUser{
		Username: "testUser",
		Email:    "testUser@example.com",
	})
	r.NoError(err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/user", bytes.NewReader(body))
	r.NoError(err)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	gin.ServeHTTP(w, req)

	r.Equal(http.StatusOK, w.Code)

	var actualResp response.Response[response.User]
	r.NoError(json.Unmarshal(w.Body.Bytes(), &actualResp))

	r.NotEqual(0, actualResp.Result.ID)
	r.NotEqual(time.Time{}, actualResp.Result.CreatedAt)
	r.NotEqual(time.Time{}, actualResp.Result.UpdatedAt)

	r.Equal(response.Response[response.User]{
		Result: response.User{
			ID:          actualResp.Result.ID,
			Username:    "testUser",
			Email:       "testUser@example.com",
			DogPhotoURL: "https://example.org",
			CreatedAt:   actualResp.Result.CreatedAt,
			UpdatedAt:   actualResp.Result.UpdatedAt,
		},
		Errors: nil,
	}, actualResp)
}

func TestGet(t *testing.T) {
	r, _, ctx, app, mocks := application.InitTest(t, SqlDB)

	mocks.DogClient.EXPECT().GetRandomDogURL(gomock.Any()).Return("https://example.org", nil).Times(2)

	gin := server.NewRouter(app)

	usr, err := app.CreateUser(ctx, &user.CreateUserParams{
		Username: "testUser",
		Email:    "testUser@example.com",
	})
	r.NoError(err)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("/user/%d", usr.ID), nil)
	r.NoError(err)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	gin.ServeHTTP(w, req)

	var actualResp response.Response[response.User]
	r.NoError(json.Unmarshal(w.Body.Bytes(), &actualResp))

	r.Equal(http.StatusOK, w.Code)
	r.Equal(response.Response[response.User]{
		Result: response.User{
			ID:          actualResp.Result.ID,
			Username:    "testUser",
			Email:       "testUser@example.com",
			DogPhotoURL: "https://example.org",
			CreatedAt:   actualResp.Result.CreatedAt,
			UpdatedAt:   actualResp.Result.UpdatedAt,
		},
	}, actualResp)
}

func TestUpdate(t *testing.T) {
	r, _, ctx, app, mocks := application.InitTest(t, SqlDB)

	mocks.DogClient.EXPECT().GetRandomDogURL(gomock.Any()).Return("https://example.org", nil).Times(2)

	gin := server.NewRouter(app)

	usr, err := app.CreateUser(ctx, &user.CreateUserParams{
		Username: "testUser",
		Email:    "testUser@example.com",
	})
	r.NoError(err)

	body, err := json.Marshal(request.UpdateUserBody{
		Username: "updatedTestUser",
		Email:    "updatedTestUser@example.com",
	})
	r.NoError(err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, fmt.Sprintf("/user/%d", usr.ID), bytes.NewReader(body))
	r.NoError(err)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	gin.ServeHTTP(w, req)

	var actualResp response.Response[response.User]
	r.NoError(json.Unmarshal(w.Body.Bytes(), &actualResp))

	r.Equal(http.StatusOK, w.Code)
	r.Equal(response.Response[response.User]{
		Result: response.User{
			ID:          actualResp.Result.ID,
			Username:    "updatedTestUser",
			Email:       "updatedTestUser@example.com",
			DogPhotoURL: "https://example.org",
			CreatedAt:   actualResp.Result.CreatedAt,
			UpdatedAt:   actualResp.Result.UpdatedAt,
		},
	}, actualResp)
}

func TestDelete(t *testing.T) {
	r, _, ctx, app, mocks := application.InitTest(t, SqlDB)

	mocks.DogClient.EXPECT().GetRandomDogURL(gomock.Any()).Return("https://example.org", nil).Times(1)

	gin := server.NewRouter(app)

	usr, err := app.CreateUser(ctx, &user.CreateUserParams{
		Username: "testUser",
		Email:    "testUser@example.com",
	})
	r.NoError(err)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("/user/%d", usr.ID), nil)
	r.NoError(err)
	w := httptest.NewRecorder()

	gin.ServeHTTP(w, req)

	r.Equal("{}", w.Body.String())

	remainingUsers, err := app.FindAllUsersByFilter(ctx, nil)
	r.NoError(err)
	r.Equal(0, len(remainingUsers))
}

func TestGetFiltered(t *testing.T) {
	r, _, ctx, app, mocks := application.InitTest(t, SqlDB)

	mocks.DogClient.EXPECT().GetRandomDogURL(gomock.Any()).Return("https://example.org", nil).Times(5)

	gin := server.NewRouter(app)

	usr1, err := app.CreateUser(ctx, &user.CreateUserParams{
		Username: "testUser1",
		Email:    "testUser1@example.com",
	})
	r.NoError(err)

	usr2, err := app.CreateUser(ctx, &user.CreateUserParams{
		Username: "testUser2",
		Email:    "testUser2@example.com",
	})
	r.NoError(err)

	_, err = app.CreateUser(ctx, &user.CreateUserParams{
		Username: "testUser3",
		Email:    "testUser3@example.com",
	})
	r.NoError(err)

	body, err := json.Marshal(request.GetFilteredUsers{
		IdsIn: []int{usr1.ID, usr2.ID},
	})
	r.NoError(err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/search-users", bytes.NewReader(body))
	r.NoError(err)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	gin.ServeHTTP(w, req)

	var actualResp response.Response[[]response.User]
	r.NoError(json.Unmarshal(w.Body.Bytes(), &actualResp), w.Body.String())

	r.Equal(http.StatusOK, w.Code)
	r.Equal(response.Response[[]response.User]{
		Result: []response.User{response.User(*usr1), response.User(*usr2)},
	}, actualResp)
}
