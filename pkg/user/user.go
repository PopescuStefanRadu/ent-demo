package user

import (
	"context"
	"golang.org/x/sync/errgroup"
	"time"
)

//go:generate mockgen -source user.go -destination mock/user.go

type User struct {
	Id          int
	Username    string
	Email       string
	DogPhotoUrl string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateUserParams struct {
	Username string
	Email    string
}

type UpdateUserParams struct {
	Id       int
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
	GetById(context.Context, int) (*User, error)
	FindAllByFilter(context.Context, *FindAllFilter) ([]User, error)
	Create(context.Context, *CreateUserParams) (*User, error)
	Update(context.Context, *UpdateUserParams) (*User, error)
	DeleteById(context.Context, int) error
	DeleteAll(ctx context.Context) (int, error)
}

type Dog interface {
	GetRandomDogUrl(ctx context.Context) (string, error)
}

func (s *Service) GetUserById(ctx context.Context, id int) (*User, error) {
	user, err := s.UserRepository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.enrichWithDogUrl(ctx, user)
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
	return s.enrichWithDogUrl(ctx, created)
}

func (s *Service) UpdateUser(ctx context.Context, u *UpdateUserParams) (*User, error) {
	updated, err := s.UserRepository.Update(ctx, u)
	if err != nil {
		return nil, err
	}
	return s.enrichWithDogUrl(ctx, updated)
}

func (s *Service) DeleteUserById(ctx context.Context, id int) error {
	return s.UserRepository.DeleteById(ctx, id)
}

func (s *Service) enrichWithDogUrl(ctx context.Context, user *User) (*User, error) {
	url, err := s.DogClient.GetRandomDogUrl(ctx)
	if err != nil {
		return nil, err
	}
	user.DogPhotoUrl = url

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
			users[v.pos].DogPhotoUrl = v.url
		}
		resultCh <- users
	}()

	for i := range users {
		iCpy := i
		g.Go(func() error {
			url, err := s.DogClient.GetRandomDogUrl(groupCtx)
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
	close(urlsCh)
	if err != nil {
		return nil, err
	}

	return <-resultCh, nil
}
