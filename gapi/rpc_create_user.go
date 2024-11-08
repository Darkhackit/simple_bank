package gapi

import (
	"context"
	db "github.com/Darkhackit/simplebank/db/sqlc"
	"github.com/Darkhackit/simplebank/pb"
	"github.com/Darkhackit/simplebank/util"
	"github.com/Darkhackit/simplebank/val"
	"github.com/Darkhackit/simplebank/worker"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	_ "net/http"
	"time"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	violations := ValidCreateUserRequest(req)
	if len(violations) > 0 {
		return nil, invalidArgumentError(violations)
	}
	name := pgtype.Text{
		String: req.GetName(),
		Valid:  true,
	}
	active := pgtype.Bool{
		Bool:  req.GetActive(),
		Valid: true,
	}
	password, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to hash password: %v", err)
	}
	arg := db.CreateUserParams{
		Name:     name,
		Username: req.GetUsername(),
		Password: password,
		Email:    req.GetEmail(),
		Active:   active,
	}

	user, err := server.q.CreateUser(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to created user: %v", err)
	}
	taskPayload := &worker.PayloadSendVerifyEmail{
		Username: user.Username,
	}
	ops := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(worker.QueueCritical),
	}
	err = server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, ops...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to send verify email: %v", err)
	}

	response := &pb.CreateUserResponse{
		User: &pb.User{
			Username: user.Username,
			Email:    user.Email,
		},
	}

	return response, nil
}

func ValidCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}
	if err := val.ValidateName(req.GetName()); err != nil {
		violations = append(violations, fieldViolation("name", err))
	}
	if err := val.ValidateName(req.GetName()); err != nil {
		violations = append(violations, fieldViolation("name", err))
	}
	return
}
