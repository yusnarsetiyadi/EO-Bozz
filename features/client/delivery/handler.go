package delivery

import (
	cfg "capstone-alta1/config"
	"capstone-alta1/features/client"
	"capstone-alta1/middlewares"
	"capstone-alta1/utils/helper"
	"fmt"

	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type ClientDelivery struct {
	clientService client.ServiceInterface
}

func New(service client.ServiceInterface, e *echo.Echo) {
	handler := &ClientDelivery{
		clientService: service,
	}

	e.GET("/clients", handler.GetAll, middlewares.JWTMiddleware())
	e.GET("/clients/:id", handler.GetById, middlewares.JWTMiddleware())
	e.POST("/clients", handler.Create)
	e.PUT("/clients", handler.Update, middlewares.JWTMiddleware())
	e.DELETE("/clients", handler.Delete, middlewares.JWTMiddleware())
	e.GET("/clients/orders", handler.GetOrderById, middlewares.JWTMiddleware())
	e.PUT("/clients/orders/:id/complete", handler.UpdateCompleteOrder, middlewares.JWTMiddleware())

}

func (delivery *ClientDelivery) GetAll(c echo.Context) error {
	query := c.QueryParam("name")
	helper.LogDebug("isi query = ", query)
	results, err := delivery.clientService.GetAll(query)
	if err != nil {
		if strings.Contains(err.Error(), "Get data success. No data.") {
			return c.JSON(http.StatusOK, helper.SuccessWithDataResponse(err.Error(), results))
		}
		return c.JSON(http.StatusBadRequest, helper.FailedResponse(err.Error()))
	}

	dataResponse := fromCoreList(results)

	return c.JSON(http.StatusOK, helper.SuccessWithDataResponse("Success read all data", dataResponse))
}

func (delivery *ClientDelivery) GetById(c echo.Context) error {
	idParam := c.Param("id")
	id, errConv := strconv.Atoi(idParam)
	if errConv != nil {
		return c.JSON(http.StatusBadRequest, helper.FailedResponse("Error. Id must integer."))
	}
	results, err := delivery.clientService.GetById(uint(id))
	if err != nil {
		if strings.Contains(err.Error(), "Get data success. No data.") {
			return c.JSON(http.StatusOK, helper.SuccessWithDataResponse(err.Error(), results))
		}
		return c.JSON(http.StatusBadRequest, helper.FailedResponse(err.Error()))
	}

	dataResponse := fromCore(results)

	return c.JSON(http.StatusOK, helper.SuccessWithDataResponse("Success read data.", dataResponse))
}

func (delivery *ClientDelivery) Create(c echo.Context) error {
	userInput := ClientRequest{}
	errBind := c.Bind(&userInput) // menangkap data yg dikirim dari req body dan disimpan ke variabel
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, helper.FailedResponse("Error binding data. "+errBind.Error()))
	}

	dataCore := toCore(userInput)
	err := delivery.clientService.Create(dataCore, c)
	if err != nil {
		if strings.Contains(err.Error(), "Error:Field validation") {
			return c.JSON(http.StatusBadRequest, helper.FailedResponse("Some field cannot empty. Details : "+err.Error()))
		}
		if strings.Contains(err.Error(), "Please pick another email.") {
			return c.JSON(http.StatusBadRequest, helper.FailedResponse("Failed insert data "+err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, helper.FailedResponse("Failed insert data. "+err.Error()))
	}
	return c.JSON(http.StatusCreated, helper.SuccessResponse("Success create data"))
}

func (delivery *ClientDelivery) Update(c echo.Context) error {
	clientID := uint(middlewares.ExtractTokenClientID(c))
	userID := uint(middlewares.ExtractTokenUserId(c))
	userInput := ClientRequest{}
	errBind := c.Bind(&userInput) // menangkap data yg dikirim dari req body dan disimpan ke variabel
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, helper.FailedResponse("Error binding data. "+errBind.Error()))
	}

	dataCore := toCore(userInput)

	fmt.Println("handler update data ", dataCore)

	err := delivery.clientService.Update(dataCore, clientID, userID, c)
	if err != nil {
		if strings.Contains(err.Error(), "Error:Field validation") {
			return c.JSON(http.StatusBadRequest, helper.FailedResponse("Some field cannot Empty. Details : "+err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, helper.FailedResponse("Failed update data. "+err.Error()))
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Success update data."))
}

func (delivery *ClientDelivery) Delete(c echo.Context) error {
	clientID := uint(middlewares.ExtractTokenClientID(c))
	userId := middlewares.ExtractTokenUserId(c)
	err := delivery.clientService.Delete(clientID, uint(userId))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.FailedResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Success delete data."))
}

func (delivery *ClientDelivery) GetOrderById(c echo.Context) error {
	clientId := middlewares.ExtractTokenClientID(c)
	results, err := delivery.clientService.GetOrderById(uint(clientId))
	if err != nil {
		if strings.Contains(err.Error(), "Get data success. No data.") {
			return c.JSON(http.StatusOK, helper.SuccessWithDataResponse(err.Error(), results))
		}
		return c.JSON(http.StatusBadRequest, helper.FailedResponse(err.Error()))
	}

	dataResponse := fromCoreListOrder(results)

	return c.JSON(http.StatusOK, helper.SuccessWithDataResponse("Success read data.", dataResponse))
}

func (delivery *ClientDelivery) UpdateCompleteOrder(c echo.Context) error {
	clientId := middlewares.ExtractTokenClientID(c)
	idParam := c.Param("id")
	orderId, errConv := strconv.Atoi(idParam)
	if errConv != nil {
		return c.JSON(http.StatusBadRequest, helper.FailedResponse("Error. Id must integer."))
	}
	inputComplete := cfg.ORDER_STATUS_COMPLETE_ORDER
	dataOrder := toOrderStatus(inputComplete)
	err := delivery.clientService.UpdateCompleteOrder(dataOrder, uint(orderId), uint(clientId))
	if err != nil {
		if strings.Contains(err.Error(), "Error:Field validation") {
			return c.JSON(http.StatusBadRequest, helper.FailedResponse("Some field cannot Empty. Details : "+err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, helper.FailedResponse("Failed update data. "+err.Error()))
	}
	return c.JSON(http.StatusOK, helper.SuccessResponse("Success update data."))
}
