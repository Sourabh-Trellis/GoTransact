package handlers

import (
	"GoTransact/apps/accounts/models"
	"GoTransact/apps/accounts/utils"
	validator "GoTransact/apps/accounts/validators"
	basemodels "GoTransact/apps/base"
	db "GoTransact/pkg/db"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8" validate:"password_complexity"`
}

func Login(c *gin.Context) {

	var loginInput LoginInput
	if err := c.ShouldBindJSON(&loginInput); err != nil {
		c.JSON(http.StatusBadRequest, basemodels.Response{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    map[string]interface{}{"data": err.Error()},
		})
		return
	}

	// custom validator for additional password validation
	if err := validator.GetValidator().Struct(loginInput); err != nil {
		c.JSON(http.StatusBadRequest, basemodels.Response{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    map[string]interface{}{"data": "Password should contain atleast one upper case character,one lower case character,one number and one special character"},
		})
		return
	}

	var user models.User
	if err := db.DB.Where("email = ?", loginInput.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, basemodels.Response{
			Status:  http.StatusUnauthorized,
			Message: "error",
			Data:    map[string]interface{}{"data": "invalid username or password"},
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginInput.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, basemodels.Response{
			Status:  http.StatusUnauthorized,
			Message: "error",
			Data:    map[string]interface{}{"data": "invalid username or password"},
		})
		return
	}

	token, err := utils.CreateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, basemodels.Response{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    map[string]interface{}{"data": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, basemodels.Response{
		Status:  http.StatusOK,
		Message: "success",
		Data:    map[string]interface{}{"token": token},
	})
}
