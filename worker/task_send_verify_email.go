package worker

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

const (
	TaskSenderVerifyEmail = "task:send_verify_email"
)

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(ctx context.Context, payload *PayloadSendVerifyEmail, opts ...asynq.Option) error {
	jsomPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("encode payload: %w", err)
	}
	task := asynq.NewTask(TaskSenderVerifyEmail, jsomPayload, opts...)

	enqueueContext, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("Failed to enqueue task: %w", err)
	}
	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("queue", enqueueContext.Queue).
		Int("max_retry", enqueueContext.MaxRetry).
		Msg("Enqueue task")

	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSenderVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("Failed to unmarshal payload: %w", asynq.SkipRetry)
	}
	user, err := processor.q.GetUser(ctx, payload.Username)
	if err != nil {
		return fmt.Errorf("Failed to get user from queue: %w", err)
	}
	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("email", user.Email).
		Msg("Enqueue task")

	return nil
}
