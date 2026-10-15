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
func GetAreaAqis() ([]AreaAqi, error) {
	areas := make([]AreaAqi, 0, len(areaIDs))

	for _, id := range areaIDs {
		resp, err := aqicn.GetAreaAqi(id)
		if err != nil {
			return nil, fmt.Errorf("sgaqicn: area %s: %w", id, err)
		}

		msg := firstMsg(resp)
		if msg == nil {
			return nil, fmt.Errorf("sgaqicn: area %s: no observation in response", id)
		}

		areas = append(areas, AreaAqi{
			City: msg.City.Name,
			AQI:  msg.AQI,
			Time: msg.Time.S,
			Tz:   msg.Time.Tz,
		})
	}

	return areas, nil
}

func firstMsg(resp *aqicn.AqiResponse) *aqicn.Msg {
	for _, obs := range resp.Rxs.Obs {
		if obs.Msg != nil {
			return obs.Msg
		}
	}
	return nil
}
