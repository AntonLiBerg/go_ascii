# go_ascii

A small turn-based terminal game engine written in Go.

## Run

```sh
go run .
```

## Test

```sh
go test ./...
```

## Structure

- `internal/world`: world state and components.
- `internal/game`: update scheduling.
- `internal/service`: pure game updates and rendering.
- `internal/scenario`: scenario parsing.
- `scenarios`: maps, entities, input profiles, and UI layouts.

See [`scenarios/README.md`](scenarios/README.md) for the scenario format.
