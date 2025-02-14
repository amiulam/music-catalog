package tracks

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	trackactivities "github.com/amiulam/music-catalog/internal/models/track_activities"
	"github.com/amiulam/music-catalog/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestHandler_UpsertTrackActivities(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockService := NewMockservice(mockCtrl)
	isLikedTrue := true

	tests := []struct {
		name               string
		mockFn             func()
		expectedStatusCode int
	}{
		{
			name: "success",
			mockFn: func() {
				mockService.EXPECT().UpsertUserActivities(gomock.Any(), uint(1), trackactivities.TrackActivityRequest{
					SpotifyID: "spotifyID",
					IsLiked:   &isLikedTrue,
				}).Return(nil)
			},
			expectedStatusCode: http.StatusCreated,
		},
		{
			name: "failed",
			mockFn: func() {
				mockService.EXPECT().UpsertUserActivities(gomock.Any(), uint(1), trackactivities.TrackActivityRequest{
					SpotifyID: "spotifyID",
					IsLiked:   &isLikedTrue,
				}).Return(assert.AnError)
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			api := gin.New()

			h := &Handler{
				Engine:   api,
				trackSvc: mockService,
			}

			payload := trackactivities.TrackActivityRequest{
				SpotifyID: "spotifyID",
				IsLiked:   &isLikedTrue,
			}

			payloadBytes, err := json.Marshal(payload)
			assert.NoError(t, err)

			h.RegisterRoute()
			w := httptest.NewRecorder()
			endpoint := `/tracks/track-activity`

			req, err := http.NewRequest(http.MethodPost, endpoint, io.NopCloser(bytes.NewBuffer(payloadBytes)))
			assert.NoError(t, err)

			token, err := jwt.CreateToken(1, "username", "")
			assert.NoError(t, err)

			req.Header.Set("Authorization", token)

			h.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
		})
	}
}
