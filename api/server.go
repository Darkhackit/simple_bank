package api

import (
	db "github.com/Darkhackit/simplebank/db/sqlc"
	"github.com/Darkhackit/simplebank/token"
	"github.com/Darkhackit/simplebank/util"
	"github.com/gin-gonic/gin"
)

type Server struct {
	q          *db.Queries
	s          db.SQLStore
	tokenMaker token.Maker
	router     *gin.Engine
	config     util.Config
}

func NewServer(config util.Config, store *db.Queries) (*Server, error) {
	tokenMaker, err := token.NewPasetoToken(config.TokenSymmetryKey)
	if err != nil {
		return nil, err
	}
	server := &Server{q: store, tokenMaker: tokenMaker, config: config}

	router := gin.Default()
	authRoutes := router.Group("/").Use(authMiddleware(tokenMaker))

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)

	authRoutes.POST("/transfers", server.createTransfer)
	router.POST("/users", server.createUser)
	router.POST("/login", server.loginUser)

	server.router = router
	return server, nil
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
