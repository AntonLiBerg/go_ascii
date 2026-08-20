package control_test

import (
	component "go_ascii/internal"
	"go_ascii/internal/game"
	"go_ascii/internal/service/control"
	"go_ascii/internal/world"
	"testing"
)

func TestControl(t *testing.T) {
	tests := []struct {
		Name      string
		MakeWorld func() world.World
		Assert    func(world.World, game.UpdateFunc, *testing.T)
	}{
		{
			Name: "not control input profile",
			MakeWorld: func() world.World {
				return world.NewWorldEmpty()
			},
			Assert: func(_ world.World, update game.UpdateFunc, t *testing.T) {
				assertEmptyUpdateFunc(t, update)
			},
		},
		{
			Name: "no key down",
			MakeWorld: func() world.World {
				return newControlTestWorld()
			},
			Assert: func(_ world.World, update game.UpdateFunc, t *testing.T) {
				assertEmptyUpdateFunc(t, update)
			},
		},
		{
			Name: "not editing, next key down",
			MakeWorld: func() world.World {
				gameWorld := newControlTestWorld()
				gameWorld.KeyDown = "n"
				return gameWorld
			},
			Assert: func(gameWorld world.World, update game.UpdateFunc, t *testing.T) {
				assertNotEmptyUpdateFunc(t, update)
				originalActive := gameWorld.ActiveControl
				next := applyUpdate(t, update, gameWorld)

				if gameWorld.ActiveControl != originalActive {
					t.Fatal("expected update not to mutate input focus")
				}
				if next.ActiveControl.SelectableEntityID != 2 {
					t.Fatalf("expected control 2 to be focused, got %d", next.ActiveControl.SelectableEntityID)
				}
				if next.EditingControl {
					t.Fatal("expected navigation not to enter editing mode")
				}
				if next.KeyDown != "" || !next.HasChanged {
					t.Fatal("expected navigation key to be consumed and world marked changed")
				}
			},
		},
		{
			Name: "not editing, previous key down",
			MakeWorld: func() world.World {
				gameWorld := newControlTestWorld()
				gameWorld.KeyDown = "p"
				return gameWorld
			},
			Assert: func(gameWorld world.World, update game.UpdateFunc, t *testing.T) {
				assertNotEmptyUpdateFunc(t, update)
				originalActive := gameWorld.ActiveControl
				next := applyUpdate(t, update, gameWorld)

				if gameWorld.ActiveControl != originalActive {
					t.Fatal("expected update not to mutate input focus")
				}
				if next.ActiveControl.SelectableEntityID != 3 {
					t.Fatalf("expected control 3 to be focused, got %d", next.ActiveControl.SelectableEntityID)
				}
				if next.EditingControl {
					t.Fatal("expected navigation not to enter editing mode")
				}
				if next.KeyDown != "" || !next.HasChanged {
					t.Fatal("expected navigation key to be consumed and world marked changed")
				}
			},
		},
		{
			Name: "not editing, select key down",
			MakeWorld: func() world.World {
				gameWorld := newControlTestWorld()
				gameWorld.KeyDown = "s"
				return gameWorld
			},
			Assert: func(gameWorld world.World, update game.UpdateFunc, t *testing.T) {
				assertNotEmptyUpdateFunc(t, update)
				originalActive := gameWorld.ActiveControl
				next := applyUpdate(t, update, gameWorld)

				if gameWorld.EditingControl {
					t.Fatal("expected update not to mutate input editing state")
				}
				if !next.EditingControl {
					t.Fatal("expected control to enter editing mode")
				}
				if next.ActiveControl != originalActive {
					t.Fatal("expected active control to remain focused")
				}
				if next.KeyDown != "" || !next.HasChanged {
					t.Fatal("expected select key to be consumed and world marked changed")
				}
			},
		},
		{
			Name: "editing, next key down",
			MakeWorld: func() world.World {
				gameWorld := newControlTestWorld()
				gameWorld.EditingControl = true
				gameWorld.KeyDown = "n"
				return gameWorld
			},
			Assert: func(gameWorld world.World, update game.UpdateFunc, t *testing.T) {
				assertNotEmptyUpdateFunc(t, update)
				next := applyUpdate(t, update, gameWorld)

				if gameWorld.ControlNumber[4].ValueCurrent != 1 {
					t.Fatal("expected update not to mutate input control value")
				}
				if next.ActiveControl != gameWorld.ActiveControl || !next.EditingControl {
					t.Fatal("expected value update to preserve active editing state")
				}
				if next.ControlNumber[4].ValueCurrent != 2 {
					t.Fatalf("expected control value 2, got %d", next.ControlNumber[4].ValueCurrent)
				}
				if next.KeyDown != "" || !next.HasChanged {
					t.Fatal("expected value key to be consumed and world marked changed")
				}
			},
		},
		{
			Name: "editing, select key down",
			MakeWorld: func() world.World {
				gameWorld := newControlTestWorld()
				gameWorld.EditingControl = true
				gameWorld.KeyDown = "s"
				return gameWorld
			},
			Assert: func(gameWorld world.World, update game.UpdateFunc, t *testing.T) {
				assertNotEmptyUpdateFunc(t, update)
				next := applyUpdate(t, update, gameWorld)

				if next.EditingControl {
					t.Fatal("expected control to exit editing mode")
				}
				if next.ActiveControl != gameWorld.ActiveControl {
					t.Fatal("expected active control to remain focused")
				}
			},
		},
		{
			Name: "room input profile resets control state",
			MakeWorld: func() world.World {
				gameWorld := newControlTestWorld()
				gameWorld.InputProfiles["control"] = gameWorld.UserInputProfile
				gameWorld.InputProfileByRoom["testroom"] = "control"
				gameWorld.ActiveControl = gameWorld.ActiveControl.Next
				gameWorld.EditingControl = true
				return gameWorld
			},
			Assert: func(gameWorld world.World, _ game.UpdateFunc, t *testing.T) {
				var profileSet bool
				gameWorld, profileSet = gameWorld.WithInputProfileForRoom("testroom")
				if !profileSet {
					t.Fatal("expected room input profile to be set")
				}
				if gameWorld.ActiveControl != gameWorld.ControlOrder["testroom"] {
					t.Fatal("expected focus to reset to the room's first control")
				}
				if gameWorld.EditingControl {
					t.Fatal("expected editing mode to be cleared")
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			gameWorld := test.MakeWorld()
			update := control.ServiceControl{}.GetUpdateFunc(gameWorld)
			test.Assert(gameWorld, update, t)
		})
	}
}

func TestControlLabelAppendsAndRemoves(t *testing.T) {
	gameWorld := world.NewWorldEmpty()
	gameWorld.UserInputProfile = world.UserInputProfile{
		ProfileType:       world.ProfileTypeControl,
		KeyMoveSelectNext: "n",
		KeyMoveSelectPrev: "p",
	}
	gameWorld.EditingControl = true
	gameWorld.KeyDown = "n"
	gameWorld.ActiveControl = &world.ControlNode{TargetEntityID: 7}
	gameWorld.ControlLabels[7] = component.ControlLabel{
		EntityIDs:       []int{7, 8, 9, 10},
		Width:           4,
		Height:          1,
		MaxLength:       4,
		Operation:       "append",
		Sources:         []string{"facing", "speed", ", "},
		SourceEntityIDs: []int{4, 5, -1},
	}
	for _, entityID := range []int{7, 8, 9, 10} {
		gameWorld.Ascii[entityID] = component.Ascii{Ascii: ' '}
	}
	gameWorld.ControlOptions[4] = component.ControlOptions{Current: 'N', Options: []rune{'N', 'E'}}
	gameWorld.ControlNumber[5] = component.ControlNumber{ValueCurrent: 0}

	next := applyUpdate(t, control.ServiceControl{}.GetUpdateFunc(gameWorld), gameWorld)
	if got := string([]rune{next.Ascii[7].Ascii, next.Ascii[8].Ascii, next.Ascii[9].Ascii, next.Ascii[10].Ascii}); got != "N0, " {
		t.Fatalf("expected appended control label %q, got %q", "N0, ", got)
	}

	next.KeyDown = "p"
	removed := applyUpdate(t, control.ServiceControl{}.GetUpdateFunc(next), next)
	if got := string([]rune{removed.Ascii[7].Ascii, removed.Ascii[8].Ascii, removed.Ascii[9].Ascii, removed.Ascii[10].Ascii}); got != "    " {
		t.Fatalf("expected control label to be cleared, got %q", got)
	}
}

func newControlTestWorld() world.World {
	n1 := &world.ControlNode{SelectableEntityID: 1, TargetEntityID: 4}
	n2 := &world.ControlNode{SelectableEntityID: 2, TargetEntityID: 5}
	n3 := &world.ControlNode{SelectableEntityID: 3, TargetEntityID: 6}
	n1.Next, n1.Prev = n2, n3
	n2.Next, n2.Prev = n3, n1
	n3.Next, n3.Prev = n1, n2

	gameWorld := world.NewWorldEmpty()
	gameWorld.ControlOrder["testroom"] = n1
	gameWorld.ActiveControl = n1
	gameWorld.UserInputProfile = world.UserInputProfile{
		ProfileType:       world.ProfileTypeControl,
		KeyMoveSelectNext: "n",
		KeyMoveSelectPrev: "p",
		KeySelect:         "s",
	}
	gameWorld.ControlNumber[4] = component.ControlNumber{ValueStart: 1, ValueCurrent: 1, ValueMax: 3}
	return gameWorld
}

func assertEmptyUpdateFunc(t *testing.T, update game.UpdateFunc) {
	t.Helper()
	if update.Order != 0 || update.UpdateFunc != nil {
		t.Fatal("expected empty update func")
	}
}

func assertNotEmptyUpdateFunc(t *testing.T, update game.UpdateFunc) {
	t.Helper()
	if update.Order == 0 || update.UpdateFunc == nil {
		t.Fatal("expected non-empty update func")
	}
}
