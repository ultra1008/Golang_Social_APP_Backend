package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/niklod/highload-social-network/config"
)

type HTTPServer struct {
	Server          *http.Server
	BaseRouterGroup *gin.RouterGroup
}

func NewHTTPServer(cfg *config.HTTPServerConfig) *HTTPServer {
	engine := gin.Default()
	group := engine.Group("/")

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      engine,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	return &HTTPServer{
		BaseRouterGroup: group,
		Server:          srv,
	}
}

func (h *HTTPServer) Start() {
	log.Printf("running HTTP server on port %s", h.Server.Addr)
	go func() {
		if err := h.Server.ListenAndServe(); err != nil {
			log.Printf("http server: %v", err)
		}
	}()
}

func (h *HTTPServer) Shutdown() {
	log.Println("Server shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.Server.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
}
