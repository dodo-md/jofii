package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Station struct {
	Name string
	URL  string
}

func getStations(limit, offset int) ([]Station, error) {
	u, err := url.Parse("https://de1.api.radio-browser.info/json/stations/search")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("tag", "j-rock")
	q.Set("limit", strconv.Itoa(limit))
	q.Set("offset", strconv.Itoa(offset))
	q.Set("hidebroken", "true")
	q.Set("order", "clickcount")
	q.Set("reverse", "true")
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "jofii/0.1.0")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("radio browser: %s", resp.Status)
	}

	var raw []struct {
		Name        string `json:"name"`
		URL         string `json:"url"`
		URLResolved string `json:"url_resolved"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	stations := make([]Station, 0, len(raw))
	for _, s := range raw {
		stream := s.URLResolved
		if stream == "" {
			stream = s.URL
		}
		if s.Name == "" || stream == "" {
			continue
		}
		stations = append(stations, Station{Name: s.Name, URL: stream})
	}
	return stations, nil
}
