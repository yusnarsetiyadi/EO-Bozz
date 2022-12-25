package repository

import (
	partner "capstone-alta1/features/partner"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type partnerRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) partner.RepositoryInterface {
	return &partnerRepository{
		db: db,
	}
}

// Create implements user.Repository
func (repo *partnerRepository) Create(input partner.Core) error {
	partnerGorm := fromCore(input)
	tx := repo.db.Create(&partnerGorm) // proses insert data
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("insert failed")
	}
	return nil
}

// GetAll implements user.Repository
func (repo *partnerRepository) GetAll() (data []partner.Core, err error) {
	var partner []Partner

	tx := repo.db.Preload("User").Find(&partner)
	if tx.Error != nil {
		return nil, tx.Error
	}
	var dataCore = toCoreList(partner)
	return dataCore, nil
}

func (repo *partnerRepository) GetAllWithSearch(query string) (data []partner.Core, err error) {
	var partner []Partner
	tx := repo.db.Preload("User").Where("name LIKE ?", "%"+query+"%").Find(&partner)
	if tx.Error != nil {
		return nil, tx.Error
	}

	if tx.RowsAffected == 0 {
		return data, tx.Error
	}

	fmt.Println("\n\n 2 getall partner = ", partner)

	var dataCore = toCoreList(partner)
	return dataCore, nil
}

// GetById implements user.RepositoryInterface
func (repo *partnerRepository) GetById(id int) (data partner.Core, err error) {
	var partner Partner

	tx := repo.db.Preload("User").First(&partner, id)

	if tx.Error != nil {
		return data, tx.Error
	}

	if tx.RowsAffected == 0 {
		return data, tx.Error
	}

	var dataCore = partner.toCore()
	return dataCore, nil
}

// Update implements user.Repository
func (repo *partnerRepository) Update(input partner.Core, id int) error {
	partnerGorm := fromCore(input)
	var partner Partner
	var user User
	repo.db.Transaction(func(tx *gorm.DB) error {
		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		if err := tx.Model(&partner).Where("ID = ?", id).Updates(&partnerGorm).Error; err != nil {
			// return any error will rollback
			return err
		}

		if err := tx.Model(&user).Where("ID = ?", id).Updates(&partnerGorm).Error; err != nil {
			return err
		}

		if tx.RowsAffected == 0 {
			return errors.New("update failed")
		}

		// return nil will commit the whole transaction
		return nil
	})

	return nil
}

// Delete implements user.Repository
func (repo *partnerRepository) Delete(id int) error {
	var partner Partner
	var user User
	repo.db.Transaction(func(tx *gorm.DB) error {
		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		if err := tx.Delete(&partner, id).Error; err != nil {
			// return any error will rollback
			return err
		}

		if err := tx.Delete(&user, id).Error; err != nil {
			return err
		}

		if tx.RowsAffected == 0 {
			return errors.New("update failed")
		}

		// return nil will commit the whole transaction
		return nil
	})

	return nil
}

func (repo *partnerRepository) FindUser(email string) (result partner.Core, err error) {
	var partnerData Partner
	tx := repo.db.Where("email", email).First(&partnerData.User)
	if tx.Error != nil {
		return partner.Core{}, tx.Error
	}

	result = partnerData.toCore()

	return result, nil
}