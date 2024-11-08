package worker

import (
	db "github.com/Darkhackit/simplebank/db/sqlc"
	"github.com/hibiken/asynq"
	"golang.org/x/net/context"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSenderVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	q      db.Queries
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSenderVerifyEmail, processor.ProcessTaskSenderVerifyEmail)

	err := processor.server.Start(mux)
	if err != nil {
		return err
	}
	return nil
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Queries) TaskProcessor {
	server := asynq.NewServer(redisOpt, asynq.Config{
		Queues: map[string]int{
			QueueCritical: 10,
			QueueDefault:  5,
		},
	})
	return &RedisTaskProcessor{server: server, q: store}
}
