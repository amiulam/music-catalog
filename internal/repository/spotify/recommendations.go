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

func (o *outbound) GetRecommendation(ctx context.Context, limit int, trackID string) (*SpotifyRecommendationResponse, error) {
	params := url.Values{}
	params.Set("limit", strconv.Itoa(limit))
	params.Set("market", "ID")
	params.Set("seed_tracks", trackID)

	basePath := `https://api.spotify.com/v1/recommendations`
	endpoint := fmt.Sprintf("%s?%s", basePath, params.Encode())
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)

	if err != nil {
		log.Error().Err(err).Msg("error while create recommendations request to spotify")
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
		log.Error().Err(err).Msg("error while execute recommendation request to spotify")
		return nil, err
	}

	defer resp.Body.Close()

	var response SpotifyRecommendationResponse

	err = json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		log.Error().Err(err).Msg("error unmarshal recommendation response from spotify")
		return nil, err
	}

	return &response, nil
}
