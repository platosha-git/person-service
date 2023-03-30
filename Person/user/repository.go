package user

import (
	"context"
	"person/models"
)

type Repository interface {
	CreateUser(ctx context.Context, user *models.User) (int, error)
	GetUser(ctx context.Context, id int) (*models.User, error)
	ChangeUser(ctx context.Context, user *models.User, id int) (*models.User, error)
	DeleteUser(ctx context.Context, id int) error
	GetAll(ctx context.Context) ([]*models.User, error)
}
