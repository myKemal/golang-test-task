package redis

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/myKemal/golang-test-task/cmd/models"
)

type RedisClient struct {
	Client redis.Conn
	Error  error
}

func NewRedisClient(connectionUrl string) *RedisClient {

	conn, err := redis.Dial("tcp", connectionUrl)
	if err != nil {
		return &RedisClient{
			Error: err,
		}
	}
	s, err := redis.String(conn.Do("PING"))
	if err != nil {
		return &RedisClient{
			Error: err,
		}
	}

	fmt.Println("Redis - > ", s)

	return &RedisClient{
		Client: conn,
	}
}

func (red *RedisClient) GetMessageList(query models.GetMessagePayload) ([]models.NewMessagePayload, error) {

	keys, err := redis.Strings(red.Client.Do("KEYS", "*"))

	if err != nil {
		return nil, err
	}

	var res []models.NewMessagePayload

	for _, key := range keys {

		sender, _ := redis.String(red.Client.Do("HGET", key, "sender"))
		reciever, _ := redis.String(red.Client.Do("HGET", key, "reciever"))
		message, _ := redis.String(red.Client.Do("HGET", key, "message"))

		fmt.Print(sender, "-|", query.Sender, "|----", reciever, "-|", query.Receiver)

		if sender == query.Sender && reciever == query.Receiver {
			object := models.NewMessagePayload{
				Sender:   sender,
				Receiver: reciever,
				Message:  message,
			}

			res = append(res, object)
		}
	}

	defer red.Client.Close()
	return res, nil
}

func (redis *RedisClient) SetMessage(newmessage models.NewMessagePayload) (interface{}, error) {

	res, err := redis.Client.Do("HSET", uuid.NewString(), "sender", newmessage.Sender, "reciever", newmessage.Receiver, "message", newmessage.Message)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println(newmessage.Message, "--->", res)

	return res, nil

}

func (client *RedisClient) Close() {
	client.Client.Close()
}
