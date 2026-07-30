package main

import (
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
}
