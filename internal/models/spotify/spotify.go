package spotify

type SearchResponse struct {
	Limit  int                   `json:"limit"`
	Offset int                   `json:"offset"`
	Items  []SpotifyTrackObjects `json:"items"`
	Total  int                   `json:"total"`
}

type SpotifyTrackObjects struct {
	// Album fields
	AlbumType         string   `json:"albumType"`
	AlbumnTotalTracks int      `json:"totalTracks"`
	AlbumImagesURL    []string `json:"albumImagesURL"`
	AlbumName         string   `json:"albumName"`

	// Artists field
	ArtistsName []string `json:"artistsName"`

	// Track fields
	Explicit bool   `json:"explicit"`
	ID       string `json:"id"`
	Name     string `json:"name"`
	IsLiked  *bool  `json:"isLiked"`
}

type RecommendationResponse struct {
	Items []SpotifyTrackObjects `json:"items"`
}
