package functions

import (
	"GoTransact/apps/accounts/models"
	"GoTransact/apps/accounts/utils"
	validator "GoTransact/apps/accounts/validators"
	db "GoTransact/pkg/db"
	log "GoTransact/settings"
	"net/http"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func Login(loginuser utils.LoginInput) (int, string, map[string]interface{}) {

	log.InfoLogger.WithFields(logrus.Fields{}).Info("Attempted to login with ", loginuser.Email)

	// custom validator for additional password validation
	if err := validator.GetValidator().Struct(loginuser); err != nil {

		return http.StatusBadRequest, "Password should contain atleast one upper case character,one lower case character,one number and one special character", map[string]interface{}{}
	}

	var user models.User
	if err := db.DB.Where("email = ?", loginuser.Email).First(&user).Error; err != nil {
		log.ErrorLogger.WithFields(logrus.Fields{}).Error("Failed to login with ", loginuser.Email)
		return http.StatusUnauthorized, "invalid username or password", map[string]interface{}{}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginuser.Password)); err != nil {
		log.ErrorLogger.WithFields(logrus.Fields{}).Error("Failed to login with ", loginuser.Email)
		return http.StatusUnauthorized, "invalid username or password", map[string]interface{}{}
	}

	token, err := utils.CreateToken(user)
	if err != nil {

		log.ErrorLogger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("error creating token")

		return http.StatusInternalServerError, "error creating token", map[string]interface{}{"data": err.Error()}
	}

	log.InfoLogger.WithFields(logrus.Fields{}).Info("User logged in with ", loginuser.Email, " and id ", user.Internal_id)

	return http.StatusOK, "Logged in successfull", map[string]interface{}{"token": token}
}
