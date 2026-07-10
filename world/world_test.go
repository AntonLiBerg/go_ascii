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
		cmp.C_PLAYER:     {},
		cmp.C_VISIBLE:    {},
	})
	if err != nil {
		t.Fatalf("AddEntity returned error: %v", err)
	}

	if !world.HasComponent(0, cmp.C_PLAYER) {
		t.Fatal("expected entity 0 to have player component")
	}
	if !world.HasComponent(0, cmp.C_VISIBLE) {
		t.Fatal("expected entity 0 to have visible component")
	}
	if !world.HasComponent(0, cmp.C_POS) {
		t.Fatal("expected entity 0 to have position component")
	}
	if !world.HasComponent(0, cmp.C_ASCII) {
		t.Fatal("expected entity 0 to have ascii component")
	}
	if !world.HasComponent(0, cmp.C_IMPASSABLE) {
		t.Fatal("expected entity 0 to have impassable component")
	}
	if !world.HasComponent(0, cmp.C_MACHINE) {
		t.Fatal("expected entity 0 to have machine component")
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
		AddMenuChoices(MenuChoices{
			ShouldShow: true,
			Header:     "test menu",
			Choices:    []MenuChoice{{Text: "one"}},
		}).
		AddPosition(eID, cmp.Position{X: 2, Y: 3}).
		AddAscii(eID, cmp.Ascii{Ascii: 'o'}).
		AddImpassable(eID).
		AddMachine(eID, cmp.Machine{MachineType: cmp.MACHINENAME_RADIO}).
		AddPlayer(eID).
		AddVisible(eID)

	if len(world.Entities) != 0 {
		t.Fatal("expected original world to keep no entities")
	}
	if _, ok := world.UserInput["q"]; ok {
		t.Fatal("expected original world to keep no user input")
	}
	if _, ok := world.Pos[eID]; ok {
		t.Fatal("expected original world to keep no position")
	}
	if world.HasComponent(eID, cmp.C_PLAYER) {
		t.Fatal("expected original world to keep no player component")
	}
	if world.MenuChoices.ShouldShow {
		t.Fatal("expected original world to keep empty menu choices")
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
	if !updated.MenuChoices.ShouldShow || updated.MenuChoices.Header != "test menu" || len(updated.MenuChoices.Choices) != 1 {
		t.Fatalf("expected updated world to store menu choices, got %+v", updated.MenuChoices)
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
	if !updated.HasComponent(eID, cmp.C_PLAYER) {
		t.Fatal("expected updated world to store player component")
	}
	if !updated.HasComponent(eID, cmp.C_VISIBLE) {
		t.Fatal("expected updated world to store visible component")
	}
	if !updated.HasComponent(eID, cmp.C_POS) {
		t.Fatal("expected updated world to store position component")
	}
	if !updated.HasComponent(eID, cmp.C_ASCII) {
		t.Fatal("expected updated world to store ascii component")
	}
	if !updated.HasComponent(eID, cmp.C_IMPASSABLE) {
		t.Fatal("expected updated world to store impassable component")
	}
	if !updated.HasComponent(eID, cmp.C_MACHINE) {
		t.Fatal("expected updated world to store machine component")
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

func TestAddPlayerClonesAndKeepsOriginalWorld(t *testing.T) {
	world := NewWorldEmpty()
	world, eID := world.AddNewEntity()

	updated := world.AddPlayer(eID)

	if !updated.HasComponent(eID, cmp.C_PLAYER) {
		t.Fatal("expected updated world to store player component")
	}
	if world.HasComponent(eID, cmp.C_PLAYER) {
		t.Fatal("expected original world to stay unchanged")
	}
}

func TestAddMenuChoicesClonesAndKeepsOriginalWorld(t *testing.T) {
	world := NewWorldEmpty()

	updated := world.AddMenuChoices(MenuChoices{
		ShouldShow: true,
		Header:     "choices",
		Choices:    []MenuChoice{{Text: "pick me"}},
	})

	if !updated.MenuChoices.ShouldShow || updated.MenuChoices.Header != "choices" || len(updated.MenuChoices.Choices) != 1 {
		t.Fatalf("expected updated world to store menu choices, got %+v", updated.MenuChoices)
	}
	if world.MenuChoices.ShouldShow || world.MenuChoices.Header != "" || len(world.MenuChoices.Choices) != 0 {
		t.Fatalf("expected original world to stay unchanged, got %+v", world.MenuChoices)
	}

	updated.MenuChoices.Choices[0].Text = "changed"
	if len(world.MenuChoices.Choices) != 0 {
		t.Fatal("expected original world menu choices slice to stay independent")
	}
}

func TestCloneCopiesComponents(t *testing.T) {
	world := NewWorldEmpty()
	world.UserInputProfile = usr.UserInputProfile{KeyQuitGame: "q", KeyMoveDown: "s"}
	world.SetUserState(usr.S_quit, true)
	world.UserInput["q"] = true
	world.MenuChoices = MenuChoices{
		ShouldShow: true,
		IsOpen:     true,
		Header:     "cloned menu",
		Choices:    []MenuChoice{{Text: "first"}},
	}
	world.Entities = []int{1}
	world.NextEnt = 2
	world.Pos[1] = cmp.Position{X: 4, Y: 5}
	world.Ascii[1] = cmp.Ascii{Ascii: 'o'}
	world.Impassable[1] = cmp.Impassable{}
	world.Player[1] = cmp.Player{}
	world.Machine[1] = cmp.Machine{MachineType: cmp.MACHINENAME_RADIO}
	world.EByPos[cmp.Position{X: 4, Y: 5}] = 1

	clone := world.Clone()

	if clone.UserInputProfile.KeyQuitGame != "q" || clone.UserInputProfile.KeyMoveDown != "s" {
		t.Fatalf("expected user input profile to be copied, got %+v", clone.UserInputProfile)
	}
	if !clone.HasUserState(usr.S_quit) {
		t.Fatalf("expected state user to be copied, got %v", clone.StateUser)
	}
	if !clone.MenuChoices.ShouldShow || !clone.MenuChoices.IsOpen || clone.MenuChoices.Header != "cloned menu" || len(clone.MenuChoices.Choices) != 1 {
		t.Fatalf("expected menu choices to be copied, got %+v", clone.MenuChoices)
	}
	if !clone.HasComponent(1, cmp.C_PLAYER) {
		t.Fatal("expected cloned world to keep player component")
	}
	if !clone.HasComponent(1, cmp.C_POS) {
		t.Fatal("expected cloned world to keep position component")
	}
	if !clone.HasComponent(1, cmp.C_ASCII) {
		t.Fatal("expected cloned world to keep ascii component")
	}
	if !clone.HasComponent(1, cmp.C_IMPASSABLE) {
		t.Fatal("expected cloned world to keep impassable component")
	}
	if !clone.HasComponent(1, cmp.C_MACHINE) {
		t.Fatal("expected cloned world to keep machine component")
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

	clone.Visible[1] = cmp.Visible{}
	if _, ok := world.Visible[1]; ok {
		t.Fatal("expected cloned visible map to be independent")
	}
	clone.MenuChoices.Choices[0].Text = "changed"
	if world.MenuChoices.Choices[0].Text == "changed" {
		t.Fatal("expected cloned menu choices to be independent")
	}
	clone.SetUserState(usr.S_INTERACT, true)
	if world.HasUserState(usr.S_INTERACT) {
		t.Fatal("expected cloned state user map to be independent")
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
