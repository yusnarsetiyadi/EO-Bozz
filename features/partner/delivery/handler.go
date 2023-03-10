package delivery

import (
	cfg "capstone-alta1/config"
	"capstone-alta1/features/partner"
	"capstone-alta1/middlewares"
	"capstone-alta1/utils/helper"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type PartnerDelivery struct {
	partnerService partner.ServiceInterface
}

func New(service partner.ServiceInterface, e *echo.Echo) {
	handler := &PartnerDelivery{
		partnerService: service,
	}

	e.GET("/partners", handler.GetAll, middlewares.JWTMiddleware())
	e.GET("/partners/:id", handler.GetById, middlewares.JWTMiddleware())
	e.POST("/partners", handler.Create)
	e.PUT("/partners", handler.Update, middlewares.JWTMiddleware())
	e.DELETE("/partners", handler.Delete, middlewares.JWTMiddleware())
	e.PUT("/partners/orders/:id/confirm", handler.ConfirmOrder, middlewares.JWTMiddleware())
	e.GET("/partners/:id/services", handler.GetPartnerServices)
	e.GET("/partners/orders", handler.GetPartnerOrders, middlewares.JWTMiddleware())
	e.GET("/partners/additionals", handler.GetPartnerAdditionals, middlewares.JWTMiddleware())
	e.GET("/partners/register", handler.GetPartnerRegisterData, middlewares.JWTMiddleware())
	e.GET("/partners/:id/register", handler.GetPartnerRegisterDataByID, middlewares.JWTMiddleware())
	e.PUT("/partners/verify", handler.VerifyPartner, middlewares.JWTMiddleware())

}

func (delivery *PartnerDelivery) GetAll(c echo.Context) error {
	// userRole := middlewares.ExtractTokenUserRole(c)
	// if userRole != "Admin" {
	// 	return c.JSON(http.StatusUnauthorized, helper.FailedResponse("this action only admin"))
	// }
	query := c.QueryParam("name")
	helper.LogDebug("isi query = ", query)
	results, err := delivery.partnerService.GetAll(query)
	if err != nil {
		if strings.Contains(err.Error(), "Get data success. No data.") {
			return c.JSON(http.StatusOK, helper.SuccessWithDataResponse(err.Error(), results))
		}
		return c.JSON(http.StatusBadRequest, helper.FailedResponse(err.Error()))
	}

	dataResponse := fromListCoreList(results)

	return c.JSON(http.StatusOK, helper.SuccessWithDataResponse("Success read all data.", dataResponse))
}

func (delivery *PartnerDelivery) GetById(c echo.Context) error {
	idParam := c.Param("id")
	id, errConv := strconv.Atoi(idParam)
	if errConv != nil {
		return c.JSON(http.StatusBadRequest, helper.FailedResponse("Error. Id must integer."))
	}
	results, err := delivery.partnerService.GetById(uint(id))
	if err != nil {
		if strings.Contains(err.Error(), "Get data success. No data.") {
			return c.JSON(http.StatusOK, helper.SuccessWithDataResponse(err.Error(), results))
		}
		return c.JSON(http.StatusBadRequest, helper.FailedResponse(err.Error()))
	}

	dataResponse := fromCore(results)

	return c.JSON(http.StatusOK, helper.SuccessWithDataResponse("Success read data.", dataResponse))
}

func (delivery *PartnerDelivery) Create(c echo.Context) error {
	userInput := PartnerRequest{}
	errBind := c.Bind(&userInput) // menangkap data yg dikirim dari req body dan disimpan ke variabel
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, helper.FailedResponse("Error binding data. "+errBind.Error()))
	}

	dataCore := toCore(userInput)
	err := delivery.partnerService.Create(dataCore, c)
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

func (delivery *PartnerDelivery) Update(c echo.Context) error {
	partnerID := uint(middlewares.ExtractTokenPartnerID(c))
	userID := uint(middlewares.ExtractTokenUserId(c))
	userInput := PartnerRequest{}
	errBind := c.Bind(&userInput) // menangkap data yg dikirim dari req body dan disimpan ke variabel
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, helper.FailedResponse("Error binding data. "+errBind.Error()))
	}

	dataCore := toCore(userInput)
	err := delivery.partnerService.Update(dataCore, partnerID, userID, c)
	if err != nil {
		if strings.Contains(err.Error(), "Error:Field validation") {
			return c.JSON(http.StatusBadRequest, helper.FailedResponse("Some field cannot Empty. Details : "+err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, helper.FailedResponse("Failed update data. "+err.Error()))
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Success update data."))
}

func (delivery *PartnerDelivery) Delete(c echo.Context) error {
	partnerID := uint(middlewares.ExtractTokenPartnerID(c))
	userID := uint(middlewares.ExtractTokenUserId(c))
	err := delivery.partnerService.Delete(partnerID, userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.FailedResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Success delete data."))
}

func (delivery *PartnerDelivery) ConfirmOrder(c echo.Context) error {
	idParam := c.Param("id")
	orderID, errConv := strconv.Atoi(idParam)
	if errConv != nil {
		return c.JSON(http.StatusBadRequest, helper.FailedResponse("Error. Id must integer."))
	}

	if orderID < 1 {
		helper.LogDebug("Partner - handler - ConfirmOrder | order id  = ", orderID)
		return c.JSON(http.StatusBadRequest, helper.FailedResponse("Failed load order id from parameter, please check again."))
	}

	partnerID := middlewares.ExtractTokenPartnerID(c)
	helper.LogDebug("Partner - handler - ConfirmOrder | partner id = ", partnerID)
	if partnerID < 1 {
		helper.LogDebug("Partner - handler - ConfirmOrder | validasi id. id = ", partnerID)
		return c.JSON(http.StatusBadRequest, helper.FailedResponse("Failed load id from JWT token, please check again."))
	}

	// userInput := "Order Confirmed"

	// helper.LogDebug("Partner - handler - ConfirmOrder | order id = ", orderID)
	// helper.LogDebug("Partner - handler - ConfirmOrder | userInput = ", userInput)
	// helper.LogDebug("Partner - handler - ConfirmOrder | partner id = ", partnerID)
	// helper.LogDebug("Partner - handler - ConfirmOrder | mau masuk proses =")

	err := delivery.partnerService.UpdateOrderConfirmStatus(uint(orderID), uint(partnerID))
	if err != nil {
		helper.LogDebug(err.Error())
		if strings.Contains(err.Error(), "Order data no need partner confirmation.") {
			return c.JSON(http.StatusBadRequest, helper.FailedResponse("Order data no need partner confirmation."))
		}
		return c.JSON(http.StatusInternalServerError, helper.FailedResponse("Failed to update partner status. Please try again."))
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Success update data status partner to Verified."))
}

func (delivery *PartnerDelivery) GetPartnerServices(c echo.Context) error {
	idParam := c.Param("id")
	partnerID, errConv := strconv.Atoi(idParam)
	if errConv != nil {
		return c.JSON(http.StatusBadRequest, helper.FailedResponse("Error. Id must integer."))
	}

	if partnerID < 1 {
		helper.LogDebug("Partner - handler - GetPartnerServices | partner id  = ", partnerID)
		return c.JSON(http.StatusBadRequest, helper.FailedResponse("Failed load id from parameter, please check again."))
	}

	results, err := delivery.partnerService.GetServices(uint(partnerID))
	if err != nil {
		if strings.Contains(err.Error(), "Get data success. No data.") {
			return c.JSON(http.StatusOK, helper.SuccessWithDataResponse(err.Error(), results))
		}
		return c.JSON(http.StatusBadRequest, helper.FailedResponse(err.Error()))
	}

	dataResponse := fromCoreServiceList(results)

	return c.JSON(http.StatusOK, helper.SuccessWithDataResponse("Success read all data.", dataResponse))
}

func (delivery *PartnerDelivery) GetPartnerOrders(c echo.Context) error {
	// userRole := middlewares.ExtractTokenUserRole(c)
	// if userRole != "Admin" {
	// 	return c.JSON(http.StatusUnauthorized, helper.FailedResponse("this action only admin"))
	// }

	partnerID := middlewares.ExtractTokenPartnerID(c)

	helper.LogDebug("Partner - handler - get partnerorders | partner id = ", partnerID)

	if partnerID < 1 {
		helper.LogDebug("Partner - handler - get partnerorders | validasi id. id = ", partnerID)
		return c.JSON(http.StatusBadRequest, helper.FailedResponse("Failed load id from JWT token, please check again."))
	}

	helper.LogDebug("Partner - handler - get partnerorders | mau mamsuk proses =")

	results, err := delivery.partnerService.GetOrders(uint(partnerID))
	if err != nil {
		if strings.Contains(err.Error(), "Get data success. No data.") {
			return c.JSON(http.StatusOK, helper.SuccessWithDataResponse(err.Error(), results))
		}
		return c.JSON(http.StatusBadRequest, helper.FailedResponse(err.Error()))
	}

	dataResponse := fromOrderCoreList(results)

	return c.JSON(http.StatusOK, helper.SuccessWithDataResponse("Success read all data.", dataResponse))
}

func (delivery *PartnerDelivery) GetPartnerAdditionals(c echo.Context) error {
	// userRole := middlewares.ExtractTokenUserRole(c)
	// if userRole != "Admin" {
	// 	return c.JSON(http.StatusUnauthorized, helper.FailedResponse("this action only admin"))
	// }
	partnerID := middlewares.ExtractTokenPartnerID(c)

	helper.LogDebug("Partner - handler - get partner additionals | partner id = ", partnerID)

	if partnerID < 1 {
		helper.LogDebug("Partner - handler - get  partner additionals | validasi id. id = ", partnerID)
		return c.JSON(http.StatusBadRequest, helper.FailedResponse("Failed load id from JWT token, please check again."))
	}

	helper.LogDebug("Partner - handler - get  partner additionals | mau mamsuk proses =")

	results, err := delivery.partnerService.GetAdditionals(uint(partnerID))
	if err != nil {
		if strings.Contains(err.Error(), "Get data success. No data.") {
			return c.JSON(http.StatusOK, helper.SuccessWithDataResponse(err.Error(), results))
		}
		return c.JSON(http.StatusBadRequest, helper.FailedResponse(err.Error()))
	}

	dataResponse := fromAdditionalCoreList(results)

	return c.JSON(http.StatusOK, helper.SuccessWithDataResponse("Success read all data.", dataResponse))
}

func (delivery *PartnerDelivery) GetPartnerRegisterData(c echo.Context) error {
	// userRole := middlewares.ExtractTokenUserRole(c)
	// if userRole != "Admin" {
	// 	return c.JSON(http.StatusUnauthorized, helper.FailedResponse("this action only admin"))
	// }
	queryCompanyName := c.QueryParam("company_name")
	queryPICName := c.QueryParam("pic_name")
	queryPartnerStatus := c.QueryParam("partner_status")
	helper.LogDebug("Partner Handler - GetPartnerRegisterData | queryCompanyName = ", queryCompanyName)
	helper.LogDebug("Partner Handler - GetPartnerRegisterData | queryPICName = ", queryPICName)
	helper.LogDebug("Partner Handler - GetPartnerRegisterData | queryPartnerStatus = ", queryPartnerStatus)

	results, err := delivery.partnerService.GetPartnerRegisterData(queryCompanyName, queryPICName, queryPartnerStatus)
	if err != nil {
		if strings.Contains(err.Error(), "Get data success. No data.") {
			return c.JSON(http.StatusOK, helper.SuccessWithDataResponse(err.Error(), results))
		}
		return c.JSON(http.StatusBadRequest, helper.FailedResponse(err.Error()))
	}

	dataResponse := fromListCoreList(results)

	return c.JSON(http.StatusOK, helper.SuccessWithDataResponse("Success read all data.", dataResponse))
}

func (delivery *PartnerDelivery) GetPartnerRegisterDataByID(c echo.Context) error {
	// userRole := middlewares.ExtractTokenUserRole(c)
	// if userRole != "Admin" {
	// 	return c.JSON(http.StatusUnauthorized, helper.FailedResponse("this action only admin"))
	// }
	idParam := c.Param("id")
	id, errConv := strconv.Atoi(idParam)
	if errConv != nil {
		return c.JSON(http.StatusBadRequest, helper.FailedResponse("Error. Id must integer."))
	}
	results, err := delivery.partnerService.GetPartnerRegisterDataByID(uint(id))
	if err != nil {
		if strings.Contains(err.Error(), "Get data success. No data.") {
			return c.JSON(http.StatusOK, helper.SuccessWithDataResponse(err.Error(), results))
		}
		return c.JSON(http.StatusBadRequest, helper.FailedResponse(err.Error()))
	}

	dataResponse := fromCore(results)

	return c.JSON(http.StatusOK, helper.SuccessWithDataResponse("Success read data.", dataResponse))
}

func (delivery *PartnerDelivery) VerifyPartner(c echo.Context) error {

	userInput := VerifyPartnerRequest{}
	errBind := c.Bind(&userInput) // menangkap data yg dikirim dari req body dan disimpan ke variabel
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, helper.FailedResponse("Error binding data. "+errBind.Error()))
	}
	if (userInput.Status == cfg.PARTNER_VERIFICATION_STATUS_NOT_VERIFIED) || (userInput.Status == cfg.PARTNER_VERIFICATION_STATUS_REVISION) || (userInput.Status == cfg.PARTNER_VERIFICATION_STATUS_VERIFIED) {
		helper.LogDebug("Partner - handler - VerifyPartner | userInput = ", userInput)
		helper.LogDebug("Partner - handler - VerifyPartner | partner id = ", userInput.PartnerID)

		if userInput.PartnerID < 1 {
			helper.LogDebug("Partner - handler - VerifyPartner | validasi id. id = ", userInput.PartnerID)
			return c.JSON(http.StatusBadRequest, helper.FailedResponse("Failed load id from JWT token, please check again."))
		}

		helper.LogDebug("Partner - handler - VerifyPartner | mau masuk proses =")

		err := delivery.partnerService.UpdatePartnerVerifyStatus(userInput.VerifyLog, userInput.Status, uint(userInput.PartnerID))
		if err != nil {
			helper.LogDebug(err.Error())
			if strings.Contains(err.Error(), "Failed. Cannot update status that already Verified.") {
				return c.JSON(http.StatusBadRequest, helper.FailedResponse("Failed. Cannot update status that already Verified."))
			}
			return c.JSON(http.StatusBadRequest, helper.FailedResponse("Failed to update partner status. Please try again."))
		}

		return c.JSON(http.StatusOK, helper.SuccessResponse("Success update data status partner to "+userInput.Status+"."))

	} else {
		return c.JSON(http.StatusBadRequest, helper.FailedResponse("Failed. Incorrect status. Please input either Not Verified, Revision or Verified."))
	}
}
