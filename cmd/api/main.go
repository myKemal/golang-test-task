package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	_const "github.com/myKemal/golang-test-task/cmd/constant"
	"github.com/myKemal/golang-test-task/cmd/handlers"
	_rb "github.com/myKemal/golang-test-task/cmd/services/rabbitmq"
)

func main() {
	fmt.Println("Starting API")
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	go func() {
		rabbitClient := _rb.NewRabbitClient(_const.RabbitmqURL)

		if rabbitClient.Error != nil {
			fmt.Println(rabbitClient.Error.Error())
		} else {
			rabbitClient.Listen()
		}
	}()

	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, "worked")
	})

	r.POST("/message", handlers.SendMessage)

	r.GET("/message/list", handlers.GetMessage)

	r.Run()
}
