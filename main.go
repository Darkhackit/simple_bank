package main

import (
	"context"
	"github.com/Darkhackit/simplebank/api"
	db "github.com/Darkhackit/simplebank/db/sqlc"
	"github.com/Darkhackit/simplebank/gapi"
	"github.com/Darkhackit/simplebank/pb"
	"github.com/Darkhackit/simplebank/util"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	var err error
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load configuration file", err)
	}
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, config.DbSource)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	queries := db.New(conn)
	runGrpcServer(config, queries)

}
func runGrpcServer(config util.Config, queries *db.Queries) {
	grpcServer := grpc.NewServer()
	server, err := gapi.NewServer(config, queries)
	if err != nil {
		log.Fatal(err)
	}
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Starting gRPC server on", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}
}

func runGinServer(config util.Config, queries *db.Queries) {
	server, err := api.NewServer(config, queries)
	if err != nil {
		log.Fatal(err)
	}
	err = server.Start(config.ServerAddr)
	if err != nil {
		log.Fatal(err)
	}
}
