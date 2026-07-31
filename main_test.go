package main

import (
	component "go_ascii/internal"
	"go_ascii/internal/scenario"
	"go_ascii/internal/world"
	"testing"
)

func TestDemoScenarioLoadsAllMapEntities(t *testing.T) {
	asciiMap, entities, components, _, err := scenario.GetScenarioFromFiles(
		"./scenarios/demo/map.txt",
		"./scenarios/demo/content.txt",
	)
	if err != nil {
		t.Fatalf("GetScenarioFromFiles returned error: %v", err)
	}

	gameWorld, err := world.NewWorld(asciiMap, entities, components)
	if err != nil {
		t.Fatalf("NewWorld returned error: %v", err)
	}
	if len(gameWorld.Helm) != 1 || len(gameWorld.CommandTable) != 1 {
		t.Fatalf("expected one helm and command table, got helm=%d commandTable=%d", len(gameWorld.Helm), len(gameWorld.CommandTable))
	}
	if len(gameWorld.BunkBed) == 0 || len(gameWorld.PrisonBars) == 0 || len(gameWorld.Wall) == 0 || len(gameWorld.Window) == 0 {
		t.Fatal("expected demo structure components to be populated")
	}
	if len(asciiMap.Rooms) != 2 || asciiMap.Ground != "floor" {
		t.Fatalf("expected two rooms with floor ground, got rooms=%d ground=%q", len(asciiMap.Rooms), asciiMap.Ground)
	}
	shipPortal := component.Position{Room: "ship", X: 19, Y: 9}
	enginePortal := component.Position{Room: "engine room", X: 3, Y: 0}
	if asciiMap.Portals[shipPortal] != enginePortal || asciiMap.Portals[enginePortal] != shipPortal {
		t.Fatalf("expected paired demo portals, got %v", asciiMap.Portals)
	}
	for playerID := range gameWorld.Player {
		if gameWorld.Pos[playerID].Room != "ship" {
			t.Fatalf("expected player to start in ship, got %+v", gameWorld.Pos[playerID])
		}
	}
}

func TestSkyshipScenarioLoadsRooms(t *testing.T) {
	asciiMap, entities, components, _, err := scenario.GetScenarioFromFiles(
		"./scenarios/skyship/map.txt",
		"./scenarios/skyship/content.txt",
	)
	if err != nil {
		t.Fatalf("GetScenarioFromFiles returned error: %v", err)
	}

	gameWorld, err := world.NewWorld(asciiMap, entities, components)
	if err != nil {
		t.Fatalf("NewWorld returned error: %v", err)
	}
	if len(asciiMap.Rooms) != 2 || len(asciiMap.Portals) != 2 {
		t.Fatalf("expected two rooms and paired portals, got rooms=%d portals=%d", len(asciiMap.Rooms), len(asciiMap.Portals))
	}
	if len(gameWorld.Helm) != 1 || len(gameWorld.CommandTable) != 1 {
		t.Fatalf("expected skyship instruments, got helm=%d commandTable=%d", len(gameWorld.Helm), len(gameWorld.CommandTable))
	}
	for commandTableID := range gameWorld.CommandTable {
		if gameWorld.Pos[commandTableID].Room != "comms" {
			t.Fatalf("expected command table in comms, got %+v", gameWorld.Pos[commandTableID])
		}
	}
}
