package utils

import (
	accountmodels "GoTransact/apps/accounts/models"
	transactionmodels "GoTransact/apps/transaction/models"
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"log"

	gomail "gopkg.in/mail.v2"
)

type TemplateData struct {
	Username string
	Amount   float64
}

func SendMail(user accountmodels.User, request transactionmodels.TransactionRequest) {
	fmt.Println("start of mail")
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", "sourabhtrellis@gmail.com")

	// Set E-Mail receivers
	// m.SetHeader("To", user.Email)
	m.SetHeader("To", "sourabhsd87@gmail.com")

	// Set E-Mail subject
	m.SetHeader("Subject", "Payment Initated")

	// Parse the HTML template
	tmpl, err := template.ParseFiles("/home/trellis/Sourabh/GoTransact/Backend/apps/transaction/utils/email_template.html")
	if err != nil {
		log.Fatal("Error parsing email template: ", err)
	}

	// Create a buffer to hold the executed template
	var body bytes.Buffer

	// Execute the template with the data
	TemplateData := TemplateData{
		Username: user.FirstName,
		Amount:   request.Amount,
	}
	fmt.Println(TemplateData)
	if err := tmpl.Execute(&body, TemplateData); err != nil {
		log.Fatal("Error executing email template: ", err)
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
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("end of mail")
}
