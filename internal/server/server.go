package httpServer

import (
	"fmt"
	"log"
	"net/http"

	userHttp "github.com/bakare-dev/simple-bank-api/internal/user/http"
	"github.com/bakare-dev/simple-bank-api/internal/user/model"
	"github.com/bakare-dev/simple-bank-api/pkg/config"
	"github.com/bakare-dev/simple-bank-api/pkg/middleware"
	"github.com/bakare-dev/simple-bank-api/pkg/response"
	"github.com/bakare-dev/simple-bank-api/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Server struct {
	DB          *gorm.DB
	Redis       *redis.Client
	router      *gin.Engine
	rateLimiter *middleware.ClientRateLimiter
}

func NewServer() *Server {
	if config.Settings.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	db := util.ConnectDB()
	if db == nil {
		log.Fatalf("Database connection is nil")
	}

	if err := db.AutoMigrate(&model.User{}, &model.Profile{}); err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	redis := util.ConnectRedis()
	if redis == nil {
		log.Fatalf("Redis connection is nil")
	}

	rateLimiter := middleware.NewClientRateLimiter(5, 10)

	server := &Server{
		DB:          db,
		Redis:       redis,
		router:      gin.Default(),
		rateLimiter: rateLimiter,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	_ = s.router.SetTrustedProxies(nil)

	s.router.Use(
		middleware.CORSMiddleware(),
		middleware.SecurityHeadersMiddleware(),
		middleware.LoggerMiddleware(),
		s.rateLimiter.GinMiddleware(),
	)

	authInterceptor := middleware.NewAuthInterceptor()
	s.router.Use(middleware.RouteProtectionMiddleware(authInterceptor))

	s.router.GET("/", func(ctx *gin.Context) {
		response.JSON(ctx, http.StatusOK, nil, "server running fine")
	})

	if err := s.MapRoutes(); err != nil {
		log.Fatalf("MapRoutes Error: %v", err)
	}
}

func (s *Server) Run() error {
	host := config.Settings.Server.Host
	port := config.Settings.Server.Port

	address := fmt.Sprintf("%s:%d", host, port)
	log.Printf("Starting HTTP server at %s", address)

	return s.router.Run(address)
}

func (s *Server) MapRoutes() error {
	v1 := s.router.Group("/api/v1")
	userHttp.RegisterUserRoutes(v1, s.DB)
	return nil
}
