package handlers

import (
	"net/http"

	// "GoTransact/apps/accounts/models"
	"GoTransact/apps/accounts/functions"
	"GoTransact/apps/accounts/utils"
	basemodels "GoTransact/apps/base"
	log "GoTransact/settings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Signup_handler handles user registration
// @Summary Register a new user
// @Description Register a new user with email, password, etc.
// @Tags Accounts
// @Accept json
// @Produce json
// @Param 			registerInput body    	utils.RegisterInput true "User Registration Input"
// @Success 200 {object} basemodels.Response "Successfully registered"
// @Failure 400 {object} basemodels.Response "Invalid input"
// @Failure 500 {object} basemodels.Response "Internal server error"
// @Router /register [post]
func Signup_handler(c *gin.Context) {

	log.InfoLogger.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"url":    c.Request.URL.String(),
	}).Info("Register Request received")
	//
	var registerInput utils.RegisterInput
	if err := c.ShouldBindJSON(&registerInput); err != nil {
		c.JSON(http.StatusBadRequest, basemodels.Response{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    map[string]interface{}{"data": err.Error()},
		})
		return
	}

	statusCode, message, data := functions.Signup(registerInput)

	c.JSON(statusCode, basemodels.Response{
		Status:  statusCode,
		Message: message,
		Data:    data,
	})
}
