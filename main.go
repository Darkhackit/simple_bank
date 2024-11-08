package main

import (
	"context"
	"github.com/Darkhackit/simplebank/api"
	db "github.com/Darkhackit/simplebank/db/sqlc"
	"github.com/Darkhackit/simplebank/gapi"
	"github.com/Darkhackit/simplebank/pb"
	"github.com/Darkhackit/simplebank/util"
	"github.com/Darkhackit/simplebank/worker"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"net"
	"net/http"
	"os"
)

func main() {

	var err error
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Msg("Cannot load configuration file")
	}
	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, config.DbSource)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	defer conn.Close()

	queries := db.New(conn)
	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
	go runTaskProcessor(redisOpt, queries)
	go runGatewayServer(config, queries, taskDistributor)
	runGrpcServer(config, queries, taskDistributor)

}
func runGrpcServer(config util.Config, queries *db.Queries, taskDistributor worker.TaskDistributor) {

	server, err := gapi.NewServer(config, queries, taskDistributor)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	log.Print("Starting gRPC server on", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, store *db.Queries) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, *store)
	log.Info().Msg("Starting task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}
func runGatewayServer(config util.Config, queries *db.Queries, taskDistributor worker.TaskDistributor) {

	server, err := gapi.NewServer(config, queries, taskDistributor)
	if err != nil {
		log.Fatal().Msg(err.Error())
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
		log.Fatal().Msg(err.Error())
	}
	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)
	listener, err := net.Listen("tcp", config.ServerAddr)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	log.Print("Starting Http gateway server on", listener.Addr().String())
	handler := gapi.HttpLogger(mux)
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}
func runGinServer(config util.Config, queries *db.Queries) {
	server, err := api.NewServer(config, queries)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	err = server.Start(config.ServerAddr)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}
