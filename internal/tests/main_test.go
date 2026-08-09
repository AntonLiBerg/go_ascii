package tests

import (
	component "go_ascii/internal"
	"go_ascii/internal/scenario"
	"go_ascii/internal/service/interaction"
	"go_ascii/internal/service/render"
	"go_ascii/internal/world"
	"slices"
	"strings"
	"testing"
)

func TestSkyshipScenarioLoadsRooms(t *testing.T) {
	asciiMap, entities, components, inputProfiles, uiLayout, uis, err := scenario.GetScenarioFromFiles(
		"../../scenarios/skyship/map.txt",
		"../../scenarios/skyship/content.txt",
		"../../scenarios/skyship/ui.txt",
	)
	if err != nil {
		t.Fatalf("GetScenarioFromFiles returned error: %v", err)
	}
	if !slices.Equal(uiLayout, []string{"room", "infobox"}) || !slices.Equal(uis, []string{"infobox"}) {
		t.Fatalf("expected skyship UI layout, got layout=%v UIs=%v", uiLayout, uis)
	}

	gameWorld, err := world.NewWorldWithUI(asciiMap, entities, components, inputProfiles, uiLayout, uis)
	if err != nil {
		t.Fatalf("NewWorld returned error: %v", err)
	}
	if len(asciiMap.Rooms) != 4 || len(asciiMap.Portals) != 2 || len(asciiMap.Terminals) != 2 {
		t.Fatalf("expected four rooms, paired portals, and two terminals, got rooms=%d portals=%d terminals=%d", len(asciiMap.Rooms), len(asciiMap.Portals), len(asciiMap.Terminals))
	}
	if len(gameWorld.Interactable) != 2 {
		t.Fatalf("expected two interactable instruments, got %d", len(gameWorld.Interactable))
	}
	for _, entity := range []rune{'.', '@', 'H', 'T'} {
		if _, ok := asciiMap.EntityGroups["bridge"][entity]; !ok {
			t.Fatalf("expected bridge group composition to include %q", entity)
		}
	}
	if got := gameWorld.UIContent["infobox"]; len(got) != 6 {
		t.Fatalf("expected six infobox lines, got %d", len(got))
	}
	commandTableID := -1
	for entityID, interactable := range gameWorld.Interactable {
		if interactable.InteractionType == component.InteractionTypeTerminal && gameWorld.Pos[entityID].Room == "bridge" {
			commandTableID = entityID
			break
		}
	}
	if commandTableID == -1 || gameWorld.Pos[commandTableID].Room != "bridge" {
		t.Fatalf("expected terminal in bridge, got %+v", gameWorld.Pos[commandTableID])
	}
	if got := inputProfiles["terminal_scan"]["exit"]; got != "q" {
		t.Fatalf("expected terminal exit binding q, got %q", got)
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
	asciiMap, entities, components, inputProfiles, _, _, err := scenario.GetScenarioFromFiles(
		"../../scenarios/skyship/map.txt",
		"../../scenarios/skyship/content.txt",
		"../../scenarios/skyship/ui.txt",
	)
	if err != nil {
		t.Fatalf("GetScenarioFromFiles returned error: %v", err)
	}
	gameWorld, err := world.NewWorld(asciiMap, entities, components, inputProfiles)
	if err != nil {
		t.Fatalf("NewWorld returned error: %v", err)
	}
	selectablesByRoom := make(map[string]int)
	for entityID := range gameWorld.Selectable {
		selectablesByRoom[gameWorld.Pos[entityID].Room]++
	}
	if selectablesByRoom["terminal_helm"] != 2 || selectablesByRoom["terminal_scan"] != 0 {
		t.Fatalf("expected only two Helm selectables, got %v", selectablesByRoom)
	}

	playerID := world.GetPlayerID(gameWorld)
	commandTableID := -1
	for position, targetRoom := range gameWorld.Terminals {
		if targetRoom != "terminal_helm" {
			continue
		}
		for _, entityID := range world.GetEntitiesAtPosition(gameWorld, position) {
			if interactable, ok := gameWorld.Interactable[entityID]; ok && interactable.InteractionType == component.InteractionTypeTerminal {
				commandTableID = entityID
				break
			}
		}
	}
	if commandTableID == -1 {
		t.Fatal("expected player and bridge terminal")
	}
	commandTablePosition := gameWorld.Pos[commandTableID]
	physicalPosition := component.Position{
		Room: commandTablePosition.Room,
		X:    commandTablePosition.X - 1,
		Y:    commandTablePosition.Y,
	}
	gameWorld.Pos[playerID] = physicalPosition
	if !gameWorld.SetInputProfileForRoom(physicalPosition.Room) {
		t.Fatal("expected bridge input profile")
	}
	gameWorld.KeyDown = gameWorld.UserInputProfile.KeyInteract

	open := interaction.ServiceInteraction{}.GetUpdateFunc(gameWorld)
	gameWorld = applyUpdate(t, open, gameWorld)
	if gameWorld.Room != "terminal_helm" || gameWorld.UserInputProfile.KeyExit != "q" {
		t.Fatalf("expected terminal_helm exit profile, got room=%q profile=%+v", gameWorld.Room, gameWorld.UserInputProfile)
	}
	if gameWorld.Pos[playerID] != physicalPosition {
		t.Fatal("expected player to remain on the bridge while terminal is open")
	}
	if gameWorld.ActiveControl != gameWorld.ControlLists["terminal_helm"] {
		t.Fatal("expected first Helm control to be focused")
	}
	frame := render.TerminalFrame(gameWorld)
	if !strings.Contains(frame, "\033[8;9Hö") || !strings.Contains(frame, "\033[8;19Ho") {
		t.Fatalf("expected focused first control and unfocused second control in frame %q", frame)
	}

	gameWorld.KeyDown = gameWorld.UserInputProfile.KeyExit
	closeTerminal := interaction.ServiceInteraction{}.GetUpdateFunc(gameWorld)
	gameWorld = applyUpdate(t, closeTerminal, gameWorld)
	if gameWorld.UserInputProfile.KeyInteract != "e" {
		t.Fatalf("expected topdown profile after exit, got %+v", gameWorld.UserInputProfile)
	}
}
