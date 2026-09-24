package main

import (
	"encoding/json"
	"net/http"
)

type Station struct {
	Name string `json:"name"`
	URL  string `json:"url_resolved"`
}

func getStations() ([]Station, error) {
	resp, err := http.Get("https://de1.api.radio-browser.info/json/stations/bytag/j-rock")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var stations []Station

	if err := json.NewDecoder(resp.Body).Decode(&stations); err != nil {
		return nil, err
	}
	return stations, nil
}
