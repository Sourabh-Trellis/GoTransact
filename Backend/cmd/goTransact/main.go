package main

import (
	account "GoTransact/apps/accounts/handlers"
	accountModels "GoTransact/apps/accounts/models"
	accountValidator "GoTransact/apps/accounts/validators"
	basemodels "GoTransact/apps/base"
	transaction "GoTransact/apps/transaction/handlers"
	transactionModels "GoTransact/apps/transaction/models"
	transactionValidator "GoTransact/apps/transaction/validators"
	"GoTransact/config"
	"GoTransact/middlewares"
	db "GoTransact/pkg/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	db.InitDB()
	accountValidator.Init()
	transactionValidator.InitValidation()

	db.DB.AutoMigrate(&accountModels.User{}, &accountModels.Company{}, &transactionModels.Payment_Gateway{}, &transactionModels.TransactionRequest{}, &transactionModels.TransactionHistory{})

	r := gin.Default()

	r.POST("/api/register", account.Signup)
	r.POST("/api/login", account.Login)
	r.POST("/api/create", func(c *gin.Context) {
		gatway := transactionModels.Payment_Gateway{
			Slug:  "card",
			Label: "Card",
		}
		if err := db.DB.Create(&gatway).Error; err != nil {
			c.JSON(http.StatusInternalServerError, basemodels.Response{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    map[string]interface{}{"data": err.Error()},
			})
			return
		}
	})
	protected := r.Group("/protected")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.POST("/api/post-payment", transaction.PaymentRequest)
	}

	r.Run(":8080")
}
