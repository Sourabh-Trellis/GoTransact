package handlers

import (
	basemodels "GoTransact/apps/base"
	"GoTransact/apps/transaction/functions"
	log "GoTransact/settings"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func ConfirmPayment(c *gin.Context) {

	log.InfoLogger.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"url":    c.Request.URL.String(),
	}).Info("confirm payment Request received")

	transactionIdStr := c.Query("transaction_id")
	statusStr := c.Query("status")

	status, message, data := functions.ConfirmPayment(transactionIdStr, statusStr)

	c.JSON(http.StatusOK, basemodels.Response{
		Status:  status,
		Message: message,
		Data:    data,
	})
}
