package gapi

import (
	"context"
	db "github.com/Darkhackit/simplebank/db/sqlc"
	"github.com/Darkhackit/simplebank/pb"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

	user, err := server.q.GetUser(ctx, req.GetUsername())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "method LoginUser not implemented")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.GetPassword()))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Username or Password is incorrect")
	}
	token, payload, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		return nil, err
	}
	uuid := pgtype.UUID{
		Bytes: payload.ID,
		Valid: true,
	}
	mtdt := server.extractMetadata(ctx)
	_, err = server.q.CreateSession(ctx, db.CreateSessionParams{
		ID:           uuid,
		RefreshToken: token,
		Username:     user.Username,
		UserAgent:    mtdt.UserAgent,
		ClientIp:     mtdt.ClientIp,
	})
	if err != nil {
		return nil, err
	}
	resp := &pb.LoginUserResponse{
		User: &pb.User{
			Username: user.Username,
			Email:    user.Email,
		},
		AccessToken: token,
	}
	return resp, nil
}
