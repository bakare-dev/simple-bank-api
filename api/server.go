package api

import (
	"net/http"

	db "github.com/bakare-dev/simple-bank-api/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}

	router := gin.Default()

	router.POST("/api/v1/user", server.createUser)
	router.GET("/api/v1/user/all", server.getUsers)
	router.GET("/api/v1/user", server.getUser)
	router.PUT("/api/v1/user", server.updateUserpassword)
	router.DELETE("/api/v1/user", server.deleteUsers)

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "server running"})
	})

	server.router = router

	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
