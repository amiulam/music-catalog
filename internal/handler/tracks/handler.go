package tracks

import (
	"context"

	"github.com/amiulam/music-catalog/internal/middleware"
	"github.com/amiulam/music-catalog/internal/models/spotify"
	"github.com/gin-gonic/gin"
)

//go:generate mockgen -source=handler.go -destination=handler_mock.go -package=tracks
type service interface {
	Search(ctx context.Context, query string, pageSize, pageIndex int) (*spotify.SearchResponse, error)
}

type Handler struct {
	*gin.Engine
	trackSvc service
}

func NewHandler(api *gin.Engine, trackSvc service) *Handler {
	return &Handler{
		Engine:   api,
		trackSvc: trackSvc,
	}
}

func (h *Handler) RegisterRoute() {
	route := h.Group("/tracks")
	route.Use(middleware.AuthMiddleware())
	route.GET("/search", h.Search)
}
