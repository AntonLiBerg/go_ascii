package interaction

import (
	wrld "go_ascii/internal/world"
	cmp "go_ascii/internal/world/component"
	usr "go_ascii/internal/world/user"
	"testing"
)

func TestServiceTurnsOnSingleNeighborAndClearsInteractKey(t *testing.T) {
	world := wrld.NewWorldEmpty()
	world.UserInputProfile = usr.UserInputProfile{KeyInteract: "e"}

	addTestEntity(t, &world, [2]int{1, 1}, map[cmp.ComponentName][]string{
		cmp.C_POS: {}, cmp.C_ASCII: {"o"}, cmp.C_PLAYER: {},
	})
	addTestEntity(t, &world, [2]int{2, 1}, map[cmp.ComponentName][]string{
		cmp.C_POS: {}, cmp.C_ASCII: {"R"}, cmp.C_IMPASSABLE: {},
		cmp.C_MACHINE: {string(cmp.MACHINENAME_RADIO)},
	})
	world.SetKeyDown("e")

	result := ServiceTurnOnMachine{}.GetUpdateFunc(world)
	if result.UpdateFunc == nil {
		t.Fatal("expected interact update func")
	}
	result.UpdateFunc(&world)

	for _, machine := range world.Machine {
		if !machine.IsOn {
			t.Fatal("expected machine to be on after interaction")
		}
	}
	if world.UserInput["e"] {
		t.Fatal("expected interact key to be cleared after interaction")
	}
	if !world.HasChanged {
		t.Fatal("expected world to be marked changed after turning on machine")
	}
}

func addTestEntity(t *testing.T, world *wrld.World, position [2]int, components map[cmp.ComponentName][]string) {
	t.Helper()
	if err := world.AddEntity(position, components); err != nil {
		t.Fatalf("AddEntity returned error: %v", err)
	}
}
