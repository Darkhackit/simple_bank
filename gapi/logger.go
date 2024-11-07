package gapi

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

func GrpcLogger(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	log.Println("Received a GRP request")
	result, err := handler(ctx, req)
	if err != nil {
		return nil, err
	}
	return result, err
}
