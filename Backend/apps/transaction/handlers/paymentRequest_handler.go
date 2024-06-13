package handlers

import (
	accountmodels "GoTransact/apps/accounts/models"
	basemodels "GoTransact/apps/base"
	transactionmodels "GoTransact/apps/transaction/models"
	"GoTransact/apps/transaction/validators"
	"GoTransact/pkg/db"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PostPaymentInput struct {
	CardNumber  string `json:"cardnumber" binding:"required" validate:"card_number" `
	ExpiryDate  string `json:"expirydate" binding:"required" validate:"expiry_date" `
	Cvv         string `json:"cvv" validate:"cvv" binding:"required"`
	Amount      string `json:"amount" binding:"required" validate:"amount"`
	Description string `json:"description" `
}

func PaymentRequest(c *gin.Context) {
	var Postpaymentinput PostPaymentInput
	if err := c.ShouldBindJSON(&Postpaymentinput); err != nil {
		c.JSON(http.StatusBadRequest, basemodels.Response{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    map[string]interface{}{"data": err.Error()},
		})
		return
	}

	if err := validators.GetValidator().Struct(Postpaymentinput); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errors := make(map[string]string)
		for _, fieldErr := range validationErrors {
			fieldName := fieldErr.Field()
			tag := fieldErr.Tag()
			errors[fieldName] = validators.CustomErrorMessages[tag]
		}

		c.JSON(http.StatusBadRequest, basemodels.Response{
			Status:  http.StatusBadRequest,
			Message: "error while validating",
			Data: map[string]interface{}{
				"validation_errors": errors,
			},
		})
		return
	}
	fmt.Println(c.Keys)
	UserFromRequest, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusBadRequest, basemodels.Response{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    map[string]interface{}{"data": "User not found in token"},
		})
		return
	}

	user, ok := UserFromRequest.(accountmodels.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assert user type"})
		return
	}

	floatAmount, _ := strconv.ParseFloat(Postpaymentinput.Amount, 64)

	var gateway transactionmodels.Payment_Gateway
	if err := db.DB.Where("slug = ?", "card").First(&gateway).Error; err != nil {
		c.JSON(http.StatusBadRequest, basemodels.Response{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    map[string]interface{}{"data": "invalid payment type"},
		})
		return
	}

	TransactionRequest := transactionmodels.TransactionRequest{
		UserID:             user.ID,
		Status:             transactionmodels.StatusProcessing,
		Description:        Postpaymentinput.Description,
		Amount:             floatAmount,
		Payment_Gateway_id: gateway.ID,
	}

	if err := db.DB.Create(&TransactionRequest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, basemodels.Response{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    map[string]interface{}{"data": err.Error()},
		})
		return
	}

	TransactionHistory := transactionmodels.TransactionHistory{
		TransactionID: TransactionRequest.ID,
		Status:        TransactionRequest.Status,
		Description:   TransactionRequest.Description,
		Amount:        TransactionRequest.Amount,
	}

	if err := db.DB.Create(&TransactionHistory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, basemodels.Response{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    map[string]interface{}{"data": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, basemodels.Response{
		Status:  http.StatusOK,
		Message: "success",
		Data:    map[string]interface{}{"data1": TransactionRequest, "data2": TransactionHistory},
	})
}
