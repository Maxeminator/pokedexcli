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

type LocationArea struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func GetLocationAreaDetails(baseURL string, name string, cache *pokecache.Cache) (LocationArea, error) {
	var pokemons LocationArea
	fullURL := baseURL + name

	val, ok := cache.Get(fullURL)
	if ok {
		err := json.Unmarshal(val, &pokemons)
		if err != nil {
			return pokemons, err
		}
		return pokemons, nil
	}

	if baseURL == "" || name == "" {
		return pokemons, fmt.Errorf("invalid url or name")
	}

	resp, err := http.Get(fullURL)
	if err != nil {
		return pokemons, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return pokemons, err
	}

	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return pokemons, fmt.Errorf("response failed with status code: %d", resp.StatusCode)
	}

	cache.Add(fullURL, body)

	err = json.Unmarshal(body, &pokemons)
	if err != nil {
		return pokemons, err
	}

	return pokemons, nil
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
		return locations, fmt.Errorf("response failed with status code: %d ", resp.StatusCode)
	}

	cache.Add(url, body)

	err = json.Unmarshal(body, &locations)
	if err != nil {
		return locations, err
	}
	fmt.Println("from network")
	return locations, nil
}
