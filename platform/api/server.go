package api

import (
	db "github.com/UitsHabib/ecommerce-microservice/platform/db/sqlc"
	"github.com/UitsHabib/ecommerce-microservice/platform/util"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config util.Config
	store  db.Store
	router *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	server := &Server{
		config: config,
		store:  store,
	}

	// if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	// 	v.RegisterValidation("")
	// }
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.GET("/features", server.getFeatures)
	router.GET("/features/:id", server.getFeature)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
