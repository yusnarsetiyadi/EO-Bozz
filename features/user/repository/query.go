package repository

import (
	"capstone-alta1/features/user"
	"errors"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) user.RepositoryInterface {
	return &userRepository{
		db: db,
	}
}

// Create implements user.Repository
func (repo *userRepository) Create(input user.Core) error {
	userGorm := fromCore(input)
	tx := repo.db.Create(&userGorm) // proses insert data
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("insert failed")
	}
	return nil
}

// GetAll implements user.Repository
func (repo *userRepository) GetAll() (data []user.Core, err error) {
	var users []User

	tx := repo.db.Find(&users)
	if tx.Error != nil {
		return nil, tx.Error
	}
	var dataCore = toCoreList(users)
	return dataCore, nil
}

// GetById implements user.RepositoryInterface
func (repo *userRepository) GetById(id uint) (data user.Core, err error) {
	var user User

	tx := repo.db.First(&user, id)

	if tx.Error != nil {
		return data, tx.Error
	}

	if tx.RowsAffected == 0 {
		return data, tx.Error
	}

	var dataCore = user.toCore()
	return dataCore, nil
}

// Update implements user.Repository
func (repo *userRepository) Update(input user.Core, id uint) error {
	userGorm := fromCore(input)
	var user User
	tx := repo.db.Model(&user).Where("ID = ?", id).Updates(&userGorm) // proses update
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("update failed")
	}
	return nil
}

// Delete implements user.Repository
func (repo *userRepository) Delete(id uint) error {
	var user User
	tx := repo.db.Delete(&user, id) // proses delete
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("delete failed")
	}
	return nil
}

func (repo *userRepository) FindUser(email string) (result user.Core, err error) {
	var userData User
	tx := repo.db.Where("email = ?", email).First(&userData)
	if tx.Error != nil {
		return user.Core{}, tx.Error
	}

	result = userData.toCore()

	return result, nil
}

// Update Password implements user.Repository
func (repo *userRepository) UpdatePassword(input user.Core, id uint) error {
	userGorm := fromCore(input)
	var user User
	tx := repo.db.Model(&user).Where("ID = ?", id).Updates(&userGorm.Password) // proses update
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("update failed")
	}
	return nil
}
