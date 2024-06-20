package handlers

import (
	basemodels "GoTransact/apps/base"
	"GoTransact/apps/transaction/functions"
	log "GoTransact/settings"
	"fmt"
	"html/template"
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

	_, message, data := functions.ConfirmPayment(transactionIdStr, statusStr)

	// Convert data to map to extract transaction details

	// Create a map for template data
	tmplData := map[string]interface{}{
		"TransactionID": transactionIdStr,
		"Amount":        data["Amount"],
		"Message":       message,
	}
	fmt.Println("---------------", message, "---------------")
	// Select the template based on the message
	var tmpl *template.Template
	var err error

	if message == "Transaction successful" {
		tmpl, err = template.ParseFiles("/home/trellis/Sourabh/GoTransact/Backend/apps/transaction/templates/payment_success.html")
	} else if message == "Transaction Canceled" {
		tmpl, err = template.ParseFiles("/home/trellis/Sourabh/GoTransact/Backend/apps/transaction/templates/payment_fail.html")
	} else {
		c.JSON(http.StatusInternalServerError, basemodels.Response{
			Status:  http.StatusInternalServerError,
			Message: "Unknown transaction status",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, basemodels.Response{
			Status:  http.StatusInternalServerError,
			Message: "Template parsing error",
		})
		return
	}

	// Render the template
	c.Writer.Header().Set("Content-Type", "text/html")
	tmpl.Execute(c.Writer, tmplData)
}
