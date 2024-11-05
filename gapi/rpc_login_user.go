package gapi

import (
	"context"
	"github.com/Darkhackit/simplebank/pb"
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
	token, _, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
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
