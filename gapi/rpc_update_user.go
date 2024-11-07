package gapi

import (
	"context"
	db "github.com/Darkhackit/simplebank/db/sqlc"
	"github.com/Darkhackit/simplebank/pb"
	"github.com/Darkhackit/simplebank/util"
	"github.com/Darkhackit/simplebank/val"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	_ "net/http"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, err
	}
	if authPayload.Username != req.Username {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}
	violations := ValidUpdateUserRequest(req)
	if len(violations) > 0 {
		return nil, invalidArgumentError(violations)
	}
	name := pgtype.Text{
		String: req.GetName(),
		Valid:  req.GetName() != "",
	}
	email := pgtype.Text{
		String: req.GetEmail(),
		Valid:  req.GetEmail() != "",
	}
	password, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to hash password: %v", err)
	}
	p := pgtype.Text{
		String: password,
		Valid:  true,
	}
	arg := db.UpdateUserParams{
		Name:     name,
		Username: req.GetUsername(),
		Password: p,
		Email:    email,
	}

	user, err := server.q.UpdateUser(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to Update user: %v", err)
	}

	response := &pb.UpdateUserResponse{
		User: &pb.User{
			Username: user.Username,
			Email:    user.Email,
		},
	}

	return response, nil
}

func ValidUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
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
