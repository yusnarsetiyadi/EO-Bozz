package repository

import (
	_service "capstone-alta1/features/service"
	"capstone-alta1/utils/helper"
	"errors"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type serviceRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) _service.RepositoryInterface {
	return &serviceRepository{
		db: db,
	}
}

func (repo *serviceRepository) Create(input _service.Core) error {
	serviceGorm := fromCore(input)
	tx := repo.db.Create(&serviceGorm) // proses insert data
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("insert failed")
	}
	return nil
}

func (repo *serviceRepository) GetAll() (data []_service.Core, err error) {
	var results []Service

	tx := repo.db.Find(&results)
	if tx.Error != nil {
		return nil, tx.Error
	}
	var dataCore = toCoreList(results)
	return dataCore, nil
}

func (repo *serviceRepository) GetAllWithSearch(queryName, queryCategory, queryCity, queryMinPrice, queryMaxPrice string) (data []_service.Core, err error) {
	var services, services2 []Service

	helper.LogDebug("\n isi queryName = ", queryName)
	helper.LogDebug("\n isi queryCategory= ", queryCategory)
	helper.LogDebug("\n isi queryCity = ", queryCity)
	helper.LogDebug("\n isi queryMinPrice = ", queryMinPrice)
	helper.LogDebug("\n isi queryMaxPrice = ", queryMaxPrice)

	intMinPrice, errConv1 := strconv.Atoi(queryMinPrice)
	intMaxPrice, errConv2 := strconv.Atoi(queryMaxPrice)
	if errConv1 != nil || errConv2 != nil {
		return nil, errors.New("error conver service price to filter")
	}

	fmt.Println("\n\nServices 1", services)
	tx := repo.db.Where("service_name LIKE ?", "%"+queryName+"%").Where(&Service{ServiceCategory: queryCategory, City: queryCity, ServicePrice: uint(intMinPrice) + uint(intMaxPrice)}).Find(&services2)
	fmt.Println("\n\nServices 2", services2)

	if tx.Error != nil {
		return nil, tx.Error
	}
	var dataCore = toCoreList(services2)
	return dataCore, nil
}

func (repo *serviceRepository) GetById(id uint) (data _service.Core, err error) {
	var service Service

	tx := repo.db.First(&service, id)

	if tx.Error != nil {
		return data, tx.Error
	}

	if tx.RowsAffected == 0 {
		return data, tx.Error
	}

	var dataCore = service.toCore()
	return dataCore, nil
}

func (repo *serviceRepository) Update(input _service.Core, id uint) error {
	resultGorm := fromCore(input)
	var result Service
	tx := repo.db.Model(&result).Where("ID = ?", id).Updates(&resultGorm) // proses update
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("update failed")
	}
	return nil
}

func (repo *serviceRepository) Delete(id uint) error {
	var result Service
	tx := repo.db.Delete(&result, id) // proses delete
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("delete failed")
	}
	return nil
}

func (repo *serviceRepository) GetAdditionalById(id uint) (data []_service.Additional, err error) {
	var clientadditional []Additional

	tx := repo.db.Find(&clientadditional, id)

	if tx.Error != nil {
		return data, tx.Error
	}

	if tx.RowsAffected == 0 {
		return data, tx.Error
	}

	var dataCore = toCoreListAdditional(clientadditional)
	return dataCore, nil
}

func (repo *serviceRepository) AddAdditionalToService(input _service.ServiceAdditional) error {
	additionalGorm := fromCoreServiceAdditional(input)
	var service Service
	tx := repo.db.Model(&service).Where("ID = ?", input.ServiceID).Create(&additionalGorm) // proses insert data
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("insert failed")
	}
	return nil
}

func (repo *serviceRepository) GetReviewById(id uint) (data []_service.Review, err error) {
	var clientreview []Review

	tx := repo.db.Find(&clientreview, id)

	if tx.Error != nil {
		return data, tx.Error
	}

	if tx.RowsAffected == 0 {
		return data, tx.Error
	}

	var dataCore = toCoreListReview(clientreview)
	return dataCore, nil
}

func (repo *serviceRepository) GetDiscussionById(id uint) (data []_service.Discussion, err error) {
	var clientdiscussion []Discussion

	tx := repo.db.Find(&clientdiscussion, id)

	if tx.Error != nil {
		return data, tx.Error
	}

	if tx.RowsAffected == 0 {
		return data, tx.Error
	}

	var dataCore = toCoreListDiscussion(clientdiscussion)
	return dataCore, nil
}

func (repo *serviceRepository) CheckAvailability(serviceId uint, queryStart, queryEnd string) (data _service.Order, err error) {
	//convert string to format time
	layoutFormat := "2006-01-02"
	timeStart, _ := time.Parse(layoutFormat, queryStart)
	timeEnd, _ := time.Parse(layoutFormat, queryEnd)

	//check available
	var service []Service
	queryBuilder := fmt.Sprintf("SELECT * FROM orders WHERE service_id = %d AND (('%s' BETWEEN start_date AND end_date) OR ('%s' BETWEEN start_date AND end_date));", serviceId, timeStart, timeEnd)
	fmt.Println("\n\n query ", queryBuilder)
	tx := repo.db.Raw(queryBuilder).Find(&service)

	//create return
	var services Service
	var orders Order
	nameService := services.ServiceName
	statusAvailable := "Available"
	statusNotvalable := "Not Available"

	if tx.Error != nil {
		return orders.toCoreNotAvailable(nameService, timeStart, timeEnd, statusNotvalable), tx.Error
	}

	affectedRow := tx.RowsAffected
	fmt.Println("\n\nHasil check availbility, \n Checkin date = ", timeStart, " \n Checkout date = ", timeEnd, " \n Hasil Row = ", affectedRow)

	if affectedRow == 0 {
		return orders.toCoreAvailable(nameService, timeStart, timeEnd, statusAvailable), nil
	}

	return orders.toCoreNotAvailable(nameService, timeStart, timeEnd, statusNotvalable), nil
}