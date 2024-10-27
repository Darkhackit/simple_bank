package api

import (
	db "github.com/Darkhackit/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	q      *db.Queries
	router *gin.Engine
}

func NewServer(store *db.Queries) *Server {

	server := &Server{q: store}

	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	server.router = router
	return server
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
