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

func (ur *UserRepository) Insert(user *models.User) {
	ur.open()
	defer ur.close()
	ur.db.Create(user)
}

func (ur *UserRepository) CountByEmail(email string) int64 {
	ur.open()
	defer ur.close()
	var count int64
	ur.db.Table(UsersTable).Where("email = ?", email).Count(&count)
	return count
}

func (ur *UserRepository) SelectByEmail(email string) *models.User {
	ur.open()
	defer ur.close()
	var user models.User
	ur.db.Table(UsersTable).Where("email = ? ", email).Find(&user)
	return &user
}
