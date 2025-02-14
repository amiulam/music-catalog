package tracks

import (
	"context"

	trackactivities "github.com/amiulam/music-catalog/internal/models/track_activities"
	"github.com/amiulam/music-catalog/internal/repository/spotify"
)

//go:generate mockgen -source=service.go -destination=service_mock.go -package=tracks
type spotifyOutbound interface {
	Search(ctx context.Context, query string, limit, offset int) (*spotify.SpotifySearchResponse, error)
}

type trackActivitesRepository interface {
	Create(ctx context.Context, model trackactivities.TrackActivity) error
	Update(ctx context.Context, model trackactivities.TrackActivity) error
	Get(ctx context.Context, userID uint, spotifyID string) (*trackactivities.TrackActivity, error)
	GetBulkSpotifyIDs(ctx context.Context, userID uint, spotifyIDs []string) (map[string]trackactivities.TrackActivity, error)
}

type service struct {
	spotifyOutbound    spotifyOutbound
	trackActivitesRepo trackActivitesRepository
}

func NewService(spotifyOutbound spotifyOutbound, trackActivitesRepo trackActivitesRepository) *service {
	return &service{
		spotifyOutbound:    spotifyOutbound,
		trackActivitesRepo: trackActivitesRepo,
	}
}
