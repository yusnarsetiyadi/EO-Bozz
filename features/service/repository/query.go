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

func (repo *serviceRepository) GetAll(queryName, queryCategory, queryCity, queryMinPrice, queryMaxPrice string) (data []_service.Core, err error) {
	var results []Service
	maxInt, _ := strconv.Atoi(queryMaxPrice)
	minInt, _ := strconv.Atoi(queryMinPrice)
	tx := repo.db.Where("service_name LIKE ?", "%"+queryName+"%").Where(&Service{City: queryCity, ServiceCategory: queryCategory}).Find(&results)
	// tx := repo.db.Where("service_name LIKE ?", "%"+queryName+"%").Where("OR service_category LIKE ? OR city LIKE ? OR service_price BETWEEN ? AND ?", "%"+queryName+"%", "%"+queryCategory+"%", "%"+queryCity+"%", uint(minInt), uint(maxInt)).Find(&results)

	if tx.Error != nil {
		helper.LogDebug("Order - query - GetAll | Error Query Find = ", tx.Error)
		return nil, tx.Error
	}

	helper.LogDebug("Order - query - GetAll | Result data : ", results)
	helper.LogDebug("Order - query - GetAll | Row Affected : ", tx.RowsAffected)
	if tx.RowsAffected == 0 {
		return nil, tx.Error
	}

	helper.LogDebug("Order - query - GetAll | Min : ", minInt, " Max ", maxInt)
	var ty *gorm.DB
	if maxInt > minInt {
		ty = repo.db.Where("service_price >= ? AND service_price <= ?", uint(minInt), uint(maxInt)).Find(&results)
	} else {
		ty = repo.db.Where("service_price >= ?", uint(minInt)).Find(&results)
	}

	if tx.Error != nil {
		helper.LogDebug("Order - query - GetAll | Error Query Find = ", ty.Error)
		return nil, ty.Error
	}

	helper.LogDebug("Order - query - GetAll | Result data 2: ", results)
	helper.LogDebug("Order - query - GetAll | Row Affected 2 : ", ty.RowsAffected)
	if ty.RowsAffected == 0 {
		return nil, ty.Error
	}

	var dataCore = toCoreListGetAll(results)
	return dataCore, nil

}

func (repo *serviceRepository) GetById(id uint) (data _service.ServiceDetailJoinPartner, err error) {
	var serviceDetailJoinPartner ServiceDetailJoinPartner

	tx := repo.db.Raw("SELECT `services`.*, `partners`.`company_name`, `partners`.`company_phone`, `partners`.`company_city`, `partners`.`company_image_file`, `partners`.`company_address`, `partners`.`link_website`, `partners`.`verification_status`, `partners`.`user_id` FROM `services` JOIN `partners` ON `services`.`partner_id` = `partners`.`id` WHERE `services`.`id` = ?", id).Scan(&serviceDetailJoinPartner)

	if tx.Error != nil {
		helper.LogDebug("Order - query - GetAll | Error Query Find = ", tx.Error)
		return _service.ServiceDetailJoinPartner{}, tx.Error
	}

	helper.LogDebug("Order - query - GetAll | Result data : ", serviceDetailJoinPartner)
	helper.LogDebug("Order - query - GetAll | Row Affected query get additional data : ", tx.RowsAffected)
	if tx.RowsAffected == 0 {
		return _service.ServiceDetailJoinPartner{}, tx.Error
	}

	var dataCore = serviceDetailJoinPartner.toCoreGetById()
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

func (repo *serviceRepository) GetServiceAdditionalById(serviceId uint) (data []_service.JoinServiceAdditional, err error) {
	var modelData []JoinServiceAdditional

	// tx := repo.db.Model(&serviceAdditional).Where("service_id = ?", serviceId).Find(&additional, serviceAdditionalId.AdditionalID)
	tx := repo.db.Raw("SELECT `service_additionals`.`id` AS `service_additional_id`, `additionals`.`additional_name`, `additionals`.`additional_price`, `services`.`service_name`, `service_additionals`.`service_id`, `service_additionals`.`additional_id`, `services`.`partner_id`  FROM `additionals` JOIN `service_additionals` ON `additionals`.`id` = `service_additionals`.`additional_id` JOIN `services` ON `services`.`id` = `service_additionals`.`service_id` WHERE `services`.`id` = ?", serviceId).Scan(&modelData)

	helper.LogDebug("service-query-GetAdditional | Data ServiceAdditional Model Data :", modelData)

	if tx.Error != nil {
		helper.LogDebug("service-query-GetAdditional | Error execute query. Error :", tx.Error)
		return data, tx.Error
	}

	helper.LogDebug("service-query-GetAdditional  | Row Affected : ", tx.RowsAffected)
	if tx.RowsAffected == 0 {
		return data, tx.Error
	}

	var dataCore = toCoreJoinServiceAdditionalsList(modelData)
	helper.LogDebug("service-query-GetAdditional | Data ServiceAdditional DataCore :", dataCore)
	return dataCore, nil
}

func (repo *serviceRepository) AddAdditionalToService(input []_service.ServiceAdditional) error {
	serviceadditionalGorm := fromCoreServiceAdditionalList(input)

	helper.LogDebug("service-query-AddAdditionalToService | Service Additionals Data :", helper.ConvToJson(serviceadditionalGorm))
	tx := repo.db.Save(&serviceadditionalGorm) // proses insert data
	if tx.Error != nil {
		helper.LogDebug("service-query-AddAdditionalToService | Error execute query. Error :", tx.Error)
		return tx.Error
	}

	helper.LogDebug("service-query-AddAdditionalToService | Row Affected : ", tx.RowsAffected)
	if tx.RowsAffected == 0 {
		return errors.New("insert failed")
	}
	return nil
}

func (repo *serviceRepository) GetReviewById(serviceId uint) (data []_service.Review, err error) {
	var clientreview []Review

	tx := repo.db.Where("service_id = ?", serviceId).Find(&clientreview)

	if tx.Error != nil {
		helper.LogDebug("service-query-Getreview | Error execute query. Error :", tx.Error)
		return data, tx.Error
	}

	helper.LogDebug("service-query-Getreview   | Row Affected : ", tx.RowsAffected)
	if tx.RowsAffected == 0 {
		return data, tx.Error
	}

	var dataCore = toCoreListReview(clientreview)
	return dataCore, nil
}

func (repo *serviceRepository) GetDiscussionById(serviceId uint) (data []_service.Discussion, err error) {
	var clientdiscussion []Discussion

	tx := repo.db.Where("service_id = ?", serviceId).Find(&clientdiscussion)

	if tx.Error != nil {
		helper.LogDebug("service-query-Getdiscussion | Error execute query. Error :", tx.Error)
		return data, tx.Error
	}

	helper.LogDebug("service-query-Getdiscussion  | Row Affected : ", tx.RowsAffected)
	if tx.RowsAffected == 0 {
		return data, tx.Error
	}

	var dataCore = toCoreListDiscussion(clientdiscussion)
	return dataCore, nil
}

func (repo *serviceRepository) CheckAvailability(serviceId uint, queryStart, queryEnd time.Time) (data _service.Order, err error) {
	//check available
	var services []Service
	var service Service
	queryBuilder := fmt.Sprintf("SELECT * FROM orders WHERE service_id = %d AND (('%s' BETWEEN start_date AND end_date) OR ('%s' BETWEEN start_date AND end_date));", serviceId, queryStart, queryEnd)
	tx := repo.db.Raw(queryBuilder).Find(&services)

	//get data service
	yx := repo.db.First(&service, serviceId)
	if yx.Error != nil {
		return _service.Order{}, err
	}

	//create return
	var orders Order
	serviceName := service.ServiceName
	statusAvailable := "Available"
	statusNotvalable := "Not Available"

	if tx.Error != nil {
		return orders.toCoreNotAvailable(serviceName, queryStart, queryEnd, statusNotvalable), tx.Error
	}

	affectedRow := tx.RowsAffected
	fmt.Println("\n\nHasil check availbility, \n Checkin date = ", queryStart, " \n Checkout date = ", queryEnd, " \n Hasil Row = ", affectedRow)

	if affectedRow == 0 {
		return orders.toCoreAvailable(serviceName, queryStart, queryEnd, statusAvailable), nil
	}

	return orders.toCoreNotAvailable(serviceName, queryStart, queryEnd, statusNotvalable), nil
}
