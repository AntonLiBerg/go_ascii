package tests

import (
	component "go_ascii/internal"
	"go_ascii/internal/scenario"
	"go_ascii/internal/world"
	"testing"
)

func TestAddEntityStoresComponents(t *testing.T) {
	gameWorld := world.NewWorldEmpty()

	err := gameWorld.AddEntity([2]int{2, 3}, map[string][]string{
		"pos":          {},
		"ascii":        {"å"},
		"impassable":   {},
		"interactable": {component.InteractionTypeDoor},
		"player":       {},
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
	if got := gameWorld.Interactable[0].InteractionType; got != component.InteractionTypeDoor {
		t.Fatalf("expected door interaction, got %q", got)
	}
}

func TestAddEntityRejectsUnknownComponent(t *testing.T) {
	gameWorld := world.NewWorldEmpty()

	if err := gameWorld.AddEntity([2]int{}, map[string][]string{"visible": {}}); err == nil {
		t.Fatal("expected unknown component error")
	}
}

func TestAddEntityDoesNotPartiallyMutateOnError(t *testing.T) {
	gameWorld := world.NewWorldEmpty()
	err := gameWorld.AddEntity([2]int{2, 3}, map[string][]string{
		"pos":     {},
		"ascii":   {"."},
		"unknown": {},
	})
	if err == nil {
		t.Fatal("expected unknown component error")
	}
	if len(gameWorld.Entities) != 0 || len(gameWorld.Pos) != 0 || len(gameWorld.EByPos) != 0 {
		t.Fatalf("expected failed entity insertion to leave no state, got entities=%v pos=%v index=%v", gameWorld.Entities, gameWorld.Pos, gameWorld.EByPos)
	}
}

func TestAddEntityRequiresOneInteractionType(t *testing.T) {
	for _, values := range [][]string{nil, {"door", "terminal"}} {
		gameWorld := world.NewWorldEmpty()
		if err := gameWorld.AddEntity([2]int{}, map[string][]string{"interactable": values}); err == nil {
			t.Fatalf("expected interaction values %v to be rejected", values)
		}
	}
}

func TestNewWorldBuildsGroundAndRoomEntities(t *testing.T) {
	gameWorld, err := world.NewWorld(
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
	if gameWorld.Room != "bridge" {
		t.Fatalf("expected active room bridge, got %q", gameWorld.Room)
	}
}

func TestNewWorldRejectsEmptyRooms(t *testing.T) {
	if _, err := world.NewWorld(scenario.Map{}, nil, nil, nil); err == nil {
		t.Fatal("expected empty room error")
	}
}

func TestNewWorldRequiresExactlyOnePlayer(t *testing.T) {
	baseMap := scenario.Map{
		Rooms:         map[string]map[[2]int]rune{"room": {{0, 0}: '.'}},
		Ground:        map[string]string{"room": "floor"},
		InputProfiles: map[string]string{"room": "topdown"},
	}
	components := map[string]map[string][]string{
		"floor":  {"pos": {}, "ascii": {"."}},
		"player": {"pos": {}, "ascii": {"o"}, "player": {}},
	}
	if _, err := world.NewWorld(baseMap, nil, components, map[string]map[string]string{"topdown": {}}); err == nil {
		t.Fatal("expected missing player error")
	}
	twoPlayerMap := baseMap
	twoPlayerMap.Rooms = map[string]map[[2]int]rune{"room": {{0, 0}: 'a', {1, 0}: 'b'}}
	if _, err := world.NewWorld(twoPlayerMap, map[rune]string{'a': "player", 'b': "player"}, components, map[string]map[string]string{"topdown": {}}); err == nil {
		t.Fatal("expected multiple player error")
	}
}

func TestNewWorldRejectsReferencedEntityWithoutDefinition(t *testing.T) {
	_, err := world.NewWorld(
		scenario.Map{
			Rooms:         map[string]map[[2]int]rune{"room": {{0, 0}: 'x'}},
			Ground:        map[string]string{"room": "floor"},
			InputProfiles: map[string]string{"room": "topdown"},
		},
		map[rune]string{'x': "missing"},
		map[string]map[string][]string{"floor": {"pos": {}, "ascii": {"."}}},
		map[string]map[string]string{"topdown": {}},
	)
	if err == nil {
		t.Fatal("expected missing entity definition error")
	}
}

func TestNewWorldSelectsPlayerRoom(t *testing.T) {
	gameWorld, err := world.NewWorld(
		scenario.Map{
			Rooms:         map[string]map[[2]int]rune{"aaa": {{0, 0}: '.'}, "zzz": {{0, 0}: 'o'}},
			Ground:        map[string]string{"aaa": "floor", "zzz": "floor"},
			InputProfiles: map[string]string{"aaa": "topdown", "zzz": "topdown"},
		},
		map[rune]string{'.': "floor", 'o': "player"},
		map[string]map[string][]string{
			"floor":  {"pos": {}, "ascii": {"."}},
			"player": {"pos": {}, "ascii": {"o"}, "player": {}},
		},
		map[string]map[string]string{"topdown": {}},
	)
	if err != nil {
		t.Fatalf("NewWorld returned error: %v", err)
	}
	if gameWorld.Room != "zzz" {
		t.Fatalf("expected player room zzz, got %q", gameWorld.Room)
	}
}

func TestNewWorldRejectsUndefinedMapKey(t *testing.T) {
	_, err := world.NewWorld(
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
	gameWorld, err := world.NewWorld(
		scenario.Map{
			Rooms: map[string]map[[2]int]rune{
				"comms": {{0, 0}: 'T', {1, 0}: '@'},
				"scan":  {{0, 0}: '+', {1, 0}: 'A'},
			},
			Ground:        map[string]string{"comms": "floor", "scan": "void"},
			InputProfiles: map[string]string{"comms": "topdown", "scan": "scan"},
			Terminals:     map[component.Position]string{terminalPosition: "scan"},
		},
		map[rune]string{'T': "terminal", '@': "player"},
		map[string]map[string][]string{
			"floor":    {"pos": {}, "ascii": {"."}},
			"void":     {"pos": {}, "ascii": {" "}},
			"terminal": {"pos": {}, "ascii": {"T"}, "interactable": {component.InteractionTypeTerminal}},
			"player":   {"pos": {}, "ascii": {"o"}, "player": {}},
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
		for _, entityID := range world.GetEntitiesAtPosition(gameWorld, position) {
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
	_, err := world.NewWorld(
		scenario.Map{
			Rooms: map[string]map[[2]int]rune{
				"comms": {{0, 0}: 'T', {1, 0}: '@'},
				"scan":  {{0, 0}: '+'},
			},
			Ground:        map[string]string{"comms": "floor", "scan": "void"},
			InputProfiles: map[string]string{"comms": "topdown", "scan": "scan"},
			Terminals:     map[component.Position]string{terminalPosition: "scan"},
		},
		map[rune]string{'T': "terminal", '@': "player"},
		map[string]map[string][]string{
			"floor":    {"pos": {}, "ascii": {"."}},
			"void":     {"pos": {}, "ascii": {" "}},
			"terminal": {"pos": {}, "ascii": {"T"}, "interactable": {component.InteractionTypeTerminal}},
			"player":   {"pos": {}, "ascii": {"o"}, "player": {}},
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

func TestCloneCopiesMutableState(t *testing.T) {
	gameWorld := world.NewWorldEmpty()
	gameWorld.UserInputProfile = world.UserInputProfile{KeyQuitGame: "q"}
	gameWorld.KeyDown = "q"
	gameWorld.ShouldQuit = true
	gameWorld.HasChanged = true
	gameWorld.IterationNr = 4
	gameWorld.Room = "scan"
	gameWorld.InputProfiles["topdown"] = world.UserInputProfile{KeyMoveUp: "w"}
	gameWorld.InputProfileByRoom["bridge"] = "topdown"
	portalFrom := component.Position{Room: "bridge", X: 1, Y: 2}
	portalTo := component.Position{Room: "comms", X: 3, Y: 4}
	gameWorld.Portals[portalFrom] = portalTo
	gameWorld.Terminals[portalFrom] = "scan"
	if err := gameWorld.AddEntity([2]int{4, 5}, map[string][]string{
		"pos": {}, "ascii": {"o"}, "impassable": {}, "player": {}, "interactable": {component.InteractionTypeDoor},
	}); err != nil {
		t.Fatalf("AddEntity returned error: %v", err)
	}

	clone := gameWorld.Clone()
	if clone.UserInputProfile != gameWorld.UserInputProfile || clone.KeyDown != "q" || !clone.ShouldQuit {
		t.Fatalf("expected scalar state to be copied, got %+v", clone)
	}
	if !clone.HasChanged || clone.IterationNr != 4 || clone.Room != "scan" {
		t.Fatalf("expected update state to be copied, got %+v", clone)
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
	clone.Interactable[0] = component.Interactable{InteractionType: component.InteractionTypeTerminal}

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
	if gameWorld.Interactable[0].InteractionType != component.InteractionTypeDoor {
		t.Fatal("expected interactable map to be independent")
	}
}
