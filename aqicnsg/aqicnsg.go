package aqicnsg

import (
	"fmt"

	"github.com/aikchun/aqicnsg-telegram/aqicn"
)

// areaIDs are the station ids for North, South, East, West and Central Singapore.
var areaIDs = []string{"1662", "1663", "1664", "1665", "1666"}

type AreaAqi struct {
	City string
	AQI  int
	Time string
	Tz   string
}

// GetAreaAqis fetches the current AQI for every Singapore area. It fails if
// any of the requests fail.
func GetAreaAqis(apiToken string) ([]AreaAqi, error) {
	areas := make([]AreaAqi, 0, len(areaIDs))

	for _, id := range areaIDs {
		resp, err := aqicn.GetAreaAqi(id, apiToken)
		if err != nil {
			return nil, fmt.Errorf("sgaqicn: area %s: %w", id, err)
		}

		data := resp.Data

		areas = append(areas, AreaAqi{
			City: data.City.Name,
			AQI:  data.AQI,
			Time: data.Time.S,
			Tz:   data.Time.Tz,
		})
	}

	return areas, nil
}
