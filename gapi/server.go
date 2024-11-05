package gapi

import (
	db "github.com/Darkhackit/simplebank/db/sqlc"
	"github.com/Darkhackit/simplebank/pb"
	"github.com/Darkhackit/simplebank/token"
	"github.com/Darkhackit/simplebank/util"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	q          *db.Queries
	s          db.SQLStore
	tokenMaker token.Maker
	config     util.Config
}

func NewServer(config util.Config, store *db.Queries) (*Server, error) {
	tokenMaker, err := token.NewPasetoToken(config.TokenSymmetryKey)
	if err != nil {
		return nil, err
	}
	server := &Server{q: store, tokenMaker: tokenMaker, config: config}

	return server, nil
}
