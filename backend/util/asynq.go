package util

import (
	"os"

	"github.com/hibiken/asynq"
)

func GetAsynqClient() *asynq.Client {

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	return asynq.NewClient(asynq.RedisClientOpt{Addr: redisHost + ":" + redisPort})
}

func GetAsynqScheduler() *asynq.Scheduler {

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	return asynq.NewScheduler(asynq.RedisClientOpt{Addr: redisHost + ":" + redisPort}, &asynq.SchedulerOpts{})
}
