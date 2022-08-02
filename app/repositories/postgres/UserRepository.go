package postgres

import (
	"go-membership/app/models"
)

type UserRepository struct {
	connection
}

const UsersTable = "users"

func NewUserRepository() *UserRepository {
	return &UserRepository{
		NewConnection("postgres"),
	}
}

func (ur *UserRepository) Insert(user *models.User) error {
	err := ur.open()
	if err != nil {
		return err
	}

	defer ur.close()
	ur.db.Create(user)
	return nil
}

func (ur *UserRepository) CountByEmail(email string) (int64, error) {
	err := ur.open()
	if err != nil {
		return 0, err
	}
	defer ur.close()
	var count int64
	ur.db.Table(UsersTable).Where("email = ?", email).Count(&count)
	return count, nil
}

func (ur *UserRepository) SelectByEmail(email string) (*models.User, error) {
	err := ur.open()
	if err != nil {
		return nil, err
	}
	defer ur.close()
	var user models.User
	ur.db.Table(UsersTable).Where("email = ? ", email).Find(&user)
	return &user, err
}
