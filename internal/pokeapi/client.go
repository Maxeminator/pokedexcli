package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type LocationAreaResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []NamedAPIResource
}

type NamedAPIResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func GetLocationAreas(url string) (LocationAreaResponse, error) {
	locations := LocationAreaResponse{}

	if url == "" {
		return locations, fmt.Errorf("invalid url")
	}

	resp, err := http.Get(url)
	if err != nil {
		return locations, err
	}

	defer resp.Body.Close()
	if resp.StatusCode > 299 {
		return locations, fmt.Errorf("response failde with status code: %d ", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&locations)
	if err != nil {
		return locations, err
	}
	return locations, nil
}
