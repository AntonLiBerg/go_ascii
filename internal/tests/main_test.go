package tests

import (
	component "go_ascii/internal"
	"go_ascii/internal/scenario"
	"go_ascii/internal/service/interaction"
	"go_ascii/internal/world"
	"testing"
)

func TestDemoScenarioLoadsAllMapEntities(t *testing.T) {
	asciiMap, entities, components, inputProfiles, err := scenario.GetScenarioFromFiles(
		"../../scenarios/demo/map.txt",
		"../../scenarios/demo/content.txt",
	)
	if err != nil {
		t.Fatalf("GetScenarioFromFiles returned error: %v", err)
	}

	gameWorld, err := world.NewWorld(asciiMap, entities, components, inputProfiles)
	if err != nil {
		t.Fatalf("NewWorld returned error: %v", err)
	}
	if len(gameWorld.Interactable) != 4 {
		t.Fatalf("expected four interactable entities, got %d", len(gameWorld.Interactable))
	}
	if len(asciiMap.Rooms) != 2 || asciiMap.Ground["ship"] != "floor" || asciiMap.Ground["engine room"] != "floor" {
		t.Fatalf("expected two rooms with floor ground, got rooms=%d ground=%v", len(asciiMap.Rooms), asciiMap.Ground)
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
	asciiMap, entities, components, inputProfiles, err := scenario.GetScenarioFromFiles(
		"../../scenarios/skyship/map.txt",
		"../../scenarios/skyship/content.txt",
	)
	if err != nil {
		t.Fatalf("GetScenarioFromFiles returned error: %v", err)
	}

	gameWorld, err := world.NewWorld(asciiMap, entities, components, inputProfiles)
	if err != nil {
		t.Fatalf("NewWorld returned error: %v", err)
	}
	if len(asciiMap.Rooms) != 3 || len(asciiMap.Portals) != 2 || len(asciiMap.Terminals) != 1 {
		t.Fatalf("expected three rooms, paired portals, and a terminal, got rooms=%d portals=%d terminals=%d", len(asciiMap.Rooms), len(asciiMap.Portals), len(asciiMap.Terminals))
	}
	if len(gameWorld.Interactable) != 2 {
		t.Fatalf("expected two interactable instruments, got %d", len(gameWorld.Interactable))
	}
	commandTableID := -1
	for entityID, interactable := range gameWorld.Interactable {
		if interactable.InteractionType == component.InteractionTypeCommandTable {
			commandTableID = entityID
			break
		}
	}
	if commandTableID == -1 || gameWorld.Pos[commandTableID].Room != "comms" {
		t.Fatalf("expected command table in comms, got %+v", gameWorld.Pos[commandTableID])
	}
	if got := inputProfiles["terminal_scan"]["exit"]; got != "e" {
		t.Fatalf("expected terminal exit binding e, got %q", got)
	}
	terminalGround := component.Position{Room: "terminal_scan", X: 0, Y: 0}
	foundSpace := false
	for _, entityID := range world.GetEntitiesAtPosition(gameWorld, terminalGround) {
		if gameWorld.Layer[entityID].Nr == 0 && gameWorld.Ascii[entityID].Ascii == ' ' {
			foundSpace = true
			break
		}
	}
	if !foundSpace {
		t.Fatal("expected SPACE keyword to create blank terminal ground")
	}
}

func TestSkyshipCommandTerminalFlow(t *testing.T) {
	asciiMap, entities, components, inputProfiles, err := scenario.GetScenarioFromFiles(
		"../../scenarios/skyship/map.txt",
		"../../scenarios/skyship/content.txt",
	)
	if err != nil {
		t.Fatalf("GetScenarioFromFiles returned error: %v", err)
	}
	gameWorld, err := world.NewWorld(asciiMap, entities, components, inputProfiles)
	if err != nil {
		t.Fatalf("NewWorld returned error: %v", err)
	}

	playerID := -1
	for entityID := range gameWorld.Player {
		playerID = entityID
		break
	}
	commandTableID := -1
	for entityID, interactable := range gameWorld.Interactable {
		if interactable.InteractionType == component.InteractionTypeCommandTable {
			commandTableID = entityID
			break
		}
	}
	if playerID == -1 || commandTableID == -1 {
		t.Fatal("expected player and command table")
	}
	commandTablePosition := gameWorld.Pos[commandTableID]
	physicalPosition := component.Position{
		Room: commandTablePosition.Room,
		X:    commandTablePosition.X - 1,
		Y:    commandTablePosition.Y,
	}
	gameWorld.Pos[playerID] = physicalPosition
	if !gameWorld.SetInputProfileForRoom(physicalPosition.Room) {
		t.Fatal("expected comms input profile")
	}
	gameWorld.KeyDown = gameWorld.UserInputProfile.KeyInteract

	open := interaction.ServiceInteraction{}.GetUpdateFunc(gameWorld)
	gameWorld = applyUpdate(t, open, gameWorld)
	if gameWorld.ViewRoom != "terminal_scan" || gameWorld.UserInputProfile.KeyExit != "e" {
		t.Fatalf("expected terminal_scan view and exit profile, got view=%q profile=%+v", gameWorld.ViewRoom, gameWorld.UserInputProfile)
	}
	if gameWorld.Pos[playerID] != physicalPosition {
		t.Fatal("expected player to remain in comms while terminal is open")
	}

	gameWorld.KeyDown = gameWorld.UserInputProfile.KeyExit
	closeTerminal := interaction.ServiceInteraction{}.GetUpdateFunc(gameWorld)
	gameWorld = applyUpdate(t, closeTerminal, gameWorld)
	if gameWorld.ViewRoom != "" || gameWorld.UserInputProfile.KeyInteract != "e" {
		t.Fatalf("expected comms view and topdown profile after exit, got view=%q profile=%+v", gameWorld.ViewRoom, gameWorld.UserInputProfile)
	}
}
