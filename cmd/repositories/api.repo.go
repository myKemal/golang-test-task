package repositories

import (
	"github.com/myKemal/golang-test-task/cmd/models"
	_rb "github.com/myKemal/golang-test-task/cmd/services/rabbitmq"
	_redis "github.com/myKemal/golang-test-task/cmd/services/redis"
)

func AddMessage(message models.NewMessagePayload, client *_rb.RabbitClient) error {

	return client.SendMessage(models.RabbitMessage{Quene: "NewMessage", Message: message})

}

func GetMessageList(message models.GetMessagePayload, client *_redis.RedisClient) (interface{}, error) {

	return client.GetMessageList(models.GetMessagePayload{Sender: message.Sender, Receiver: message.Receiver})
}
