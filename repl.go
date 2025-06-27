package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Maxeminator/pokedexcli/internal/pokeapi"
	"github.com/Maxeminator/pokedexcli/internal/pokecache"
)

type config struct {
	Next     *string
	Previous *string
	Cache    *pokecache.Cache
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

var commands map[string]cliCommand

func init() {
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "View list of areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "View previous list of areas",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Description of selected location",
			callback:    commandExplore,
		},
	}
}

func cleanInput(text string) []string {
	lower := strings.ToLower(text)
	words := strings.Fields(lower)
	return words
}

func commandExit(*config, []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(*config, []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage: ")
	fmt.Println()
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(cfg *config, args []string) error {
	url := "https://pokeapi.co/api/v2/location-area"
	if cfg.Next != nil {
		url = *cfg.Next
	}
	locations, err := pokeapi.GetLocationAreas(url, cfg.Cache)
	if err != nil {
		return fmt.Errorf("failed to get location areas: %w", err)
	}
	if locations.Count == 0 {
		fmt.Println("No location areas found.")
		return nil
	}
	fmt.Println("Location Areas:")
	for _, area := range locations.Results {
		fmt.Println(area.Name)
	}
	cfg.Next = locations.Next
	cfg.Previous = locations.Previous
	if cfg.Next != nil {
		fmt.Printf("Next page: %s\n", *locations.Next)
	}
	if locations.Previous != nil {
		fmt.Printf("Previous page: %s\n", *locations.Previous)
	}
	return nil
}

func commandMapb(cfg *config, args []string) error {
	if cfg.Previous == nil {
		fmt.Println("You're on the first page.")
		return nil
	}
	url := "https://pokeapi.co/api/v2/location-area"

	url = *cfg.Previous

	locations, err := pokeapi.GetLocationAreas(url, cfg.Cache)
	if err != nil {
		return fmt.Errorf("failed to get location areas: %w", err)
	}
	if locations.Count == 0 {
		fmt.Println("No location areas found.")
		return nil
	}
	fmt.Println("Location Areas:")
	for _, area := range locations.Results {
		fmt.Println(area.Name)
	}
	cfg.Next = locations.Next
	cfg.Previous = locations.Previous
	if cfg.Next != nil {
		fmt.Printf("Next page: %s\n", *locations.Next)
	}
	if locations.Previous != nil {
		fmt.Printf("Previous page: %s\n", *locations.Previous)
	}
	return nil
}

func commandExplore(cfg *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("name is empty! Usage: explore <location name>")
	}

	area := args[0]
	fmt.Printf("Exploring %s...\n", area)

	baseUrl := "https://pokeapi.co/api/v2/location-area/"

	data, err := pokeapi.GetLocationAreaDetails(baseUrl, area, cfg.Cache)
	if err != nil {
		return err
	}

	if len(data.PokemonEncounters) == 0 {
		fmt.Println("No pokemons in this area.")
		return nil
	}

	fmt.Println("Found Pokemon:")
	for _, e := range data.PokemonEncounters {
		fmt.Printf(" - %s\n", e.Pokemon.Name)
	}
	return nil
}
