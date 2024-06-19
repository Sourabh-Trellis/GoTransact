package handlers

import (
	"GoTransact/apps/accounts/utils"
	basemodels "GoTransact/apps/base"
	log "GoTransact/settings"
	"net/http"
	"time"

	// rdb "github.com/go-redis/redis/v8"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

var (
	//ctx = context.Background()
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis server address
	})
)

func LogoutHandler(c *gin.Context) {

	log.InfoLogger.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"url":    c.Request.URL.String(),
	}).Info("Attempted to logout")

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, basemodels.Response{
			Status:  http.StatusUnauthorized,
			Message: "authorization header missing",
		})
		return
	}

	//tokenStr := authHeader[len("Bearer "):]
	_, err := utils.VerifyToken(authHeader)
	if err != nil {
		c.JSON(http.StatusUnauthorized, basemodels.Response{
			Status:  http.StatusUnauthorized,
			Message: "invalid token",
			Data: map[string]interface{}{
				"data": err.Error(),
			},
		})
		return
	}

	// Blacklist the token by storing it in Redis with an expiration time
	err = rdb.Set(authHeader, "Blacklisted", 24*time.Hour).Err() // adjust expiration time as needed
	if err != nil {
		c.JSON(http.StatusInternalServerError, basemodels.Response{
			Status:  http.StatusInternalServerError,
			Message: "failed to blacklist token",
			Data: map[string]interface{}{
				"data": err.Error(),
			},
		})
		return
	}

	log.InfoLogger.WithFields(logrus.Fields{
		// "method": c.Request.Method,
		// "url":    c.Request.URL.String(),
	}).Info("Logged out successfully")
	c.JSON(http.StatusOK, basemodels.Response{
		Status:  http.StatusOK,
		Message: "logged out successfully",
	})
}
