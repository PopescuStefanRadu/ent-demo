package user

import (
	"context"
	"time"
)

type User struct {
	Id        int
	Username  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Service struct {
	UserRepository Repository
}

type FindAllFilter struct {
	IdsIn []int
}

type Repository interface {
	GetById(context.Context, int) (*User, error)
	FindAllByFilter(context.Context, *FindAllFilter) ([]User, error)
	Create(context.Context, *User) (*User, error)
	Update(context.Context, *User) (*User, error)
	DeleteById(context.Context, int) error
}

func (s *Service) GetUserById(ctx context.Context, id int) (*User, error) {
	return s.UserRepository.GetById(ctx, id)
}

func (s *Service) FindAllUsersByFilter(ctx context.Context, filter *FindAllFilter) ([]User, error) {
	return s.UserRepository.FindAllByFilter(ctx, filter)
}

func (s *Service) CreateUser(ctx context.Context, u *User) (*User, error) {
	return s.UserRepository.Create(ctx, u)
}

func (s *Service) UpdateUser(ctx context.Context, u *User) (*User, error) {
	return s.UserRepository.Update(ctx, u)
}

func (s *Service) DeleteUserById(ctx context.Context, id int) error {
	return s.UserRepository.DeleteById(ctx, id)
}