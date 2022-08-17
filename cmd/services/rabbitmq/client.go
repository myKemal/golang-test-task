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

	conn, err := amqp.Dial(connectionUrl)
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

	fmt.Println("Listening MQ")

	ch, err := mq.Client.Channel()
	if err != nil {
		defer mq.Client.Close()
		defer ch.Close()
		fmt.Println(err.Error())
	}

	q, errq := ch.QueueDeclare(
		"NewMessage",
		false,
		false,
		false,
		false,
		nil,
	)

	if errq != nil {

		fmt.Println(err.Error())
	}
	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err.Error())
	}

	forever := make(chan bool)

	go func() {
		redisClient := _redis.NewRedisClient(_const.RedisURL)
		if redisClient.Error != nil {
			fmt.Println(redisClient.Error.Error())
			return
		}
		for d := range msgs {

			var newmessage models.NewMessagePayload
			err := json.Unmarshal(d.Body, &newmessage)
			if err == nil {
				fmt.Println(newmessage.Message, newmessage.Sender, newmessage.Receiver)

				redisClient.SetMessage(newmessage)
			}

		}
		defer redisClient.Close()

	}()

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
		fmt.Println(err)

		defer mq.Client.Close()
		defer ch.Close()
		return err
	}

	q, errQ := ch.QueueDeclare(
		publishData.Quene,
		false,
		false,
		false,
		false,
		nil,
	)

	if errQ != nil {
		fmt.Println(errQ)

		defer mq.Client.Close()
		defer ch.Close()
		return errQ
	}

	errPublish := ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(message),
		})

	if errPublish != nil {
		fmt.Println(errPublish)
		defer mq.Client.Close()
		defer ch.Close()
		return errPublish
	}

	defer mq.Client.Close()
	defer ch.Close()
	return nil
}
