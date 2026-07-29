package world

import (
	component "go_ascii/internal"
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

func TestCloneCopiesMutableState(t *testing.T) {
	gameWorld := NewWorldEmpty()
	gameWorld.UserInputProfile = UserInputProfile{KeyQuitGame: "q"}
	gameWorld.KeyDown = "q"
	gameWorld.ShouldQuit = true
	gameWorld.HasChanged = true
	gameWorld.IterationNr = 4
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

	clone.Entities[0] = 9
	clone.Pos[0] = component.Position{X: 1, Y: 1}
	clone.EByPos[component.Position{X: 1, Y: 1}] = 0
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
	if _, ok := gameWorld.EByPos[component.Position{X: 1, Y: 1}]; ok {
		t.Fatal("expected reverse position map to be independent")
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
