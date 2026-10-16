package aqicn

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// GetAreaAqi fetches the current feed for the station with the given id.
func GetAreaAqi(id string, apiToken string) (*FeedResponse, error) {
	u := fmt.Sprintf("https://api.waqi.info/feed/@%s/?token=%s", url.PathEscape(id), url.PathEscape(apiToken))

	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("aqicn: GET %s: %s", u, resp.Status)
	}

	var r FeedResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("aqicn: decode feed response: %w", err)
	}
	return &r, nil
}
