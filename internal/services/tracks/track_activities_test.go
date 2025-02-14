package tracks

import (
	"context"
	"fmt"
	"testing"

	trackactivities "github.com/amiulam/music-catalog/internal/models/track_activities"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func Test_service_UpsertUserActivities(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockTrackActivityRepo := NewMocktrackActivitesRepository(mockCtrl)
	isLikedTrue := true
	isLikedFalse := false

	type args struct {
		userID  uint
		request trackactivities.TrackActivityRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name:    "success: create",
			wantErr: false,
			args: args{
				userID: 1,
				request: trackactivities.TrackActivityRequest{
					SpotifyID: "spotifyID",
					IsLiked:   &isLikedTrue,
				},
			},
			mockFn: func(args args) {
				mockTrackActivityRepo.EXPECT().Get(gomock.Any(), args.userID, args.request.SpotifyID).Return(nil, gorm.ErrRecordNotFound)

				mockTrackActivityRepo.EXPECT().Create(gomock.Any(), trackactivities.TrackActivity{
					UserID:    args.userID,
					SpotifyID: args.request.SpotifyID,
					IsLiked:   args.request.IsLiked,
					CreatedBy: fmt.Sprintf("%d", args.userID),
					UpdatedBy: fmt.Sprintf("%d", args.userID),
				}).Return(nil)
			},
		},
		{
			name:    "success: update",
			wantErr: false,
			args: args{
				userID: 1,
				request: trackactivities.TrackActivityRequest{
					SpotifyID: "spotifyID",
					IsLiked:   &isLikedTrue,
				},
			},
			mockFn: func(args args) {
				mockTrackActivityRepo.EXPECT().Get(gomock.Any(), args.userID, args.request.SpotifyID).Return(&trackactivities.TrackActivity{
					IsLiked: &isLikedFalse,
				}, nil)

				mockTrackActivityRepo.EXPECT().Update(gomock.Any(), trackactivities.TrackActivity{
					IsLiked: args.request.IsLiked,
				}).Return(nil)
			},
		},
		{
			name:    "failed",
			wantErr: true,
			args: args{
				userID: 1,
				request: trackactivities.TrackActivityRequest{
					SpotifyID: "spotifyID",
					IsLiked:   &isLikedTrue,
				},
			},
			mockFn: func(args args) {
				mockTrackActivityRepo.EXPECT().Get(gomock.Any(), args.userID, args.request.SpotifyID).Return(nil, assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			s := &service{
				trackActivitesRepo: mockTrackActivityRepo,
			}
			if err := s.UpsertUserActivities(context.Background(), tt.args.userID, tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("service.UpsertUserActivities() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
