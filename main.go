package main

import (
	"fmt"
)

type Station struct {
	Name string `json:"name"`
	URL  string `json:"url_resolved"`
}

func main() {
	s := Station{
		Name: "kessoku band radio",
		URL:  "https://example.com/stream",
	}

	fmt.Println("station name:", s.Name)
	fmt.Println("stream url:", s.URL)
}
