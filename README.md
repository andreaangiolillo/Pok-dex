# Pokedex CLI

A command-line Pokédex written in Go. It uses the [PokéAPI](https://pokeapi.co/) to browse location areas, discover Pokémon, catch them, and inspect your collection.

## Features

- Browse location areas with `map` and `mapb`
- Explore a location area to see available Pokémon
- Attempt to catch Pokémon
- Inspect caught Pokémon
- List all Pokémon in your Pokédex
- Cache API responses in memory for a limited time

## Requirements

- Go 1.27 or newer
- Internet access for requests to PokéAPI

## Getting started

Clone the repository and run the application:

```bash
git clone https://github.com/andreaangiolillo/Personal.git
cd Personal/pokedexcli
go run .
```

You can also build an executable:

```bash
go build -o pokedex
./pokedex
```

## Commands

| Command | Description |
| --- | --- |
| `help` | Display the available commands |
| `map` | Show the next page of location areas |
| `mapb` | Show the previous page of location areas |
| `explore <location-area>` | List Pokémon found in a location area |
| `catch <pokemon>` | Attempt to catch a Pokémon |
| `inspect <pokemon>` | Show details for a caught Pokémon |
| `pokedex` | List all caught Pokémon |
| `exit` | Quit the application |

Example session:

```text
Pokedex > map
Pokedex > explore pastoria-city-area
Pokedex > catch tentacool
Pokedex > inspect tentacool
Pokedex > pokedex
```

Location-area names are provided by the PokéAPI. You can get a valid name from the output of `map`.

## Project structure

```text
.
├── main.go                 # CLI, commands, and PokéAPI integration
├── repl/                   # Input cleaning and REPL tests
└── internal/pokecache/     # In-memory response cache
```

## Testing

Run all tests with:

```bash
go test ./...
```

## Data source

Pokémon and location data comes from [PokéAPI](https://pokeapi.co/). This project is for educational purposes and is not affiliated with Nintendo, Game Freak, or The Pokémon Company.
