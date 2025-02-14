package tracks

import (
	"context"
	"reflect"
	"testing"

	"github.com/amiulam/music-catalog/internal/models/spotify"
	trackactivities "github.com/amiulam/music-catalog/internal/models/track_activities"
	spotifyRepo "github.com/amiulam/music-catalog/internal/repository/spotify"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func Test_service_Search(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockSpotifyOutbound := NewMockspotifyOutbound(ctrlMock)
	mockTrackActivityRepo := NewMocktrackActivitesRepository(ctrlMock)

	isLikedTrue := true
	isLikedFalse := false

	next := "https://api.spotify.com/v1/search?offset=10&limit=10&query=bohemian%20rhapsody&type=track&market=ID&locale=id-ID,id;q%3D0.7"

	type args struct {
		query     string
		pageSize  int
		pageIndex int
	}

	tests := []struct {
		name    string
		args    args
		want    *spotify.SearchResponse
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "success",
			args: args{
				query:     "bohemian rhapsody",
				pageSize:  10,
				pageIndex: 1,
			},
			want: &spotify.SearchResponse{
				Limit: 10,
				Items: []spotify.SpotifyTrackObjects{
					{
						AlbumType:         "album",
						AlbumnTotalTracks: 22,
						AlbumImagesURL:    []string{"https://i.scdn.co/image/ab67616d0000b273e8b066f70c206551210d902b", "https://i.scdn.co/image/ab67616d00001e02e8b066f70c206551210d902b", "https://i.scdn.co/image/ab67616d00004851e8b066f70c206551210d902b"},
						AlbumName:         "Bohemian Rhapsody (The Original Soundtrack)",
						ArtistsName:       []string{"Queen"},
						Explicit:          false,
						ID:                "3z8h0TU7ReDPLIbEnYhWZb",
						Name:              "Bohemian Rhapsody",
						IsLiked:           &isLikedTrue,
					},
					{
						AlbumType:         "compilation",
						AlbumnTotalTracks: 17,
						AlbumImagesURL:    []string{"https://i.scdn.co/image/ab67616d0000b273bb19d0c22d5709c9d73c8263", "https://i.scdn.co/image/ab67616d00001e02bb19d0c22d5709c9d73c8263", "https://i.scdn.co/image/ab67616d00004851bb19d0c22d5709c9d73c8263"},
						AlbumName:         "Greatest Hits (Remastered)",
						ArtistsName:       []string{"Queen"},
						Explicit:          false,
						ID:                "2OBofMJx94NryV2SK8p8Zf",
						Name:              "Bohemian Rhapsody - Remastered 2011",
						IsLiked:           &isLikedFalse,
					},
				},
				Total: 904,
			},
			wantErr: false,
			mockFn: func(args args) {
				mockSpotifyOutbound.EXPECT().Search(gomock.Any(), args.query, 10, 0).Return(&spotifyRepo.SpotifySearchResponse{
					Tracks: spotifyRepo.SpotifyTracks{
						Href:   "https://api.spotify.com/v1/search?offset=0&limit=10&query=bohemian%20rhapsody&type=track&market=ID&locale=id-ID,id;q%3D0.7",
						Limit:  10,
						Next:   &next,
						Offset: 0,
						Total:  904,
						Items: []spotifyRepo.SpotifyTrackObjects{
							{
								Album: spotifyRepo.SpotifyAlbumObject{
									AlbumType:   "album",
									TotalTracks: 22,
									Images: []spotifyRepo.SpotifyAlbumImage{
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
								Artists: []spotifyRepo.SpotifyArtistObject{
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
								Album: spotifyRepo.SpotifyAlbumObject{
									AlbumType:   "compilation",
									TotalTracks: 17,
									Images: []spotifyRepo.SpotifyAlbumImage{
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
								Artists: []spotifyRepo.SpotifyArtistObject{
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
				}, nil)

				mockTrackActivityRepo.EXPECT().GetBulkSpotifyIDs(gomock.Any(), uint(1), []string{"3z8h0TU7ReDPLIbEnYhWZb", "2OBofMJx94NryV2SK8p8Zf"}).Return(map[string]trackactivities.TrackActivity{
					"3z8h0TU7ReDPLIbEnYhWZb": {
						IsLiked: &isLikedTrue,
					},
					"2OBofMJx94NryV2SK8p8Zf": {
						IsLiked: &isLikedFalse,
					},
				}, nil)
			},
		},
		{
			name: "failed",
			args: args{
				query:     "bohemian rhapsody",
				pageSize:  10,
				pageIndex: 1,
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				mockSpotifyOutbound.EXPECT().Search(gomock.Any(), args.query, 10, 0).Return(nil, assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)

			s := &service{
				spotifyOutbound:    mockSpotifyOutbound,
				trackActivitesRepo: mockTrackActivityRepo,
			}
			got, err := s.Search(context.Background(), tt.args.query, tt.args.pageSize, tt.args.pageIndex, 1)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.Search() = %v, want %v", got, tt.want)
			}
		})
	}
}
