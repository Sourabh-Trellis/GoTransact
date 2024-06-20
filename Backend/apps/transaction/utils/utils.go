package utils

import (
	accountmodels "GoTransact/apps/accounts/models"
	transactionmodels "GoTransact/apps/transaction/models"
	log "GoTransact/settings"
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	gomail "gopkg.in/mail.v2"
)

type TemplateData struct {
	Username     string
	TrasactionID uuid.UUID
	Amount       float64
	ConfirmURL   string
	CancelURL    string
	DateTime     time.Time
}

func SendMail(user accountmodels.User, request transactionmodels.TransactionRequest) {

	log.InfoLogger.WithFields(logrus.Fields{
		"email": user.Email,
		"id":    user.Internal_id,
	}).Info("Attempted to send confirm payment mail")

	fmt.Println("start of mail")
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", "sourabhtrellis@gmail.com")

	// Set E-Mail receivers
	// m.SetHeader("To", user.Email)
	m.SetHeader("To", user.Email)

	// Set E-Mail subject
	m.SetHeader("Subject", "Payment Confirmation Required")

	// Parse the HTML template
	tmpl, err := template.ParseFiles("/home/trellis/Sourabh/GoTransact/Backend/apps/transaction/utils/email_template.html")
	if err != nil {
		fmt.Printf("Error parsing email template: %s", err)
	}

	// Create a buffer to hold the executed template
	var body bytes.Buffer

	baseURL := "http://localhost:8080/api/confirm-payment" // Replace with your actual domain and endpoint
	params := url.Values{}
	params.Add("transaction_id", request.Internal_id.String())
	params.Add("status", "true")
	ConfirmActionURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	baseURL = "http://localhost:8080/api/confirm-payment" // Replace with your actual domain and endpoint
	params = url.Values{}
	params.Add("transaction_id", request.Internal_id.String())
	params.Add("status", "false")
	CancelActionURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Execute the template with the data
	TemplateData := TemplateData{
		Username:     user.FirstName,
		TrasactionID: request.Internal_id,
		Amount:       request.Amount,
		ConfirmURL:   ConfirmActionURL,
		CancelURL:    CancelActionURL,
	}
	fmt.Println(TemplateData)
	if err := tmpl.Execute(&body, TemplateData); err != nil {
		fmt.Printf("Error executing email template: %s", err)
	}

	fmt.Println()
	// Set E-Mail body as HTML
	m.SetBody("text/html", body.String())

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, "sourabhtrellis@gmail.com", "nmvx vzro ehqo xwpd")

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		log.ErrorLogger.WithFields(logrus.Fields{
			"error": err.Error(),
			"email": user.Email,
		}).Error("Error sending confirm payment mail")
		panic(err)
	}
	log.InfoLogger.WithFields(logrus.Fields{
		"email": user.Email,
		"id":    user.Internal_id,
	}).Info("comfirmation mail sent")
}
