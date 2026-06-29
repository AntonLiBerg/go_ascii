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
		cmp.C_MACHINE:    {string(cmp.MACHINENAME_RADIO)},
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
	if got := world.Machine[0].MachineType; got != cmp.MACHINENAME_RADIO {
		t.Fatalf("expected entity 0 to have radio machine component, got %s", got)
	}
}

func TestWorldAddMethodsCloneAndReturnWorld(t *testing.T) {
	world := NewWorldEmpty()

	updated, eID := world.AddNewEntity()
	updated = updated.
		AddUserInput("q", true).
		AddPosition(eID, cmp.Position{X: 2, Y: 3}).
		AddAscii(eID, cmp.Ascii{Ascii: 'o'}).
		AddImpassable(eID).
		AddMachine(eID, cmp.Machine{MachineType: cmp.MACHINENAME_RADIO}).
		AddTag(eID, cmp.TAG_PLAYER).
		AddTag(eID, cmp.Tag("visible"))

	if len(world.Entities) != 0 {
		t.Fatal("expected original world to keep no entities")
	}
	if _, ok := world.UserInput["q"]; ok {
		t.Fatal("expected original world to keep no user input")
	}
	if _, ok := world.Pos[eID]; ok {
		t.Fatal("expected original world to keep no position")
	}

	if eID != 0 {
		t.Fatalf("expected first entity id to be 0, got %d", eID)
	}
	if updated.NextEnt != 1 {
		t.Fatalf("expected next entity id to be 1, got %d", updated.NextEnt)
	}
	if !updated.UserInput["q"] {
		t.Fatal("expected updated world to store user input")
	}
	if got := updated.Pos[eID]; got != (cmp.Position{X: 2, Y: 3}) {
		t.Fatalf("expected updated position 2,3, got %+v", got)
	}
	if got := updated.EByPos[cmp.Position{X: 2, Y: 3}]; got != eID {
		t.Fatalf("expected reverse position index to point at entity %d, got %d", eID, got)
	}
	if got := updated.Ascii[eID].Ascii; got != 'o' {
		t.Fatalf("expected ascii o, got %q", got)
	}
	if _, ok := updated.Impassable[eID]; !ok {
		t.Fatal("expected updated world to store impassable component")
	}
	if got := updated.Machine[eID].MachineType; got != cmp.MACHINENAME_RADIO {
		t.Fatalf("expected updated world to store radio machine component, got %s", got)
	}
	if !updated.Tags[eID].Vals[cmp.TAG_PLAYER] {
		t.Fatal("expected updated world to store player tag")
	}
	if !updated.Tags[eID].Vals[cmp.Tag("visible")] {
		t.Fatal("expected updated world to store visible tag")
	}
	if !updated.EByTag[cmp.TAG_PLAYER][eID] {
		t.Fatal("expected reverse tag index to include player tag")
	}
	if !updated.EByTag[cmp.Tag("visible")][eID] {
		t.Fatal("expected reverse tag index to include visible tag")
	}
}

func TestAddPositionRemovesOldReverseIndex(t *testing.T) {
	world := NewWorldEmpty()
	world, eID := world.AddNewEntity()
	world = world.AddPosition(eID, cmp.Position{X: 1, Y: 1})

	updated := world.AddPosition(eID, cmp.Position{X: 2, Y: 2})

	if _, ok := updated.EByPos[cmp.Position{X: 1, Y: 1}]; ok {
		t.Fatal("expected old reverse position index to be removed")
	}
	if got := updated.EByPos[cmp.Position{X: 2, Y: 2}]; got != eID {
		t.Fatalf("expected new reverse position index to point at entity %d, got %d", eID, got)
	}
	if got := world.EByPos[cmp.Position{X: 1, Y: 1}]; got != eID {
		t.Fatalf("expected original reverse position index to stay unchanged, got %d", got)
	}
}

func TestAddTagsReplacesReverseTagIndex(t *testing.T) {
	world := NewWorldEmpty()
	world, eID := world.AddNewEntity()
	world = world.AddTags(eID, cmp.Tags{Vals: map[cmp.Tag]bool{cmp.TAG_PLAYER: true}})

	updated := world.AddTags(eID, cmp.Tags{Vals: map[cmp.Tag]bool{cmp.Tag("visible"): true}})

	if _, ok := updated.EByTag[cmp.TAG_PLAYER][eID]; ok {
		t.Fatal("expected old reverse tag index to be removed")
	}
	if !updated.EByTag[cmp.Tag("visible")][eID] {
		t.Fatal("expected new reverse tag index to include entity")
	}
	if !world.EByTag[cmp.TAG_PLAYER][eID] {
		t.Fatal("expected original reverse tag index to stay unchanged")
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
	world.Machine[1] = cmp.Machine{MachineType: cmp.MACHINENAME_RADIO}
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
	if got := clone.Machine[1].MachineType; got != cmp.MACHINENAME_RADIO {
		t.Fatalf("expected cloned world to keep radio machine component, got %s", got)
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
	delete(clone.Machine, 1)
	if _, ok := world.Machine[1]; !ok {
		t.Fatal("expected cloned machine map to be independent")
	}
}
