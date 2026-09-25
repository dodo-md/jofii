package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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
	seen := map[string]bool{}
	for _, s := range raw {
		stream := s.URLResolved
		if stream == "" {
			stream = s.URL
		}
		if s.Name == "" || stream == "" || seen[stream] || !audible(stream) {
			continue
		}
		seen[stream] = true
		stations = append(stations, Station{Name: s.Name, URL: stream})
	}
	return stations, nil
}

func audible(stream string) bool {
	req, err := http.NewRequest(http.MethodGet, stream, nil)
	if err != nil {
		return false
	}
	req.Header.Set("User-Agent", "jofii/0.1.0")
	req.Header.Set("Icy-MetaData", "0")

	client := &http.Client{Timeout: 6 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false
	}

	buf := make([]byte, 512)
	n, _ := resp.Body.Read(buf)
	return looksLikeAudio(resp.Header.Get("Content-Type"), resp.Header, buf[:n])
}

func looksLikeAudio(contentType string, header http.Header, sample []byte) bool {
	body := bytes.TrimSpace(sample)
	if bytes.HasPrefix(bytes.ToLower(body), []byte("<")) {
		return false
	}
	ct := strings.ToLower(contentType)
	if strings.Contains(ct, "text/html") || strings.Contains(ct, "text/plain") {
		return false
	}
	if header.Get("icy-name") != "" || header.Get("icy-br") != "" {
		return true
	}
	if strings.HasPrefix(ct, "audio/") || strings.Contains(ct, "ogg") || strings.Contains(ct, "flac") || strings.Contains(ct, "mpegurl") || strings.Contains(ct, "scpls") {
		return true
	}
	if bytes.HasPrefix(body, []byte("ID3")) || bytes.HasPrefix(body, []byte("OggS")) || bytes.HasPrefix(body, []byte("fLaC")) || bytes.HasPrefix(body, []byte("#EXTM3U")) {
		return true
	}
	return len(body) >= 2 && body[0] == 0xFF && body[1]&0xE0 == 0xE0
}
