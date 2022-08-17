package redis

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/myKemal/golang-test-task/cmd/models"
	"github.com/nitishm/go-rejson"
)

type RedisClient struct {
	Client *redis.Conn
	Error  error
}

func NewRedisClient(connectionUrl string) *RedisClient {

	conn, err := redis.Dial("tcp", connectionUrl)
	if err != nil {
		return &RedisClient{
			Error: err,
		}
	}
	return &RedisClient{
		Client: &conn,
	}
}

func (client *RedisClient) GetMessageList(query models.GetMessagePayload) (interface{}, error) {

	rh := rejson.NewReJSONHandler()
	rh.SetRedigoClient(*client.Client)
	res, err := rh.JSONGet("Message", fmt.Sprint(query.Sender, "_", query.Receiver))

	if err != nil {
		return nil, err
	}

	return res, err
}

func (client *RedisClient) SetMessage(newmessage models.NewMessagePayload) (interface{}, error) {

	rh := rejson.NewReJSONHandler()
	rh.SetRedigoClient(*client.Client)
	res, err := rh.JSONSet("Message", fmt.Sprint(newmessage.Sender, "_", newmessage.Receiver), newmessage)

	if err != nil {
		return nil, err
	}

	return res, nil

}
