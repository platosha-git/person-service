package usecase

import (
	"context"
	"person/models"
	userPkg "person/user"
)

type UserUseCase struct {
	userRepo userPkg.Repository
}

func NewUserUseCase(userRepo userPkg.Repository) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

func (u *UserUseCase) Create(ctx context.Context, user *models.User) (int, error) {
	userId, err := u.userRepo.CreateUser(ctx, user)
	if err != nil {
		return -1, err
	}

	return userId, nil
}

func (u *UserUseCase) GetProfile(ctx context.Context, id int) (*models.User, error) {
	user, err := u.userRepo.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserUseCase) ChangeProfile(ctx context.Context, user *models.User, id int) (*models.User, error) {
	oldUser, err := u.userRepo.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	// do nothing if we dont need to change anything
	if user.Address == "" && user.Work == "" && user.Name == "" && user.Age == 0 {
		return oldUser, nil
	}
	// check empty fields
	if user.Work == "" {
		user.Work = oldUser.Work
	}
	if user.Name == "" {
		user.Name = oldUser.Name
	}
	if user.Address == "" {
		user.Address = oldUser.Address
	}
	if user.Age == 0 {
		user.Age = oldUser.Age
	}

	newUser, err := u.userRepo.ChangeUser(ctx, user, id)
	return newUser, err
}

func (u *UserUseCase) DeletePerson(ctx context.Context, id int) error {
	err := u.userRepo.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserUseCase) GetAllPersons(ctx context.Context) ([]*models.User, error) {
	persons, err := u.userRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return persons, nil
}
