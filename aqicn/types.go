package aqicn

import "encoding/json"

type Attribution struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

type City struct {
	Geo      [2]float64 `json:"geo"`
	Name     string     `json:"name"`
	URL      string     `json:"url"`
	Location string     `json:"location,omitempty"`
}

// IAQI is the latest reading for a single pollutant or weather metric,
// keyed in Msg.IAQI by its code (pm25, pm10, o3, co, so2, t, h, p, w, dew).
type IAQI struct {
	V float64 `json:"v"`
	T string  `json:"t"`
}

type Time struct {
	S   string `json:"s"`
	Tz  string `json:"tz"`
	V   int64  `json:"v"`
	ISO string `json:"iso"`
}

// Series is the delta-encoded history for a single metric. V starts with the
// latest value (scaled by M) followed by [valueDelta, secondsDelta] pairs
// going back in time from S.
type Series struct {
	E int               `json:"e"`
	S string            `json:"s"`
	D int               `json:"d"`
	M int               `json:"m"`
	V []json.RawMessage `json:"v"`
}

type Forecast struct {
	Hourly json.RawMessage `json:"hourly"`
}

type XSync struct {
	Gen int64 `json:"gen"`
}

// FeedResponse is the payload of the token-authenticated feed endpoint
// https://api.waqi.info/feed/@{idx}/?token={token}.
type FeedResponse struct {
	Status string   `json:"status"`
	Data   FeedData `json:"data"`
}

type FeedData struct {
	AQI          int             `json:"aqi"`
	Idx          int             `json:"idx"`
	Attributions []Attribution   `json:"attributions"`
	City         City            `json:"city"`
	DominentPol  string          `json:"dominentpol"`
	IAQI         map[string]IAQI `json:"iaqi"`
	Time         Time            `json:"time"`
	Forecast     FeedForecast    `json:"forecast"`
	Debug        FeedDebug       `json:"debug"`
}

type FeedForecast struct {
	// Daily is keyed by pollutant code (pm25, pm10, o3, uvi).
	Daily map[string][]DailyForecast `json:"daily"`
}

type DailyForecast struct {
	Avg int    `json:"avg"`
	Day string `json:"day"`
	Max int    `json:"max"`
	Min int    `json:"min"`
}

type FeedDebug struct {
	Sync string `json:"sync"`
}
