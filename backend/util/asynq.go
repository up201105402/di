package util

import (
	"os"

	"github.com/hibiken/asynq"
)

func ConnectToAsynq() *asynq.Client {

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	return asynq.NewClient(asynq.RedisClientOpt{Addr: redisHost + ":" + redisPort})
}
