package main

//go:generate swagger generate spec -o docs/swagger.json
import (
	accountModels "GoTransact/apps/accounts/models"
	accountValidator "GoTransact/apps/accounts/validators"
	transactionModels "GoTransact/apps/transaction/models"
	"GoTransact/apps/transaction/utils"
	transactionValidator "GoTransact/apps/transaction/validators"
	_ "GoTransact/cmd/goTransact/docs"
	"GoTransact/config"
	db "GoTransact/pkg/db"
	"GoTransact/router"
	"log"

	// log "GoTransact/pkg/log"
	logger "GoTransact/settings"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/robfig/cron"
)

// @title GoTransact
// @version 1.0
// @description This is a sample server for a project.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api

// @schemes http
func main() {
	config.LoadEnv()
	db.InitDB("prod")
	accountValidator.Init()
	transactionValidator.InitValidation()
	logger.Init()

	db.DB.AutoMigrate(&accountModels.User{}, &accountModels.Company{}, &transactionModels.Payment_Gateway{}, &transactionModels.TransactionRequest{}, &transactionModels.TransactionHistory{})

	c := cron.New()
	c.AddFunc("@every 24h", func() {
		transactions := utils.FetchTransactionsLast24Hours()
		filePath, err := utils.GenerateExcel(transactions)
		if err != nil {
			log.Fatalf("failed to generate excel: %v", err)
		}
		utils.SendMailWithAttachment("sourabhsd87@gmail.com", filePath)
	})
	c.Start()

	r := router.Router()
	// r := gin.Default()
	// // docs.SwaggerInfo.BasePath = "/api"

	// r.POST("/api/register", account.Signup_handler)
	// r.POST("/api/login", account.Login_handler)
	// r.GET("/api/confirm-payment", transaction.ConfirmPayment)

	// protected := r.Group("/api/protected")
	// protected.Use(middlewares.AuthMiddleware())
	// {
	// 	protected.POST("/post-payment", transaction.PaymentRequest)
	// 	protected.POST("/logout", account.LogoutHandler)
	// }
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8080")

}
