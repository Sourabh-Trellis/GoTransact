package handlers

import (
	"net/http"

	"GoTransact/apps/accounts/models"
	accountModels "GoTransact/apps/accounts/models"
	"GoTransact/apps/accounts/utils"
	validator "GoTransact/apps/accounts/validators"
	basemodels "GoTransact/apps/base"
	db "GoTransact/pkg/db"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RegisterInput struct {
	FirstName   string `json:"firstName" binding:"required"`
	LastName    string `json:"lastName" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Companyname string `json:"companyName" binding:"required"`
	Password    string `json:"password" binding:"required,min=8" validate:"password_complexity"`
}

func Signup(c *gin.Context) {

	//
	var registerInput RegisterInput
	if err := c.ShouldBindJSON(&registerInput); err != nil {
		c.JSON(http.StatusBadRequest, basemodels.Response{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    map[string]interface{}{"data": err.Error()},
		})
		return
	}

	// custom validator for additional password validation
	if err := validator.GetValidator().Struct(registerInput); err != nil {
		c.JSON(http.StatusBadRequest, basemodels.Response{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    map[string]interface{}{"data": "Password should contain atleast one upper case character,one lower case character,one number and one special character"},
		})
		return
	}

	//chaecking if user whit email already exist
	var existingUser accountModels.User
	if err := db.DB.Where("email = ?", registerInput.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, basemodels.Response{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    map[string]interface{}{"data": "email already exists"},
		})
		return
	} else if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	var existingCompany accountModels.Company
	if err := db.DB.Where("name = ?", registerInput.Companyname).First(&existingCompany).Error; err == nil {
		c.JSON(http.StatusBadRequest, basemodels.Response{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    map[string]interface{}{"data": "company profile already exists"},
		})
		return
	} else if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	//hashing the password to store in database
	hashedPassword, err := utils.HashPassword(registerInput.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, basemodels.Response{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    map[string]interface{}{"data": err.Error()},
		})
		return
	}

	//creating user and company model
	user := models.User{
		FirstName: registerInput.FirstName,
		LastName:  registerInput.LastName,
		Email:     registerInput.Email,
		Password:  hashedPassword,
		Company: accountModels.Company{
			Name: registerInput.Companyname,
		},
	}

	//save the user
	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, basemodels.Response{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    map[string]interface{}{"data": err.Error()},
		})
		return
	}

	go utils.SendMail(user.Email)

	c.JSON(http.StatusOK, basemodels.Response{
		Status:  http.StatusOK,
		Message: "success",
		Data:    map[string]interface{}{},
	})
}
