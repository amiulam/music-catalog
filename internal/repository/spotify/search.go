package spotify

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/rs/zerolog/log"
)

type SpotifySearchResponse struct {
	Tracks SpotifyTracks `json:"tracks"`
}

type SpotifyTracks struct {
	Href     string                `json:"href"`
	Limit    int                   `json:"limit"`
	Next     *string               `json:"next"`
	Offset   int                   `json:"offset"`
	Previous *string               `json:"previous"`
	Total    int                   `json:"total"`
	Items    []SpotifyTrackObjects `json:"items"`
}

type SpotifyTrackObjects struct {
	Album    SpotifyAlbumObject    `json:"album"`
	Artists  []SpotifyArtistObject `json:"artists"`
	Explicit bool                  `json:"explicit"`
	Href     string                `json:"href"`
	ID       string                `json:"id"`
	Name     string                `json:"name"`
}

type SpotifyAlbumObject struct {
	AlbumType   string              `json:"album_type"`
	TotalTracks int                 `json:"total_tracks"`
	Images      []SpotifyAlbumImage `json:"images"`
	Name        string              `json:"name"`
}

type SpotifyArtistObject struct {
	Href string `json:"href"`
	Name string `json:"name"`
}

type SpotifyAlbumImage struct {
	URL string `json:"url"`
}

func (o *outbound) Search(ctx context.Context, query string, limit, offset int) (*SpotifySearchResponse, error) {
	params := url.Values{}
	params.Set("q", query)
	params.Set("type", "track")
	params.Set("offset", strconv.Itoa(offset))
	params.Set("limit", strconv.Itoa(limit))

	basePath := `https://api.spotify.com/v1/search`
	endpoint := fmt.Sprintf("%s?%s", basePath, params.Encode())
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)

	if err != nil {
		log.Error().Err(err).Msg("error while create search request to spotify")
		return nil, err
	}

	accessToken, tokenType, err := o.GetTokenDetails()

	if err != nil {
		log.Error().Err(err).Msg("error get token details")
		return nil, err
	}

	bearerToken := fmt.Sprintf("%v %v", tokenType, accessToken)

	req.Header.Set("Authorization", bearerToken)

	resp, err := o.client.Do(req)

	if err != nil {
		log.Error().Err(err).Msg("error while execute search request to spotify")
		return nil, err
	}

	defer resp.Body.Close()

	var response SpotifySearchResponse

	err = json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		log.Error().Err(err).Msg("error unmarshal search response from spotify")
		return nil, err
	}

	return &response, nil
}
