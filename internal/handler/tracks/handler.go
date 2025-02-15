package tracks

import (
	"context"

	"github.com/amiulam/music-catalog/internal/middleware"
	"github.com/amiulam/music-catalog/internal/models/spotify"
	trackactivities "github.com/amiulam/music-catalog/internal/models/track_activities"
	"github.com/gin-gonic/gin"
)

//go:generate mockgen -source=handler.go -destination=handler_mock.go -package=tracks
type service interface {
	Search(ctx context.Context, query string, pageSize, pageIndex int, userID uint) (*spotify.SearchResponse, error)
	UpsertUserActivities(ctx context.Context, userID uint, request trackactivities.TrackActivityRequest) error
	GetRecommendation(ctx context.Context, limit int, trackID string, userID uint) (*spotify.RecommendationResponse, error)
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
	route.POST("/track-activity", h.UpsertTrackActivities)
	route.GET("/recommendations", h.GetRecommendation)
}
