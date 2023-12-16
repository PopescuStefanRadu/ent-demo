package entwrap

import (
	"context"

	"github.com/PopescuStefanRadu/ent-demo/pkg/ent"
	"github.com/PopescuStefanRadu/ent-demo/pkg/ent/user"
	businessUser "github.com/PopescuStefanRadu/ent-demo/pkg/user"
)

type UserRepository struct {
	Client *ent.UserClient
}

func (ur *UserRepository) GetByID(ctx context.Context, id int) (*businessUser.User, error) {
	u, err := ur.Client.Get(ctx, id)
	return toPtrBusinessModel(u), err
}

func (ur *UserRepository) Create(ctx context.Context, u *businessUser.CreateUserParams) (*businessUser.User, error) {
	createdUser, err := ur.Client.Create().
		SetUsername(u.Username).
		SetEmail(u.Email).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return toPtrBusinessModel(createdUser), err
}

func (ur *UserRepository) FindAllByFilter(
	ctx context.Context,
	filter *businessUser.FindAllFilter,
) ([]businessUser.User, error) {
	if filter == nil || len(filter.IdsIn) == 0 {
		users, err := ur.Client.Query().Where().All(ctx)
		if err != nil {
			return nil, err
		}

		return toBusinessModelSlice(users), nil
	}

	filteredUsers, err := ur.Client.Query().Where(user.IDIn(filter.IdsIn...)).All(ctx)
	if err != nil {
		return nil, err
	}

	return toBusinessModelSlice(filteredUsers), nil
}

func (ur *UserRepository) Update(ctx context.Context, u *businessUser.UpdateUserParams) (*businessUser.User, error) {
	err := ur.Client.Update().
		Where(user.ID(u.ID)).
		SetUsername(u.Username).
		SetEmail(u.Email).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return ur.GetByID(ctx, u.ID)
}

func (ur *UserRepository) DeleteByID(ctx context.Context, id int) error {
	return ur.Client.DeleteOneID(id).Exec(ctx)
}

func (ur *UserRepository) DeleteAll(ctx context.Context) (int, error) {
	return ur.Client.Delete().Exec(ctx)
}

func toBusinessModelSlice(users []*ent.User) []businessUser.User {
	if len(users) == 0 {
		return nil
	}

	res := make([]businessUser.User, len(users))

	for i, u := range users {
		res[i] = toBusinessModel(u)
	}

	return res
}

func toPtrBusinessModel(u *ent.User) *businessUser.User {
	if u == nil {
		return nil
	}

	model := toBusinessModel(u)

	return &model
}

func toBusinessModel(u *ent.User) businessUser.User {
	return businessUser.User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
