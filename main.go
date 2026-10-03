package main
import (
	"fmt"
	"bufio"
	"os"
	"github.com/andreaangiolillo/Pok-dex/pokedexcli/repl"
	"github.com/andreaangiolillo/Pok-dex/pokedexcli/internal/pokecache"
	"net/http"
	"io"
	"time"
	"encoding/json"
	"math/rand"
)


type command struct {
	name string
	description string
	callback func(*config, *pokecache.Cache, ...string) error
}

type config struct {
	commands map[string]command
	previousLocation string
	nextLocation string
	pokemons map[string]Pokemon
}

type Locations struct {
	Count int `json:"count"`
	Next *string `json:"next"`
	Previous *string `json:"previous"`
	Results []Location `json:"results"`
}

type Location struct {
	Name string `json:"name"`
	URL string  `json:"url"`
}

type PokemonEntry struct {
	Pokemon Entry `json:"pokemon"`
}

type Entry struct {
	Name string `json:"name"`
	URL string  `json:"url"`
}

type Area struct {
	Pokemons []PokemonEntry `json:"pokemon_encounters"`
}

type Pokemon struct {
	Name string `json:"name"`
	Weight int `json:"weight"`
	Height int `json:"height"`
	BaseExperience int `json:"base_experience"`
}



func mapCommand(url string, conf *config, cache *pokecache.Cache) error {
	if url == "" {
	    url = "https://pokeapi.co/api/v2/location-area"
	}
	
	body, ok := cache.Get(url); 
	if !ok {
		res, err := http.Get(url)
		if err != nil {
			return err
		}
		
		body, err = io.ReadAll(res.Body)
		res.Body.Close()
		
		if res.StatusCode > 299 {
			fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		}

		if err != nil {
			return err
		}
	}
		
	locs := Locations{}
	err := json.Unmarshal(body, &locs)
	if err != nil {
		return err
	}
	
	conf.previousLocation = ""
	if locs.Previous != nil {
		conf.previousLocation = *locs.Previous
	}

	conf.nextLocation = ""
	if locs.Next != nil{
		conf.nextLocation = *locs.Next
	}else{
		fmt.Println("\nThis is the last page")
	}

	for _, loc := range locs.Results{
		fmt.Println(loc.Name)
	}

	cache.Add(url, body)

	return nil
}

func exploreCommand(config *config, cache *pokecache.Cache, location string) error{
	url := fmt.Sprintf("%s/%s", "https://pokeapi.co/api/v2/location-area", location)
	body, ok := cache.Get(url)

	if !ok {
		res, err := http.Get(url)
		if err != nil {
			return err
		}
		body, err = io.ReadAll(res.Body)
		res.Body.Close()
			
		if res.StatusCode > 299 {
			fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		}
		
		if err != nil {
			return err
		}
	}

	var area Area
	err := json.Unmarshal(body, &area)
	if err != nil {
		return err
	}
	
	for _, p := range area.Pokemons{
		fmt.Println(p.Pokemon.Name)
	}
	
	cache.Add(url, body)

	return nil
}

func catchCommand(config *config, cache *pokecache.Cache, name string) error{
	url := fmt.Sprintf("%s/%s", "https://pokeapi.co/api/v2/pokemon", name)
	
	body, ok := cache.Get(url)
	if !ok {
		res, err := http.Get(url)
		if err != nil {
			return err
		}

		body, err = io.ReadAll(res.Body)
		res.Body.Close()

		if err != nil {
			return err
		}

		if res.StatusCode > 299 {
			fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		}
	}

	var pokemon Pokemon
	err := json.Unmarshal(body, &pokemon)
	if err != nil{
		return err
	}

	cache.Add(url, body)

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name) 
	total := pokemon.Weight + pokemon.BaseExperience + 100
	var percentage float64 = float64(pokemon.Weight + pokemon.BaseExperience) / float64(total)

	result := rand.Intn(total)
	if float64(result) > (float64(total) * percentage) {
		fmt.Printf("%s was caught!\n", pokemon.Name) 
		config.pokemons[pokemon.Name] = pokemon
		return nil
	}
	
	fmt.Printf("%s escaped!\n", pokemon.Name) 
	return nil
}
	

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	d, _ := time.ParseDuration("1h30m")
	cache := pokecache.NewCache(d)
	conf := config {
			commands: map[string]command{
			"exit": {
				name: "exit",
				description: "Exit the Pokedex",
				callback: func(_ *config, _ *pokecache.Cache, _ ...string) error{
					fmt.Println("Closing the Pokedex... Goodbye!")
					os.Exit(0)
					return nil
				},
			},
			"help":{
				name: "help",
				description: "Displays a help message",
				callback: func(_ *config, _ *pokecache.Cache, _ ...string) error{
					fmt.Println("Welcome to the Pokedex!")
					fmt.Println("Usage:")
					fmt.Println("")
					fmt.Println("help: Displays a help message")
					fmt.Println("exit: Exit the Pokedex")
					return nil
				},
			},
			"map":{
				name: "map",
				description: "Displays the next names of 20 location areas in the Pokemon world",
				callback: func(conf *config, cache *pokecache.Cache, _ ...string) error {
					return mapCommand(conf.nextLocation, conf, cache)
				},
			},
			"mapb":{
				name: "mapb",
				description: "Displays the previous names of 20 location areas in the Pokemon world",
				callback: func(conf *config, cache *pokecache.Cache, _ ...string) error {
					if conf.previousLocation == "" {
						fmt.Println("You're on the first page")
					}
					return mapCommand(conf.previousLocation, conf, cache)
				},
			},
			"explore" : {
				name: "explore",
				description: "Explore a location to find pokemons",
				callback: func(conf *config, cache *pokecache.Cache, args ...string) error {
					return exploreCommand(conf, cache, args[0])
				},
			},
			"catch": {
				name: "catch",
				description: "Catch a Pokemon",
				callback: func(conf *config, cache *pokecache.Cache, args ...string) error {
					return catchCommand(conf, cache, args[0])
				},
			},
			"inspect": {
				name: "inspect",
				description: "Inspect a Pokemon you caught",
				callback: func(conf *config, _ *pokecache.Cache, args ...string) error {
					pokemon := args[0]
					v, ok := conf.pokemons[pokemon]
					if !ok {
						fmt.Println("you have not caught that pokemon")
						return nil
					}
					fmt.Printf("Name: %s\nWeight: %d\nHeight: %d\nExperience: %d\n", v.Name, v.Weight, v.Height, v.BaseExperience)
					return nil
				},
			},
			"pokedex": {
				name: "pokedex",
				description: "Print all Pokemons you caught",
				callback: func(conf *config, _ *pokecache.Cache, _ ...string) error{
					if len(conf.pokemons) == 0 {
						fmt.Println("Your Pokedex is Empty")
						return nil
					}

					fmt.Println("Your Pokedex:")
					for k, _ := range conf.pokemons {
						fmt.Printf("- %s\n", k)
					}
					return nil

				},
			},
		},
		pokemons: map[string]Pokemon{},
	}
	

	for {
		fmt.Print("Pokedex > ")
		_= scanner.Scan()
		text := scanner.Text()
		inputs := repl.CleanInput(text)
		
		for i := 0; i < len(inputs); i++ {
			if command, ok := conf.commands[inputs[i]]; ok {
				switch command.name {
					case "explore", "catch", "inspect":
						if len(inputs) <= i + 1 {
							fmt.Println("Missing argument")
							continue
						}
						i+=1
						arg := string(inputs[i])
						if err := command.callback(&conf, cache, arg); err != nil {
							fmt.Println(err)
						}
				
					default:
						if err := command.callback(&conf, cache); err != nil {
							fmt.Println(err)
						}
		
				}

			}else{
				fmt.Printf("Unknown command: %s\n", inputs[i])
			}
		}
	}
}

