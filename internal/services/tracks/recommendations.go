package tracks

import (
	"context"

	"github.com/amiulam/music-catalog/internal/models/spotify"
	trackactivities "github.com/amiulam/music-catalog/internal/models/track_activities"
	spotifyRepo "github.com/amiulam/music-catalog/internal/repository/spotify"
	"github.com/rs/zerolog/log"
)

func (s *service) GetRecommendation(ctx context.Context, limit int, trackID string, userID uint) (*spotify.RecommendationResponse, error) {
	trackDetails, err := s.spotifyOutbound.GetRecommendation(ctx, limit, trackID)
	if err != nil {
		log.Error().Err(err).Msg("error get recommendations from spotify outbound")
		return nil, err
	}

	trackIDs := make([]string, len(trackDetails.Tracks))

	for idx, trackItem := range trackDetails.Tracks {
		trackIDs[idx] = trackItem.ID
	}

	trackActivites, err := s.trackActivitesRepo.GetBulkSpotifyIDs(ctx, userID, trackIDs)

	if err != nil {
		log.Error().Err(err).Msg("error get track activites from database")
		return nil, err
	}

	return modelToRecommendationResponse(trackDetails, trackActivites), nil
}

func modelToRecommendationResponse(data *spotifyRepo.SpotifyRecommendationResponse, mapTrackActivities map[string]trackactivities.TrackActivity) *spotify.RecommendationResponse {
	if data == nil {
		return nil
	}

	items := make([]spotify.SpotifyTrackObjects, 0)

	for _, item := range data.Tracks {
		artistsName := make([]string, len(item.Artists))

		for idx, artist := range item.Artists {
			artistsName[idx] = artist.Name
		}

		imageUrls := make([]string, len(item.Album.Images))

		for idx, image := range item.Album.Images {
			imageUrls[idx] = image.URL
		}

		items = append(items, spotify.SpotifyTrackObjects{
			// Album fields
			AlbumType:         item.Album.AlbumType,
			AlbumnTotalTracks: item.Album.TotalTracks,
			AlbumImagesURL:    imageUrls,
			AlbumName:         item.Album.Name,

			// Artists field
			ArtistsName: artistsName,

			// Track fields
			Explicit: item.Explicit,
			ID:       item.ID,
			Name:     item.Name,
			IsLiked:  mapTrackActivities[item.ID].IsLiked,
		})
	}

	return &spotify.RecommendationResponse{
		Items: items,
	}
}
