package delivery

import (
	"capstone-alta1/features/client"
	"capstone-alta1/utils/helper"
)

type ClientResponse struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	Role            string `json:"role"`
	Gender          string `json:"gender"`
	Address         string `json:"address"`
	City            string `json:"city"`
	Phone           string `json:"phone"`
	ClientImageFile string `json:"client_image_file"`
	UserID          uint   `json:"user_id"`
}

type ClientOrderResponse struct {
	ID            uint
	EventName     string
	StartDate     string
	EndDate       string
	EventLocation string
	ServiceName   string
	GrossAmmount  uint
	OrderStatus   string
	ServiceID     uint
	ClientID      uint
}

func fromCore(dataCore client.Core) ClientResponse {
	return ClientResponse{
		ID:              dataCore.ID,
		Name:            dataCore.User.Name,
		Email:           dataCore.User.Email,
		Role:            dataCore.User.Role,
		Gender:          dataCore.Gender,
		Address:         dataCore.Address,
		City:            dataCore.City,
		Phone:           dataCore.Phone,
		ClientImageFile: dataCore.ClientImageFile,
		UserID:          dataCore.User.ID,
	}
}

func fromCoreOrder(dataCore client.Order) ClientOrderResponse {
	return ClientOrderResponse{
		ID:            dataCore.ID,
		EventName:     dataCore.EventName,
		StartDate:     helper.GetDateFormated(dataCore.StartDate),
		EndDate:       helper.GetDateFormated(dataCore.EndDate),
		EventLocation: dataCore.EventLocation,
		ServiceName:   dataCore.ServiceName,
		GrossAmmount:  dataCore.GrossAmmount,
		OrderStatus:   dataCore.OrderStatus,
		ServiceID:     dataCore.ServiceID,
		ClientID:      dataCore.ClientID,
	}
}

func fromCoreList(dataCore []client.Core) []ClientResponse {
	var dataResponse []ClientResponse
	for _, v := range dataCore {
		dataResponse = append(dataResponse, fromCore(v))
	}
	return dataResponse
}

func fromCoreListOrder(dataCore []client.Order) []ClientOrderResponse {
	var dataResponse []ClientOrderResponse
	for _, v := range dataCore {
		dataResponse = append(dataResponse, fromCoreOrder(v))
	}
	return dataResponse
}
