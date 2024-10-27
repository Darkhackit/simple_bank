package main

import (
	"context"
	"github.com/Darkhackit/simplebank/api"
	db "github.com/Darkhackit/simplebank/db/sqlc"
	"github.com/Darkhackit/simplebank/util"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
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

	server := api.NewServer(queries)

	err = server.Start(config.ServerAddr)
	if err != nil {
		log.Fatal(err)
	}
}
