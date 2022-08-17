package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_const "github.com/myKemal/golang-test-task/cmd/constant"
	"github.com/myKemal/golang-test-task/cmd/models"
	"github.com/myKemal/golang-test-task/cmd/repositories"
	"github.com/myKemal/golang-test-task/cmd/services"
	_rb "github.com/myKemal/golang-test-task/cmd/services/rabbitmq"
	_redis "github.com/myKemal/golang-test-task/cmd/services/redis"
)

func SendMessage(c *gin.Context) {
	var payload models.NewMessagePayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rbClient := _rb.NewRabbitClient(_const.RabbitmqURL)

	if rbClient.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": rbClient.Error.Error()})
	}

	err := repositories.AddMessage(payload, rbClient)
	c.JSON(services.CreateHttpStatusCode(err), true)

}

func GetMessage(c *gin.Context) {
	var payload models.GetMessagePayload
	if err := c.ShouldBindQuery(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	redisClient := _redis.NewRedisClient(_const.RedisURL)

	if redisClient.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": redisClient.Error.Error()})
	}

	response, err := repositories.GetMessageList(payload, redisClient)
	c.JSON(services.CreateHttpStatusCode(err), response)
}
