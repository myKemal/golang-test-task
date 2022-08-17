package rabbitmq

import (
	"encoding/json"
	"fmt"

	_const "github.com/myKemal/golang-test-task/cmd/constant"
	"github.com/myKemal/golang-test-task/cmd/models"
	_redis "github.com/myKemal/golang-test-task/cmd/services/redis"
	"github.com/streadway/amqp"
)

type RabbitClient struct {
	Client *amqp.Connection
	Error  error
}

func NewRabbitClient(connectionUrl string) *RabbitClient {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return &RabbitClient{
			Error: err,
		}
	}

	return &RabbitClient{
		Client: conn,
		Error:  nil,
	}
}

func (mq *RabbitClient) Listen() {

	redisClient := _redis.NewRedisClient(_const.RedisURL)

	if redisClient.Error != nil {
		fmt.Errorf(redisClient.Error.Error())
	}

	ch, err := mq.Client.Channel()
	if err != nil {
		defer mq.Client.Close()
		defer ch.Close()
		fmt.Errorf(err.Error())
	}

	msgs, err := ch.Consume(
		"NewMessage", // Bu sfer dinleyeceğim kuyruk ismini kendim yazdım
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		fmt.Errorf(err.Error())
	}

	forever := make(chan bool)

	go func(redis *_redis.RedisClient) {
		//Burada eğer varsa kuyruktaki mesajları çekiyoruz
		for d := range msgs {
			//d değişkeni ile kuyruktaki mesajın bilgilerine ulaşabiliriz.
			var newmessage models.NewMessagePayload
			err := json.Unmarshal(d.Body, newmessage)

			if err == nil {
				redis.SetMessage(newmessage)
			}

			//Kuyruktaki mesaj ekrana bastırdık.
		}
	}(redisClient)

	<-forever

}

func (mq *RabbitClient) SendMessage(publishData models.RabbitMessage) error {
	// convert incomming Message to amqp.publis
	message, e := json.Marshal(publishData.Message)
	if e != nil {
		defer mq.Client.Close()
		return e
	}

	ch, err := mq.Client.Channel()
	if err != nil {
		defer mq.Client.Close()
		defer ch.Close()
		return err
	}

	errPublish := ch.Publish(
		"",
		publishData.Quene,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(message),
		})

	if errPublish != nil {
		defer mq.Client.Close()
		defer ch.Close()
		return errPublish
	}

	defer mq.Client.Close()
	defer ch.Close()
	return nil
}
