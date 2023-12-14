package user

import "time"

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
	GetById(int) (User, error)
	FindAllByFilter(FindAllFilter) ([]User, error)
	Create(User) (User, error)
	Update(User) (User, error)
	DeleteById(int) error
}

func (s *Service) GetUserById(id int) (User, error) {
	return s.UserRepository.GetById(id)
}

func (s *Service) FindAllUsersByFilter(filter FindAllFilter) ([]User, error) {
	return s.UserRepository.FindAllByFilter(filter)
}

func (s *Service) CreateUser(u User) (User, error) {
	return s.UserRepository.Create(u)
}

func (s *Service) UpdateUser(u User) (User, error) {
	return s.UserRepository.Update(u)
}

func (s *Service) DeleteUserById(id int) error {
	return s.UserRepository.DeleteById(id)
}
