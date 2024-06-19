package functions

import (
	"GoTransact/apps/accounts/utils"
	validator "GoTransact/apps/accounts/validators"

	"net/http"

	// "GoTransact/apps/accounts/models"

	accountModels "GoTransact/apps/accounts/models"

	db "GoTransact/pkg/db"
	log "GoTransact/settings"

	"github.com/sirupsen/logrus"
)

func Signup(user utils.RegisterInput) (int, string, map[string]interface{}) {

	log.InfoLogger.WithFields(logrus.Fields{}).Info("Attempted to register with ", user.Email, " and company ", user.Companyname)

	if err := validator.GetValidator().Struct(user); err != nil {
		return http.StatusBadRequest, "Password should contain atleast one upper case character,one lower case character,one number and one special character", map[string]interface{}{}
	}

	//chaecking if user with email already exist
	var count int64

	// Check if user with the email already exists
	if err := db.DB.Model(&accountModels.User{}).Where("email = ?", user.Email).Count(&count).Error; err != nil {
		return http.StatusInternalServerError, "Database error", map[string]interface{}{}
	}
	if count > 0 {
		return http.StatusBadRequest, "email already exists", map[string]interface{}{}
	}

	// Check if company with the name already exists
	if err := db.DB.Model(&accountModels.Company{}).Where("name = ?", user.Companyname).Count(&count).Error; err != nil {
		return http.StatusInternalServerError, "Database error", map[string]interface{}{}
	}
	if count > 0 {
		return http.StatusBadRequest, "company already exists", map[string]interface{}{}
	}

	//hashing the password to store in database
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return http.StatusInternalServerError, "Error while hashing password", map[string]interface{}{}
	}

	//creating user and company model
	newuser := accountModels.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  hashedPassword,
		Company: accountModels.Company{
			Name: user.Companyname,
		},
	}

	//save the user
	if err := db.DB.Create(&newuser).Error; err != nil {
		//log
		log.ErrorLogger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Error creating user in database")

		return http.StatusInternalServerError, "error creating user", map[string]interface{}{}
	}
	//log
	log.InfoLogger.WithFields(logrus.Fields{}).Info("User created in database ", user.Email, " and company ", user.Companyname)

	go utils.SendMail(user.Email)

	return http.StatusOK, "User created successfully", map[string]interface{}{}
}
