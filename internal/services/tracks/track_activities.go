package tracks

import (
	"context"
	"fmt"

	trackactivities "github.com/amiulam/music-catalog/internal/models/track_activities"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func (s *service) UpsertUserActivities(ctx context.Context, userID uint, request trackactivities.TrackActivityRequest) error {
	activity, err := s.trackActivitesRepo.Get(ctx, userID, request.SpotifyID)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error().Err(err).Msg("error get record from database")
		return err
	}

	if err == gorm.ErrRecordNotFound || activity.IsLiked == nil {
		err = s.trackActivitesRepo.Create(ctx, trackactivities.TrackActivity{
			UserID:    userID,
			SpotifyID: request.SpotifyID,
			IsLiked:   request.IsLiked,
			CreatedBy: fmt.Sprintf("%d", userID),
			UpdatedBy: fmt.Sprintf("%d", userID),
		})

		if err != nil {
			log.Error().Err(err).Msg("error create record to database")
			return err
		}
		return nil
	}

	activity.IsLiked = request.IsLiked
	err = s.trackActivitesRepo.Update(ctx, *activity)

	if err != nil {
		log.Error().Err(err).Msg("error while update the record")
		return err
	}

	return nil
}
