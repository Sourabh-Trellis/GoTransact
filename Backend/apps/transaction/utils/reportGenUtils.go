package utils

import (
	transactionmodels "GoTransact/apps/transaction/models"
	"GoTransact/pkg/db"
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	gomail "gopkg.in/mail.v2"
)

func FetchTransactionsLast24Hours() []transactionmodels.TransactionRequest {
	var transactions []transactionmodels.TransactionRequest
	last24Hours := time.Now().Add(-24 * time.Hour)
	db.DB.Where("created_at >= ?", last24Hours).Find(&transactions)
	return transactions
}

func GenerateExcel(transactions []transactionmodels.TransactionRequest) (string, error) {
	f := excelize.NewFile()
	sheetName := "Transactions"
	index := f.NewSheet(sheetName)

	f.SetCellValue(sheetName, "A1", "ID")
	f.SetCellValue(sheetName, "B1", "InternalID")
	f.SetCellValue(sheetName, "C1", "UserID")
	f.SetCellValue(sheetName, "D1", "Status")
	f.SetCellValue(sheetName, "E1", "PaymentGatewayID")
	f.SetCellValue(sheetName, "F1", "Description")
	f.SetCellValue(sheetName, "G1", "Amount")
	f.SetCellValue(sheetName, "H1", "CreatedAt")
	f.SetCellValue(sheetName, "I1", "UpdatedAt")

	for i, tr := range transactions {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), tr.ID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), tr.Internal_id)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), tr.UserID)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), tr.Status)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), tr.Payment_Gateway_id)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), tr.Description)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), tr.Amount)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), tr.CreatedAt)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), tr.UpdatedAt)
	}

	f.SetActiveSheet(index)
	filePath := "transactions.xlsx"
	if err := f.SaveAs(filePath); err != nil {
		return "", err
	}

	return filePath, nil
}

func SendMailWithAttachment(email, filePath string) {

	m := gomail.NewMessage()
	m.SetHeader("From", "sourabhtrellis@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Daily Transactions Report")
	m.SetBody("text/plain", "Please find attached the daily transactions report.")
	m.Attach(filePath)

	d := gomail.NewDialer("smtp.gmail.com", 587, "sourabhtrellis@gmail.com", "nmvx vzro ehqo xwpd")

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		log.Printf("could not send email: %v", err)
	}
	fmt.Println("Email sent successfully")
}
