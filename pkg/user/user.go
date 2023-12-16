package user

import (
	"context"
	"time"

	"golang.org/x/sync/errgroup"
)

//go:generate mockgen -source user.go -destination mock/user.go

type User struct {
	ID          int
	Username    string
	Email       string
	DogPhotoURL string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateUserParams struct {
	Username string
	Email    string
}

type UpdateUserParams struct {
	ID       int
	Username string
	Email    string
}

type Service struct {
	UserRepository Repository
	DogClient      Dog
}

type FindAllFilter struct {
	IdsIn []int
}

type Repository interface {
	GetByID(ctx context.Context, id int) (*User, error)
	FindAllByFilter(ctx context.Context, findParams *FindAllFilter) ([]User, error)
	Create(ctx context.Context, createParams *CreateUserParams) (*User, error)
	Update(ctx context.Context, updateParams *UpdateUserParams) (*User, error)
	DeleteByID(ctx context.Context, id int) error
	DeleteAll(ctx context.Context) (int, error)
}

type Dog interface {
	GetRandomDogURL(ctx context.Context) (string, error)
}

func (s *Service) GetUserByID(ctx context.Context, id int) (*User, error) {
	user, err := s.UserRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.enrichWithDogURL(ctx, user)
}

func (s *Service) FindAllUsersByFilter(ctx context.Context, filter *FindAllFilter) ([]User, error) {
	users, err := s.UserRepository.FindAllByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}

	return s.parallelEnrichWithDogUrls(ctx, users)
}

func (s *Service) CreateUser(ctx context.Context, u *CreateUserParams) (*User, error) {
	created, err := s.UserRepository.Create(ctx, u)
	if err != nil {
		return nil, err
	}

	return s.enrichWithDogURL(ctx, created)
}

func (s *Service) UpdateUser(ctx context.Context, u *UpdateUserParams) (*User, error) {
	updated, err := s.UserRepository.Update(ctx, u)
	if err != nil {
		return nil, err
	}

	return s.enrichWithDogURL(ctx, updated)
}

func (s *Service) DeleteUserByID(ctx context.Context, id int) error {
	return s.UserRepository.DeleteByID(ctx, id)
}

func (s *Service) enrichWithDogURL(ctx context.Context, user *User) (*User, error) {
	url, err := s.DogClient.GetRandomDogURL(ctx)
	if err != nil {
		return nil, err
	}

	user.DogPhotoURL = url

	return user, err
}

func (s *Service) parallelEnrichWithDogUrls(ctx context.Context, users []User) ([]User, error) {
	g, groupCtx := errgroup.WithContext(ctx)
	urlsCh := make(chan struct {
		pos int
		url string
	})

	resultCh := make(chan []User)
	go func() {
		defer close(resultCh)

		for v := range urlsCh {
			users[v.pos].DogPhotoURL = v.url
		}

		resultCh <- users
	}()

	for i := range users {
		iCpy := i

		g.Go(func() error {
			url, err := s.DogClient.GetRandomDogURL(groupCtx)
			if err != nil {
				return err
			}

			urlsCh <- struct {
				pos int
				url string
			}{pos: iCpy, url: url}

			return nil
		})
	}

	err := g.Wait()
	close(urlsCh) //nolint:wsl // special case

	if err != nil { //nolint:wsl,nolintlint
		return nil, err
	}

	return <-resultCh, nil
}
