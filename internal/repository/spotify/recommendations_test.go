package spotify

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/amiulam/music-catalog/internal/configs"
	"github.com/amiulam/music-catalog/pkg/httpclient"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_outbound_GetRecommendation(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockHTTPClient := httpclient.NewMockHTTPClient(ctrlMock)

	type args struct {
		limit   int
		trackID string
	}
	tests := []struct {
		name    string
		args    args
		want    *SpotifyRecommendationResponse
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "success",
			args: args{
				limit:   10,
				trackID: "trackID",
			},
			want: &SpotifyRecommendationResponse{
				Tracks: []SpotifyTrackObjects{
					{
						Album: SpotifyAlbumObject{
							AlbumType:   "album",
							TotalTracks: 22,
							Images: []SpotifyAlbumImage{
								{
									URL: "https://i.scdn.co/image/ab67616d0000b273e8b066f70c206551210d902b",
								}, {
									URL: "https://i.scdn.co/image/ab67616d00001e02e8b066f70c206551210d902b",
								}, {
									URL: "https://i.scdn.co/image/ab67616d00004851e8b066f70c206551210d902b",
								},
							},
							Name: "Bohemian Rhapsody (The Original Soundtrack)",
						},
						Artists: []SpotifyArtistObject{
							{
								Href: "https://api.spotify.com/v1/artists/1dfeR4HaWDbWqFHLkxsg1d",
								Name: "Queen",
							},
						},
						Explicit: false,
						Href:     "https://api.spotify.com/v1/tracks/3z8h0TU7ReDPLIbEnYhWZb",
						ID:       "3z8h0TU7ReDPLIbEnYhWZb",
						Name:     "Bohemian Rhapsody",
					},
					{
						Album: SpotifyAlbumObject{
							AlbumType:   "compilation",
							TotalTracks: 17,
							Images: []SpotifyAlbumImage{
								{
									URL: "https://i.scdn.co/image/ab67616d0000b273bb19d0c22d5709c9d73c8263",
								}, {
									URL: "https://i.scdn.co/image/ab67616d00001e02bb19d0c22d5709c9d73c8263",
								}, {
									URL: "https://i.scdn.co/image/ab67616d00004851bb19d0c22d5709c9d73c8263",
								},
							},
							Name: "Greatest Hits (Remastered)",
						},
						Artists: []SpotifyArtistObject{
							{
								Href: "https://api.spotify.com/v1/artists/1dfeR4HaWDbWqFHLkxsg1d",
								Name: "Queen",
							},
						},
						Explicit: false,
						Href:     "https://api.spotify.com/v1/tracks/2OBofMJx94NryV2SK8p8Zf",
						ID:       "2OBofMJx94NryV2SK8p8Zf",
						Name:     "Bohemian Rhapsody - Remastered 2011",
					},
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				params := url.Values{}
				params.Set("limit", strconv.Itoa(args.limit))
				params.Set("market", "ID")
				params.Set("seed_tracks", args.trackID)

				basePath := `https://api.spotify.com/v1/recommendations`
				endpoint := fmt.Sprintf("%s?%s", basePath, params.Encode())
				req, err := http.NewRequest(http.MethodGet, endpoint, nil)
				assert.NoError(t, err)

				req.Header.Set("Authorization", "Bearer accessToken")
				mockHTTPClient.EXPECT().Do(req).Return(&http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBufferString(recommendationResponse)),
				}, nil)
			},
		},
		{
			name: "failed",
			args: args{
				limit:   10,
				trackID: "trackID",
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				params := url.Values{}
				params.Set("limit", strconv.Itoa(args.limit))
				params.Set("market", "ID")
				params.Set("seed_tracks", args.trackID)

				basePath := `https://api.spotify.com/v1/recommendations`
				endpoint := fmt.Sprintf("%s?%s", basePath, params.Encode())
				req, err := http.NewRequest(http.MethodGet, endpoint, nil)
				assert.NoError(t, err)

				req.Header.Set("Authorization", "Bearer accessToken")
				mockHTTPClient.EXPECT().Do(req).Return(nil, assert.AnError)
			},
		},
		{
			name: "failed: 500",
			args: args{
				limit:   10,
				trackID: "trackID",
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				params := url.Values{}
				params.Set("limit", strconv.Itoa(args.limit))
				params.Set("market", "ID")
				params.Set("seed_tracks", args.trackID)

				basePath := `https://api.spotify.com/v1/recommendations`
				endpoint := fmt.Sprintf("%s?%s", basePath, params.Encode())
				req, err := http.NewRequest(http.MethodGet, endpoint, nil)
				assert.NoError(t, err)

				req.Header.Set("Authorization", "Bearer accessToken")
				mockHTTPClient.EXPECT().Do(req).Return(&http.Response{
					StatusCode: 500,
					Body:       io.NopCloser(bytes.NewBufferString(`internal server error`)),
				}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			o := &outbound{
				cfg:         &configs.Config{},
				client:      mockHTTPClient,
				AccessToken: "accessToken",
				TokenType:   "Bearer",
				ExpiredAt:   time.Now().Add(1 * time.Hour),
			}
			got, err := o.GetRecommendation(context.Background(), tt.args.limit, tt.args.trackID)
			if (err != nil) != tt.wantErr {
				t.Errorf("outbound.GetRecommendation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("outbound.GetRecommendation() = %v, want %v", got, tt.want)
			}
		})
	}
}
