package gapi

import (
	db "github.com/Darkhackit/simplebank/db/sqlc"
	"github.com/Darkhackit/simplebank/pb"
	"github.com/Darkhackit/simplebank/token"
	"github.com/Darkhackit/simplebank/util"
	"github.com/Darkhackit/simplebank/worker"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	q               *db.Queries
	s               db.SQLStore
	tokenMaker      token.Maker
	config          util.Config
	taskDistributor worker.TaskDistributor
}

func NewServer(config util.Config, store *db.Queries, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoToken(config.TokenSymmetryKey)
	if err != nil {
		return nil, err
	}
	server := &Server{q: store, tokenMaker: tokenMaker, config: config, taskDistributor: taskDistributor}

	return server, nil
}
