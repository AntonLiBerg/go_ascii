package world

import (
	cmp "go_ascii/component"
	usr "go_ascii/user"
	"testing"
)

func TestAddEntityStoresComponents(t *testing.T) {
	world := NewWorldEmpty()

	err := world.AddEntity([2]int{2, 3}, map[cmp.ComponentName][]string{
		cmp.C_POS:        {},
		cmp.C_ASCII:      {"o"},
		cmp.C_IMPASSABLE: {},
		cmp.C_TAGS:       {string(cmp.TAG_PLAYER), "visible"},
	})
	if err != nil {
		t.Fatalf("AddEntity returned error: %v", err)
	}

	tags, ok := world.Tags[0]
	if !ok {
		t.Fatal("expected tags for entity 0")
	}
	if !tags.Vals[cmp.TAG_PLAYER] {
		t.Fatal("expected entity 0 to have player tag")
	}
	if !tags.Vals[cmp.Tag("visible")] {
		t.Fatal("expected entity 0 to have visible tag")
	}
	if !world.EByTag[cmp.TAG_PLAYER][0] {
		t.Fatal("expected reverse tag index to include entity 0 for player")
	}
	if !world.EByTag[cmp.Tag("visible")][0] {
		t.Fatal("expected reverse tag index to include entity 0 for visible")
	}
	gotPosID, ok := world.EByPos[cmp.Position{X: 2, Y: 3}]
	if !ok || gotPosID != 0 {
		t.Fatalf("expected reverse position index to map 2,3 to entity 0, got %d, exists=%t", gotPosID, ok)
	}
	if _, ok := world.Impassable[0]; !ok {
		t.Fatal("expected entity 0 to have impassable component")
	}
}

func TestCloneCopiesComponents(t *testing.T) {
	world := NewWorldEmpty()
	world.UserInputProfile = usr.UserInputProfile{KeyQuitGame: "q", KeyMoveDown: "s"}
	world.StateUser = usr.S_quit
	world.UserInput["q"] = true
	world.Entities = []int{1}
	world.NextEnt = 2
	world.Pos[1] = cmp.Position{X: 4, Y: 5}
	world.Ascii[1] = cmp.Ascii{Ascii: 'o'}
	world.Impassable[1] = cmp.Impassable{}
	world.Tags[1] = cmp.Tags{Vals: map[cmp.Tag]bool{cmp.TAG_PLAYER: true}}
	world.EByTag[cmp.TAG_PLAYER] = map[int]bool{1: true}
	world.EByPos[cmp.Position{X: 4, Y: 5}] = 1

	clone := world.Clone()

	if clone.UserInputProfile.KeyQuitGame != "q" || clone.UserInputProfile.KeyMoveDown != "s" {
		t.Fatalf("expected user input profile to be copied, got %+v", clone.UserInputProfile)
	}
	if clone.StateUser != usr.S_quit {
		t.Fatalf("expected state user to be copied, got %v", clone.StateUser)
	}
	if !clone.Tags[1].Vals[cmp.TAG_PLAYER] {
		t.Fatal("expected cloned world to keep player tag")
	}
	if !clone.EByTag[cmp.TAG_PLAYER][1] {
		t.Fatal("expected cloned world to keep reverse player tag index")
	}
	gotPosID, ok := clone.EByPos[cmp.Position{X: 4, Y: 5}]
	if !ok || gotPosID != 1 {
		t.Fatalf("expected cloned world to keep reverse position index, got %d, exists=%t", gotPosID, ok)
	}
	if _, ok := clone.Impassable[1]; !ok {
		t.Fatal("expected cloned world to keep impassable component")
	}

	clone.Tags[1].Vals[cmp.Tag("new")] = true
	if world.Tags[1].Vals[cmp.Tag("new")] {
		t.Fatal("expected cloned tags map to be independent")
	}
	clone.EByTag[cmp.TAG_PLAYER][2] = true
	if world.EByTag[cmp.TAG_PLAYER][2] {
		t.Fatal("expected cloned reverse tag index to be independent")
	}
	clone.EByPos[cmp.Position{X: 1, Y: 1}] = 9
	if _, ok := world.EByPos[cmp.Position{X: 1, Y: 1}]; ok {
		t.Fatal("expected cloned reverse position index to be independent")
	}
	delete(clone.Impassable, 1)
	if _, ok := world.Impassable[1]; !ok {
		t.Fatal("expected cloned impassable map to be independent")
	}
}
