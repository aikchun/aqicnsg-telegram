package aqicn

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// GetAreaAqi fetches the current feed for the station with the given id.
func GetAreaAqi(id string) (*AqiResponse, error) {
	u := fmt.Sprintf("https://api2.waqi.info/api/feed/@H%s/now.json", url.PathEscape(id))

	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("aqicn: GET %s: %s", u, resp.Status)
	}

	var r AqiResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("aqicn: decode response: %w", err)
	}
	return &r, nil
}
