package main

import (
	"context"
	"github.com/Darkhackit/simplebank/api"
	db "github.com/Darkhackit/simplebank/db/sqlc"
	"github.com/Darkhackit/simplebank/gapi"
	"github.com/Darkhackit/simplebank/pb"
	"github.com/Darkhackit/simplebank/util"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"net"
	"net/http"
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
	go runGatewayServer(config, queries)
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
func runGatewayServer(config util.Config, queries *db.Queries) {

	server, err := gapi.NewServer(config, queries)
	if err != nil {
		log.Fatal(err)
	}
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})
	grpcMux := runtime.NewServeMux(jsonOption)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)
	listener, err := net.Listen("tcp", config.ServerAddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Starting Http gateway server on", listener.Addr().String())
	err = http.Serve(listener, mux)
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
