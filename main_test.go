package main

import (
	"slices"
	"testing"

	"github.com/aikchun/aqicnsg-telegram/aqicnsg"
)

func TestMaxAQI(t *testing.T) {
	tests := []struct {
		name  string
		areas []aqicnsg.AreaAqi
		want  int
	}{
		{"single", []aqicnsg.AreaAqi{{AQI: 42}}, 42},
		{"max first", []aqicnsg.AreaAqi{{AQI: 150}, {AQI: 20}, {AQI: 80}}, 150},
		{"max middle", []aqicnsg.AreaAqi{{AQI: 20}, {AQI: 150}, {AQI: 80}}, 150},
		{"max last", []aqicnsg.AreaAqi{{AQI: 20}, {AQI: 80}, {AQI: 150}}, 150},
		{"all equal", []aqicnsg.AreaAqi{{AQI: 55}, {AQI: 55}}, 55},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maxAQI(tt.areas); got != tt.want {
				t.Errorf("maxAQI() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestFormatPM25Report(t *testing.T) {
	areas := []aqicnsg.AreaAqi{
		{City: "North, Singapore", AQI: 42, Time: "2026-10-18 14:00:00"},
		{City: "South, Singapore", AQI: 101, Time: "2026-10-18 14:00:00"},
		{City: "Central", AQI: 350, Time: "2026-10-18 14:00:00"},
	}
	want := "PM2.5 AQI at Oct 18 2:00PM:\n" +
		"North: 42 🟢\n" +
		"South: 101 🟠\n" +
		"Central: 350 ⚫"

	got, err := formatPM25Report(areas)
	if err != nil {
		t.Fatalf("formatPM25Report() error = %v", err)
	}
	if got != want {
		t.Errorf("formatPM25Report() =\n%q\nwant\n%q", got, want)
	}
}

func TestFormatPM25ReportInvalidTime(t *testing.T) {
	areas := []aqicnsg.AreaAqi{{City: "North", AQI: 42, Time: "not a time"}}
	if _, err := formatPM25Report(areas); err == nil {
		t.Error("formatPM25Report() error = nil, want error")
	}
}

func TestParseText(t *testing.T) {
	tests := []struct {
		name        string
		text        string
		wantCommand string
		wantArgs    []string
		wantOK      bool
	}{
		{"command", "/now", "now", []string{}, true},
		{"command with bot name", "/now@aqibot", "now", []string{}, true},
		{"command with args", "/subscribe 80 extra", "subscribe", []string{"80", "extra"}, true},
		{"bot name and args", "/subscribe@aqibot 80", "subscribe", []string{"80"}, true},
		{"not a command", "hello world", "", []string{}, false},
		{"empty", "", "", []string{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command, args, ok := parseText(tt.text)
			if command != tt.wantCommand || !slices.Equal(args, tt.wantArgs) || ok != tt.wantOK {
				t.Errorf("parseText(%q) = (%q, %q, %t), want (%q, %q, %t)",
					tt.text, command, args, ok, tt.wantCommand, tt.wantArgs, tt.wantOK)
			}
		})
	}
}

func TestAqiEmoji(t *testing.T) {
	tests := []struct {
		aqi  int
		want string
	}{
		{0, "🟢"},
		{50, "🟢"},
		{51, "🟡"},
		{100, "🟡"},
		{101, "🟠"},
		{150, "🟠"},
		{151, "🔴"},
		{200, "🔴"},
		{201, "🟣"},
		{300, "🟣"},
		{301, "⚫"},
		{999, "⚫"},
	}
	for _, tt := range tests {
		if got := aqiEmoji(tt.aqi); got != tt.want {
			t.Errorf("aqiEmoji(%d) = %s, want %s", tt.aqi, got, tt.want)
		}
	}
}
