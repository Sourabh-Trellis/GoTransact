package utils

import (
	log "GoTransact/settings"
	"crypto/tls"
	"fmt"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	gomail "gopkg.in/mail.v2"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func SendMail(to string) {

	log.InfoLogger.WithFields(logrus.Fields{}).Info("Attempted to send mail on registrtion to ", to)

	fmt.Println("start of mail")
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", "sourabhtrellis@gmail.com")

	// Set E-Mail receivers
	m.SetHeader("To", to)

	// Set E-Mail subject
	m.SetHeader("Subject", "Registration successfull")

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", "YOU HAVE REGISTERED SUCCESSFULLY ON GOTRANSACT")

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, "sourabhtrellis@gmail.com", "nmvx vzro ehqo xwpd")

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		log.ErrorLogger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Error while sending mail")
		panic(err)
	}
	log.InfoLogger.WithFields(logrus.Fields{}).Info("Registration mail sent to ", to)
}
