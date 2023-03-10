package repository

import (
	client "capstone-alta1/features/client"

	"time"

	"gorm.io/gorm"
)

// struct gorm model
type Client struct {
	gorm.Model
	Gender          string
	Address         string
	City            string
	Phone           string
	ClientImageFile string
	UserID          uint
	User            User
	Order           []Order
	Review          []Review
}

type User struct {
	gorm.Model
	Name     string
	Email    string
	Password string
	Role     string
}

type Order struct {
	gorm.Model
	EventName     string
	StartDate     time.Time
	EndDate       time.Time
	EventLocation string
	ServiceName   string
	GrossAmmount  uint
	OrderStatus   string
	ServiceID     uint
	ClientID      uint
}

type Review struct {
	gorm.Model
	Review    string
	Rating    float64
	OrderID   uint
	ClientID  uint
	ServiceID uint
}

// mapping

// mengubah struct core ke struct model gorm
func fromCore(dataCore client.Core) Client {
	clientGorm := Client{
		User: User{
			Name:     dataCore.User.Name,
			Email:    dataCore.User.Email,
			Password: dataCore.User.Password,
			Role:     dataCore.User.Role,
		},
		Gender:          dataCore.Gender,
		Address:         dataCore.Address,
		City:            dataCore.City,
		Phone:           dataCore.Phone,
		ClientImageFile: dataCore.ClientImageFile,
		UserID:          dataCore.User.ID,
	}
	return clientGorm
}

func fromOrder(dataCore client.Order) Order {
	orderGorm := Order{
		OrderStatus: dataCore.OrderStatus,
	}
	return orderGorm
}

// mengubah struct model gorm ke struct core
func (dataModel *Client) toCore() client.Core {
	return client.Core{
		User: client.User{
			ID:       dataModel.User.ID,
			Name:     dataModel.User.Name,
			Email:    dataModel.User.Email,
			Password: dataModel.User.Password,
			Role:     dataModel.User.Role,
		},
		Gender:          dataModel.Gender,
		Address:         dataModel.Address,
		City:            dataModel.City,
		Phone:           dataModel.Phone,
		ClientImageFile: dataModel.ClientImageFile,
		UserID:          dataModel.User.ID,
		ID:              dataModel.ID,
	}
}

func (data *Order) toCoreOrder() client.Order {
	return client.Order{
		ID:            data.ID,
		EventName:     data.EventName,
		StartDate:     data.StartDate,
		EndDate:       data.EndDate,
		EventLocation: data.EventLocation,
		ServiceName:   data.ServiceName,
		GrossAmmount:  data.GrossAmmount,
		OrderStatus:   data.OrderStatus,
		ServiceID:     data.ServiceID,
		ClientID:      data.ClientID,
	}
}

// mengubah slice struct model gorm ke slice struct core
func toCoreList(dataModel []Client) []client.Core {
	var dataCore []client.Core
	for _, v := range dataModel {
		dataCore = append(dataCore, v.toCore())
	}
	return dataCore
}

func toCoreListOrder(dataModel []Order) []client.Order {
	var dataCore []client.Order
	for _, v := range dataModel {
		dataCore = append(dataCore, v.toCoreOrder())
	}
	return dataCore
}
