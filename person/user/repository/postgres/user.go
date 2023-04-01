package postgres

import (
	"context"
	"gorm.io/gorm"
	"person/models"
	userPkg "person/user"
)

type User struct {
	ID      uint
	Name    string
	Address string
	Work    string
	Age     int
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	rep := &UserRepository{db}
	err := rep.db.AutoMigrate(&User{})
	if err != nil {
		panic(err)
	}
	return rep
}

func (r UserRepository) CreateUser(ctx context.Context, user *models.User) (int, error) {
	model := toPostgresUser(user)

	result := r.db.Create(model)
	if result.Error != nil {
		return -1, result.Error
	}

	return int(model.ID), nil
}

func (r UserRepository) GetUser(ctx context.Context, id int) (*models.User, error) {
	var user User
	res := r.db.First(&user, id)
	if res.Error != nil {
		return nil, res.Error //userPkg.ErrUserNotFound
	}

	return ToModel(&user), nil
}

func (r UserRepository) ChangeUser(ctx context.Context, user *models.User, id int) (*models.User, error) {
	model := toPostgresUser(user)
	model.ID = uint(id)

	res := r.db.Save(model)
	if res.Error != nil {
		return nil, res.Error //userPkg.ErrUserAlreadyExists
	}

	number := res.RowsAffected
	if number == 0 {
		return nil, userPkg.ErrUserNotFound
	}

	return ToModel(model), nil
}

func (r UserRepository) DeleteUser(ctx context.Context, id int) error {
	var user User
	r.db.First(&user, id)

	res := r.db.Delete(&user)

	return res.Error
}

func (r UserRepository) GetAll(ctx context.Context) ([]*models.User, error) {
	var persons []User
	res := r.db.Find(&persons)
	if res.Error != nil {
		return nil, res.Error
	}

	return ToModelSlice(persons), nil
}

func toPostgresUser(u *models.User) *User {
	return &User{
		Name:    u.Name,
		Work:    u.Work,
		Address: u.Address,
		Age:     u.Age,
	}
}

func ToModel(u *User) *models.User {
	return &models.User{
		Id:      int(u.ID),
		Name:    u.Name,
		Work:    u.Work,
		Address: u.Address,
		Age:     u.Age,
	}
}

func ToModelSlice(pgPersons []User) []*models.User {
	var persons []*models.User

	for i := 0; i < len(pgPersons); i++ {
		persons = append(persons, ToModel(&pgPersons[i]))
	}

	return persons
}