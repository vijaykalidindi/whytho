package server

import (
	"github.com/gin-gonic/gin"
	"github.com/vinamra28/operator-reviewer/internal/config"
	"github.com/vinamra28/operator-reviewer/internal/handlers"
	"github.com/vinamra28/operator-reviewer/internal/services"
)

type Server struct {
	config *config.Config
	router *gin.Engine
}

func New(cfg *config.Config) *Server {
	router := gin.Default()
	
	gitlabService := services.NewGitLabService(cfg.GitLabToken, cfg.GitLabBaseURL)
	reviewService := services.NewReviewService(cfg.GeminiAPIKey)
	
	webhookHandler := handlers.NewWebhookHandler(gitlabService, reviewService, cfg.WebhookSecret)
	
	router.POST("/webhook", webhookHandler.HandleWebhook)
	router.GET("/health", handlers.HealthCheck)
	
	return &Server{
		config: cfg,
		router: router,
	}
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}