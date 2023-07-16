package util

import (
	"os"

	"github.com/hibiken/asynq"
)

func GetAsynqClient() *asynq.Client {

	redisHost, exists := os.LookupEnv("REDIS_HOST")

	if !exists {
		panic("REDIS_HOST is not defined!")
	}

	redisPort, exists := os.LookupEnv("REDIS_PORT")

	if !exists {
		panic("REDIS_PORT is not defined!")
	}

	inspector := asynq.NewInspector(asynq.RedisClientOpt{Addr: redisHost + ":" + redisPort})
	scheduledTasks, _ := inspector.ListScheduledTasks("runs")
	activeTasks, _ := inspector.ListActiveTasks("runs")

	print(scheduledTasks)
	print(activeTasks)

	return asynq.NewClient(asynq.RedisClientOpt{Addr: redisHost + ":" + redisPort})
}

func GetAsynqScheduler() *asynq.Scheduler {

	redisHost, exists := os.LookupEnv("REDIS_HOST")

	if !exists {
		panic("REDIS_HOST is not defined!")
	}

	redisPort, exists := os.LookupEnv("REDIS_PORT")

	if !exists {
		panic("REDIS_PORT is not defined!")
	}

	return asynq.NewScheduler(asynq.RedisClientOpt{Addr: redisHost + ":" + redisPort}, &asynq.SchedulerOpts{})
}
