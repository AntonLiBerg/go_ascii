package world

import "testing"

func TestAddEntityStoresComponents(t *testing.T) {
	gameWorld := NewWorldEmpty()

	err := gameWorld.AddEntity([2]int{2, 3}, map[string][]string{
		"pos":        {},
		"ascii":      {"å"},
		"impassable": {},
		"machine":    {},
		"player":     {},
	})
	if err != nil {
		t.Fatalf("AddEntity returned error: %v", err)
	}

	if len(gameWorld.Entities) != 1 || gameWorld.Entities[0] != 0 {
		t.Fatalf("expected entity 0, got %v", gameWorld.Entities)
	}
	if got := gameWorld.Pos[0]; got != (Position{X: 2, Y: 3}) {
		t.Fatalf("expected position 2,3, got %+v", got)
	}
	if got := gameWorld.EByPos[Position{X: 2, Y: 3}]; got != 0 {
		t.Fatalf("expected reverse position index to point to entity 0, got %d", got)
	}
	if got := gameWorld.Ascii[0]; got != 'å' {
		t.Fatalf("expected glyph 'å', got %q", got)
	}
	if _, ok := gameWorld.Impassable[0]; !ok {
		t.Fatal("expected impassable component")
	}
	if _, ok := gameWorld.Player[0]; !ok {
		t.Fatal("expected player component")
	}
	if isOn, ok := gameWorld.Machine[0]; !ok || isOn {
		t.Fatalf("expected an off machine, got isOn=%t, exists=%t", isOn, ok)
	}
}

func TestAddEntityRejectsUnknownComponent(t *testing.T) {
	gameWorld := NewWorldEmpty()

	if err := gameWorld.AddEntity([2]int{}, map[string][]string{"visible": {}}); err == nil {
		t.Fatal("expected unknown component error")
	}
}

func TestCloneCopiesMutableState(t *testing.T) {
	gameWorld := NewWorldEmpty()
	gameWorld.UserInputProfile = UserInputProfile{KeyQuitGame: "q"}
	gameWorld.KeyDown = "q"
	gameWorld.ShouldQuit = true
	gameWorld.HasChanged = true
	gameWorld.IterationNr = 4
	if err := gameWorld.AddEntity([2]int{4, 5}, map[string][]string{
		"pos": {}, "ascii": {"o"}, "impassable": {}, "player": {}, "machine": {},
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

	clone.Entities[0] = 9
	clone.Pos[0] = Position{X: 1, Y: 1}
	clone.EByPos[Position{X: 1, Y: 1}] = 0
	clone.Ascii[0] = 'x'
	delete(clone.Impassable, 0)
	delete(clone.Player, 0)
	clone.Machine[0] = true

	if gameWorld.Entities[0] != 0 {
		t.Fatal("expected entities slice to be independent")
	}
	if gameWorld.Pos[0] != (Position{X: 4, Y: 5}) {
		t.Fatal("expected position map to be independent")
	}
	if _, ok := gameWorld.EByPos[Position{X: 1, Y: 1}]; ok {
		t.Fatal("expected reverse position map to be independent")
	}
	if gameWorld.Ascii[0] != 'o' {
		t.Fatal("expected glyph map to be independent")
	}
	if _, ok := gameWorld.Impassable[0]; !ok {
		t.Fatal("expected impassable map to be independent")
	}
	if _, ok := gameWorld.Player[0]; !ok {
		t.Fatal("expected player map to be independent")
	}
	if gameWorld.Machine[0] {
		t.Fatal("expected machine map to be independent")
	}
}
