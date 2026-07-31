package world

import (
	component "go_ascii/internal"
	"go_ascii/internal/scenario"
	"testing"
)

func TestAddEntityStoresComponents(t *testing.T) {
	gameWorld := NewWorldEmpty()

	err := gameWorld.AddEntity([2]int{2, 3}, map[string][]string{
		"pos":        {},
		"ascii":      {"å"},
		"impassable": {},
		"door":       {},
		"player":     {},
	})
	if err != nil {
		t.Fatalf("AddEntity returned error: %v", err)
	}

	if len(gameWorld.Entities) != 1 || gameWorld.Entities[0] != 0 {
		t.Fatalf("expected entity 0, got %v", gameWorld.Entities)
	}
	if got := gameWorld.Pos[0]; got != (component.Position{X: 2, Y: 3}) {
		t.Fatalf("expected position 2,3, got %+v", got)
	}
	if got := gameWorld.EByPos[component.Position{X: 2, Y: 3}]; got != 0 {
		t.Fatalf("expected reverse position index to point to entity 0, got %d", got)
	}
	if got := gameWorld.Ascii[0].Ascii; got != 'å' {
		t.Fatalf("expected glyph 'å', got %q", got)
	}
	if _, ok := gameWorld.Impassable[0]; !ok {
		t.Fatal("expected impassable component")
	}
	if _, ok := gameWorld.Player[0]; !ok {
		t.Fatal("expected player component")
	}
	if door, ok := gameWorld.Door[0]; !ok || !door.IsInteractable {
		t.Fatalf("expected an interactable door, got %+v, exists=%t", door, ok)
	}
}

func TestAddEntityRejectsUnknownComponent(t *testing.T) {
	gameWorld := NewWorldEmpty()

	if err := gameWorld.AddEntity([2]int{}, map[string][]string{"visible": {}}); err == nil {
		t.Fatal("expected unknown component error")
	}
}

func TestNewWorldBuildsGroundAndRoomEntities(t *testing.T) {
	gameWorld, err := NewWorld(
		scenario.Map{
			Rooms:         map[string]map[[2]int]rune{"bridge": {{1, 1}: 'o'}},
			Ground:        map[string]string{"bridge": "floor"},
			InputProfiles: map[string]string{"bridge": "topdown"},
			Portals:       make(map[component.Position]component.Position),
		},
		map[rune]string{'.': "floor", 'o': "player"},
		map[string]map[string][]string{
			"floor":  {"pos": {}, "ascii": {"."}},
			"player": {"pos": {}, "ascii": {"o"}, "player": {}},
		},
		map[string]map[string]string{"topdown": {"moveup": "w"}},
	)
	if err != nil {
		t.Fatalf("NewWorld returned error: %v", err)
	}

	if gameWorld.Layer[0].Nr != 0 || gameWorld.Layer[1].Nr != 1 {
		t.Fatalf("expected entity layers 0 and 1, got %v", gameWorld.Layer)
	}
	if gameWorld.Pos[0] != gameWorld.Pos[1] {
		t.Fatalf("expected layered entities at the same position, got %v", gameWorld.Pos)
	}
	if gameWorld.Pos[0].Room != "bridge" {
		t.Fatalf("expected entities in bridge, got %+v", gameWorld.Pos)
	}
	if gameWorld.UserInputProfile.KeyMoveUp != "w" {
		t.Fatalf("expected bridge input profile, got %+v", gameWorld.UserInputProfile)
	}
}

func TestNewWorldRejectsUndefinedMapKey(t *testing.T) {
	_, err := NewWorld(
		scenario.Map{
			Rooms:         map[string]map[[2]int]rune{"room": {{0, 0}: 'X'}},
			Ground:        map[string]string{"room": "floor"},
			InputProfiles: map[string]string{"room": "topdown"},
		},
		map[rune]string{},
		map[string]map[string][]string{"floor": {"pos": {}, "ascii": {"."}}},
		map[string]map[string]string{"topdown": {}},
	)
	if err == nil {
		t.Fatal("expected undefined map key error")
	}
}

func TestNewWorldBuildsTerminalRoomFromLiteralASCII(t *testing.T) {
	terminalPosition := component.Position{Room: "comms", X: 0, Y: 0}
	gameWorld, err := NewWorld(
		scenario.Map{
			Rooms: map[string]map[[2]int]rune{
				"comms": {{0, 0}: 'T'},
				"scan":  {{0, 0}: '+', {1, 0}: 'A'},
			},
			Ground:        map[string]string{"comms": "floor", "scan": "void"},
			InputProfiles: map[string]string{"comms": "topdown", "scan": "scan"},
			Terminals:     map[component.Position]string{terminalPosition: "scan"},
		},
		map[rune]string{'T': "terminal"},
		map[string]map[string][]string{
			"floor":    {"pos": {}, "ascii": {"."}},
			"void":     {"pos": {}, "ascii": {" "}},
			"terminal": {"pos": {}, "ascii": {"T"}, "commandtable": {}},
		},
		map[string]map[string]string{
			"topdown": {"interact": "e"},
			"scan":    {"exit": "e"},
		},
	)
	if err != nil {
		t.Fatalf("NewWorld returned error: %v", err)
	}

	want := map[component.Position]rune{
		{Room: "scan", X: 0, Y: 0}: '+',
		{Room: "scan", X: 1, Y: 0}: 'A',
	}
	for position, glyph := range want {
		found := false
		for _, entityID := range GetEntitiesAtPosition(gameWorld, position) {
			if gameWorld.Layer[entityID].Nr == 1 && gameWorld.Ascii[entityID].Ascii == glyph {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected literal glyph %q at %+v", glyph, position)
		}
	}
}

func TestNewWorldRejectsTerminalProfileWithoutExit(t *testing.T) {
	terminalPosition := component.Position{Room: "comms", X: 0, Y: 0}
	_, err := NewWorld(
		scenario.Map{
			Rooms: map[string]map[[2]int]rune{
				"comms": {{0, 0}: 'T'},
				"scan":  {{0, 0}: '+'},
			},
			Ground:        map[string]string{"comms": "floor", "scan": "void"},
			InputProfiles: map[string]string{"comms": "topdown", "scan": "scan"},
			Terminals:     map[component.Position]string{terminalPosition: "scan"},
		},
		map[rune]string{'T': "terminal"},
		map[string]map[string][]string{
			"floor":    {"pos": {}, "ascii": {"."}},
			"void":     {"pos": {}, "ascii": {" "}},
			"terminal": {"pos": {}, "ascii": {"T"}, "commandtable": {}},
		},
		map[string]map[string]string{
			"topdown": {"interact": "e"},
			"scan":    {"quitgame": "q"},
		},
	)
	if err == nil {
		t.Fatal("expected missing terminal exit binding error")
	}
}

func TestAddEntityStoresStructureComponents(t *testing.T) {
	gameWorld := NewWorldEmpty()
	err := gameWorld.AddEntity([2]int{}, map[string][]string{
		"helm": {}, "commandtable": {}, "bunkbed": {}, "prisonbars": {}, "wall": {}, "window": {},
	})
	if err != nil {
		t.Fatalf("AddEntity returned error: %v", err)
	}

	if !gameWorld.Helm[0].IsInteractable || !gameWorld.CommandTable[0].IsInteractable {
		t.Fatal("expected helm and command table to be interactable")
	}
	if _, ok := gameWorld.BunkBed[0]; !ok {
		t.Fatal("expected bunk bed component")
	}
	if _, ok := gameWorld.PrisonBars[0]; !ok {
		t.Fatal("expected prison bars component")
	}
	if _, ok := gameWorld.Wall[0]; !ok {
		t.Fatal("expected wall component")
	}
	if _, ok := gameWorld.Window[0]; !ok {
		t.Fatal("expected window component")
	}
}

func TestCloneCopiesMutableState(t *testing.T) {
	gameWorld := NewWorldEmpty()
	gameWorld.UserInputProfile = UserInputProfile{KeyQuitGame: "q"}
	gameWorld.KeyDown = "q"
	gameWorld.ShouldQuit = true
	gameWorld.HasChanged = true
	gameWorld.IterationNr = 4
	gameWorld.ViewRoom = "scan"
	gameWorld.InputProfiles["topdown"] = UserInputProfile{KeyMoveUp: "w"}
	gameWorld.InputProfileByRoom["bridge"] = "topdown"
	portalFrom := component.Position{Room: "bridge", X: 1, Y: 2}
	portalTo := component.Position{Room: "comms", X: 3, Y: 4}
	gameWorld.Portals[portalFrom] = portalTo
	gameWorld.Terminals[portalFrom] = "scan"
	if err := gameWorld.AddEntity([2]int{4, 5}, map[string][]string{
		"pos": {}, "ascii": {"o"}, "impassable": {}, "player": {}, "door": {},
	}); err != nil {
		t.Fatalf("AddEntity returned error: %v", err)
	}

	clone := gameWorld.Clone()
	if clone.UserInputProfile != gameWorld.UserInputProfile || clone.KeyDown != "q" || !clone.ShouldQuit {
		t.Fatalf("expected scalar state to be copied, got %+v", clone)
	}
	if !clone.HasChanged || clone.IterationNr != 4 {
		t.Fatalf("expected update state to be copied, got %+v", clone)
	}
	if clone.ViewRoom != "scan" {
		t.Fatalf("expected view room to be copied, got %q", clone.ViewRoom)
	}

	clone.Entities[0] = 9
	clone.Pos[0] = component.Position{X: 1, Y: 1}
	clone.Layer[0] = component.Layer{Nr: 2}
	clone.EByPos[component.Position{X: 1, Y: 1}] = 0
	delete(clone.InputProfiles, "topdown")
	delete(clone.InputProfileByRoom, "bridge")
	delete(clone.Portals, portalFrom)
	delete(clone.Terminals, portalFrom)
	clone.Ascii[0] = component.Ascii{Ascii: 'x'}
	delete(clone.Impassable, 0)
	delete(clone.Player, 0)
	clone.Door[0] = component.Door{IsInteractable: false}

	if gameWorld.Entities[0] != 0 {
		t.Fatal("expected entities slice to be independent")
	}
	if gameWorld.Pos[0] != (component.Position{X: 4, Y: 5}) {
		t.Fatal("expected position map to be independent")
	}
	if gameWorld.Layer[0].Nr != 0 {
		t.Fatal("expected layer map to be independent")
	}
	if _, ok := gameWorld.EByPos[component.Position{X: 1, Y: 1}]; ok {
		t.Fatal("expected reverse position map to be independent")
	}
	if gameWorld.Portals[portalFrom] != portalTo {
		t.Fatal("expected portal map to be independent")
	}
	if gameWorld.Terminals[portalFrom] != "scan" {
		t.Fatal("expected terminal map to be independent")
	}
	if gameWorld.InputProfiles["topdown"].KeyMoveUp != "w" || gameWorld.InputProfileByRoom["bridge"] != "topdown" {
		t.Fatal("expected input profile maps to be independent")
	}
	if gameWorld.Ascii[0].Ascii != 'o' {
		t.Fatal("expected glyph map to be independent")
	}
	if _, ok := gameWorld.Impassable[0]; !ok {
		t.Fatal("expected impassable map to be independent")
	}
	if _, ok := gameWorld.Player[0]; !ok {
		t.Fatal("expected player map to be independent")
	}
	if !gameWorld.Door[0].IsInteractable {
		t.Fatal("expected door map to be independent")
	}
}
