package db

import (
	"GoTransact/config"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {

	dbURI := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=%s", config.DbHost, config.DbUser, config.DbPassword, config.DbName, config.DbPort, config.DbTimezone)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Microsecond, // Slow SQL threshold
			LogLevel:                  logger.Info,      // Log level
			IgnoreRecordNotFoundError: true,             // Ignore ErrRecord Not Found error for logger
			Colorful:                  true,             // Disable color
		},
	)

	db, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{Logger: newLogger})
	if err != nil {
		panic(err)
	}

	DB = db
}
