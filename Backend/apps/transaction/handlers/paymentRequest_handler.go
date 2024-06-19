package handlers

import (
	accountmodels "GoTransact/apps/accounts/models"
	basemodels "GoTransact/apps/base"
	"GoTransact/apps/transaction/functions"
	"GoTransact/apps/transaction/utils"
	log "GoTransact/settings"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// type PostPaymentInput struct {
// 	CardNumber  string `json:"cardnumber" binding:"required" validate:"card_number" `
// 	ExpiryDate  string `json:"expirydate" binding:"required" validate:"expiry_date" `
// 	Cvv         string `json:"cvv" validate:"cvv" binding:"required"`
// 	Amount      string `json:"amount" binding:"required" validate:"amount"`
// 	Description string `json:"description" `
// }

func PaymentRequest(c *gin.Context) {

	log.InfoLogger.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"url":    c.Request.URL.String(),
	}).Info("Post Payment Request received")

	var Postpaymentinput utils.PostPaymentInput
	if err := c.ShouldBindJSON(&Postpaymentinput); err != nil {
		c.JSON(http.StatusBadRequest, basemodels.Response{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    map[string]interface{}{"data": err.Error()},
		})
		return
	}

	UserFromRequest, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusBadRequest, basemodels.Response{
			Status:  http.StatusBadRequest,
			Message: "User not found in token",
			Data:    map[string]interface{}{},
		})
		return
	}

	user, ok := UserFromRequest.(accountmodels.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assert user type"})
		return
	}

	status, message, data := functions.PostPayment(Postpaymentinput, user)

	c.JSON(status, basemodels.Response{
		Status:  status,
		Message: message,
		Data:    data,
	})
}
