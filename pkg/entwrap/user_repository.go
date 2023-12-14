package entwrap

import (
	"context"
	"github.com/PopescuStefanRadu/ent-demo/pkg/ent"
	"github.com/PopescuStefanRadu/ent-demo/pkg/ent/user"
	businessUser "github.com/PopescuStefanRadu/ent-demo/pkg/user"
)

type UserRepository struct {
	Client    *ent.UserClient
	EntClient *ent.Client
}

func (ur *UserRepository) GetById(ctx context.Context, id int) (*businessUser.User, error) {
	u, err := ur.Client.Get(ctx, id)
	return toBusinessModel(u), err
}

func (ur *UserRepository) Create(ctx context.Context, u *businessUser.User) (*businessUser.User, error) {
	err := ur.Client.Create().
		SetUsername(u.Username).
		SetEmail(u.Email).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	createdUser, err := ur.Client.Query().Where(user.Email(u.Email)).First(ctx)

	return toBusinessModel(createdUser), err
}

func (ur *UserRepository) FindAllByFilter(context.Context, *businessUser.FindAllFilter) ([]businessUser.User, error) {
	return nil, nil
}

func (ur *UserRepository) Update(context.Context, *businessUser.User) (*businessUser.User, error) {
	return nil, nil

}
func (ur *UserRepository) DeleteById(context.Context, int) error {
	return nil
}

func (ur *UserRepository) DeleteAll(ctx context.Context) (int, error) {
	return ur.Client.Delete().Exec(ctx)
}

func toBusinessModel(u *ent.User) *businessUser.User {
	if u == nil {
		return nil
	}

	return &businessUser.User{
		Id:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
