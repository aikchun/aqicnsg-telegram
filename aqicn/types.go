package aqicn

import "encoding/json"

type AqiResponse struct {
	Dt  string `json:"dt"`
	Rxs Rxs    `json:"rxs"`
}

type Rxs struct {
	Obs    []Observation `json:"obs"`
	Status string        `json:"status"`
	Ver    string        `json:"ver"`
}

// Observation is one entry of rxs.obs. Entries with a status other than "ok"
// (e.g. "nug") carry no Msg.
type Observation struct {
	Status string `json:"status"`
	Msg    *Msg   `json:"msg,omitempty"`
	UID    string `json:"uid,omitempty"`
	Cached string `json:"cached,omitempty"`
}

type Msg struct {
	AQI          int               `json:"aqi"`
	Idx          int               `json:"idx"`
	Attributions []Attribution     `json:"attributions"`
	City         City              `json:"city"`
	DominentPol  string            `json:"dominentpol"`
	IAQI         map[string]IAQI   `json:"iaqi"`
	Time         Time              `json:"time"`
	Obs          map[string]Series `json:"obs"`
	Forecast     Forecast          `json:"forecast"`
	XSync        XSync             `json:"xsync"`
}

type Attribution struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

type City struct {
	Geo  [2]float64 `json:"geo"`
	Name string     `json:"name"`
	URL  string     `json:"url"`
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

func main() {}
