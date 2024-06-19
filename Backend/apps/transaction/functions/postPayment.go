package functions

import (
	accountmodels "GoTransact/apps/accounts/models"
	transactionmodels "GoTransact/apps/transaction/models"
	"GoTransact/apps/transaction/utils"
	valid "GoTransact/apps/transaction/validators"
	"GoTransact/pkg/db"
	log "GoTransact/settings"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
)

func PostPayment(Postpaymentinput utils.PostPaymentInput, user accountmodels.User) (int, string, map[string]interface{}) {

	log.InfoLogger.WithFields(logrus.Fields{}).Info("Attempted to create transaction request with email ", user.Email, " id ", user.Internal_id)

	if err := valid.GetValidator().Struct(Postpaymentinput); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errors := make(map[string]string)
		for _, fieldErr := range validationErrors {
			fieldName := fieldErr.Field()
			tag := fieldErr.Tag()
			errors[fieldName] = valid.CustomErrorMessages[tag]
		}
		return http.StatusBadRequest, "error while validating", map[string]interface{}{}
	}

	floatAmount, _ := strconv.ParseFloat(Postpaymentinput.Amount, 64)

	var gateway transactionmodels.Payment_Gateway
	if err := db.DB.Where("slug = ?", "card").First(&gateway).Error; err != nil {
		return http.StatusBadRequest, "invalid payment type", map[string]interface{}{}
	}

	TransactionRequest := transactionmodels.TransactionRequest{
		UserID:             user.ID,
		Status:             transactionmodels.StatusProcessing,
		Description:        Postpaymentinput.Description,
		Amount:             floatAmount,
		Payment_Gateway_id: gateway.ID,
	}

	if err := db.DB.Create(&TransactionRequest).Error; err != nil {
		log.ErrorLogger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Error creating record in transaction request transaction-id=", TransactionRequest.Internal_id)
		return http.StatusInternalServerError, "internal server error", map[string]interface{}{}
	}
	log.InfoLogger.WithFields(logrus.Fields{}).Info("created record in transaction request with email ", user.Email, " id ", user.Internal_id)

	TransactionHistory := transactionmodels.TransactionHistory{
		TransactionID: TransactionRequest.ID,
		Status:        TransactionRequest.Status,
		Description:   TransactionRequest.Description,
		Amount:        TransactionRequest.Amount,
	}

	if err := db.DB.Create(&TransactionHistory).Error; err != nil {
		log.ErrorLogger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Error creating record in transaction history")
		return http.StatusInternalServerError, "internal server error", map[string]interface{}{}
	}

	log.InfoLogger.WithFields(logrus.Fields{}).Info("Created record in transaction history with email ", user.Email, " id ", user.Internal_id)

	go utils.SendMail(user, TransactionRequest)

	return http.StatusOK, "success", map[string]interface{}{"transaction ID": TransactionRequest.Internal_id}
}
