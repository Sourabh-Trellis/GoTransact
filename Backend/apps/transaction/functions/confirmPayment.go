package functions

import (
	transactionmodels "GoTransact/apps/transaction/models"
	"GoTransact/pkg/db"
	log "GoTransact/settings"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func ConfirmPayment(transactionIdStr, statusStr string) (int, string, map[string]interface{}) {
	// Convert the string ID to a uuid.UUID
	log.InfoLogger.WithFields(logrus.Fields{
		// "email": user.Email,
		// "id":    user.Internal_id,
	}).Info("Attempted to confirm/cancel payment transaction-id=", transactionIdStr)
	transactionId, err := uuid.Parse(transactionIdStr)
	fmt.Println("parsed", transactionId)
	if err != nil {
		return http.StatusBadRequest, "Invalid transaction ID", map[string]interface{}{}
	}

	var transactionRequest transactionmodels.TransactionRequest
	if err := db.DB.Where("internal_id = ?", transactionId).First(&transactionRequest).Error; err != nil {
		return http.StatusBadRequest, "transaction request not found", map[string]interface{}{}
	}

	var trasactionHistory transactionmodels.TransactionHistory
	trasactionHistory.TransactionID = transactionRequest.ID
	trasactionHistory.Description = transactionRequest.Description
	trasactionHistory.Amount = transactionRequest.Amount

	if strings.EqualFold(statusStr, "true") {

		if err := db.DB.Model(&transactionRequest).Where("id = ?", transactionRequest.ID).Update("status", transactionmodels.StatusSuccess).Error; err != nil {
			log.ErrorLogger.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Error completing payment transaction-id=", transactionRequest.Internal_id)
			return http.StatusInternalServerError, "Failed to confirm the payment", map[string]interface{}{}
		}
		log.InfoLogger.WithFields(logrus.Fields{}).Info("Payment completed transaction-id=", transactionRequest.Internal_id)
		{
			trasactionHistory.Status = transactionmodels.StatusSuccess
		}
	} else {

		if err := db.DB.Model(&transactionRequest).Where("id = ?", transactionRequest.ID).Update("status", transactionmodels.StatusFailed).Error; err != nil {
			log.ErrorLogger.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Error canceling the payment transaction-id=", transactionRequest.Internal_id)
			return http.StatusInternalServerError, "Failed to confirm the payment", map[string]interface{}{}
		}
		log.InfoLogger.WithFields(logrus.Fields{}).Info("Payment canceled transaction-id=", transactionRequest.Internal_id)
		{
			trasactionHistory.Status = transactionmodels.StatusFailed
		}

	}

	if err := db.DB.Create(&trasactionHistory).Error; err != nil {
		log.ErrorLogger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("failed to update transaction history to confirm/cancel transaction-id=", transactionRequest.Internal_id)
		return http.StatusInternalServerError, "Failed to update transaction history", map[string]interface{}{}
	}
	log.InfoLogger.WithFields(logrus.Fields{}).Info("transaction history updated to confirm/cancel transaction-id=", transactionRequest.Internal_id)
	if strings.EqualFold(statusStr, "true") {
		return http.StatusOK, "Transaction successfull", map[string]interface{}{}
	} else {
		return http.StatusOK, "Transaction Canceled", map[string]interface{}{}
	}
}
