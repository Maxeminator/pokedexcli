package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Maxeminator/pokedexcli/internal/pokecache"
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

func GetLocationAreas(url string, cache *pokecache.Cache) (LocationAreaResponse, error) {
	var locations LocationAreaResponse

	val, ok := cache.Get(url)
	if ok {
		err := json.Unmarshal(val, &locations)
		if err != nil {
			return locations, err
		}
		fmt.Println("from cache")
		return locations, nil
	}

	if url == "" {
		return locations, fmt.Errorf("invalid url")
	}

	resp, err := http.Get(url)
	if err != nil {
		return locations, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return locations, err
	}

	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return locations, fmt.Errorf("response failde with status code: %d ", resp.StatusCode)
	}

	cache.Add(url, body)

	err = json.Unmarshal(body, &locations)
	if err != nil {
		return locations, err
	}
	fmt.Println("from network")
	return locations, nil
}
