package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/Maxeminator/pokedexcli/internal/pokeapi"
	"github.com/Maxeminator/pokedexcli/internal/pokecache"
)

type config struct {
	Next     *string
	Previous *string
	Cache    *pokecache.Cache
	Pokedex  map[string]pokeapi.Pokemon
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
		"catch": {
			name:        "catch",
			description: "Attempt to catch pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspecting a pokemon in your pokedex",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "View your pokemons",
			callback:    commandPokedex,
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

func commandCatch(cfg *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("name is empty! Usage: catch <pokemon name>")
	}

	name := args[0]
	fmt.Printf("Throwing a Pokeball at %s...\n", name)

	data, err := pokeapi.GetPokemon(name, cfg.Cache)
	if err != nil {
		return err
	}

	if _, ok := cfg.Pokedex[data.Name]; ok {
		fmt.Printf("You already caught %s.\n", data.Name)
		return nil
	}
	chanceToCatch := 50.0 / (1.0 + float64(data.BaseExperience)/50.0)
	chanceToCatch /= 100.0

	roll := rand.Float64()
	if roll < float64(chanceToCatch) {
		fmt.Printf("Caught %s!\n", data.Name)
		fmt.Println("You may now inspect it with the inspect command.")
		cfg.Pokedex[data.Name] = data
		return nil
	} else {
		fmt.Printf("%s escaped!\n", data.Name)
		return nil
	}
}

func commandInspect(cfg *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("name is empty! Usage: inspect <pokemon name>")
	}
	name := args[0]

	if pokemon, ok := cfg.Pokedex[name]; ok {
		fmt.Println("Name: ", pokemon.Name)
		fmt.Println("Height: ", pokemon.Height)
		fmt.Println("Weight: ", pokemon.Weight)
		fmt.Println("Stats:")
		for _, s := range pokemon.Stats {
			fmt.Printf("  - %s: %d\n", s.Stat.Name, s.BaseStat)
		}
		fmt.Println("Types:")
		for _, t := range pokemon.Types {
			fmt.Printf("  - %s\n", t.Type.Name)
		}
		return nil
	} else {
		fmt.Println("you have not caught that pokemon")
		return nil
	}
}

func commandPokedex(cfg *config, args []string) error {
	if len(cfg.Pokedex) == 0 {
		fmt.Println("You haven't caught any Pokemon yet.")
		return nil
	}
	fmt.Println("Your Pokedex:")
	for pokemon := range cfg.Pokedex {
		fmt.Println("-", pokemon)
	}
	return nil
}
