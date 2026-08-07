package tests

import (
	component "go_ascii/internal"
	"go_ascii/internal/game"
	"go_ascii/internal/service/control"
	"go_ascii/internal/world"
	"testing"
)

func TestControlSelectFirstSelectable(t *testing.T) {
	tests := []struct {
		Name      string
		MakeWorld func() world.World
		Assert    func(world.World, game.UpdateFunc, *testing.T)
	}{
		{
			"not control inputprofile",
			func() world.World {
				w := world.NewWorldEmpty()
				w.UserInputProfile = world.UserInputProfile{}
				return w
			},
			Assert_IsEmptyUpdateFunc,
		},
		{
			"no selected, no key down",
			func() world.World {
				w := world.NewWorldEmpty()
				w = initworld_control(w)

				w.UserInputProfile = world.UserInputProfile{ProfileType: world.ProfileTypeControl}
				return w
			},
			Assert_IsEmptyUpdateFunc,
		},
		{
			"no selected, nextkey down",
			func() world.World {
				w := world.NewWorldEmpty()
				w = initworld_control(w)

				w.UserInputProfile = world.UserInputProfile{ProfileType: world.ProfileTypeControl, KeyMoveSelectNext: "e"}
				w.KeyDown = "e"
				return w
			},
			func(w world.World, u game.UpdateFunc, t *testing.T) {
				Assert_IsNotEmptyUpdateFunc(w, u, t)
				firstFocus := w.FocusedControl
				nw, err := u.UpdateFunc(w)
				if err != nil {
					t.Fatal(err)
				}
				if firstFocus.Next.Value != nw.FocusedControl.Value {
					t.Fatal("Wrong control focused")
				}
			},
		},
		{
			"no selected, prevkey down",
			func() world.World {
				w := world.NewWorldEmpty()
				w = initworld_control(w)

				w.UserInputProfile = world.UserInputProfile{ProfileType: world.ProfileTypeControl, KeyMoveSelectPrev: "p"}
				w.KeyDown = "p"
				return w
			},
			func(w world.World, u game.UpdateFunc, t *testing.T) {
				Assert_IsNotEmptyUpdateFunc(w, u, t)
				firstFocus := w.FocusedControl
				nw, err := u.UpdateFunc(w)
				if err != nil {
					t.Fatal(err)
				}
				if firstFocus.Prev.Value != nw.FocusedControl.Value {
					t.Fatal("Wrong control focused")
				}
			},
		},
		{
			"no selected, select down",
			func() world.World {
				w := world.NewWorldEmpty()
				w = initworld_control(w)

				w.UserInputProfile = world.UserInputProfile{ProfileType: world.ProfileTypeControl, KeySelect: "a"}
				w.KeyDown = "a"
				return w
			},
			func(w world.World, u game.UpdateFunc, t *testing.T) {
				Assert_IsNotEmptyUpdateFunc(w, u, t)
				if !w.SelectedControl.IsEmpty() {
					t.Fatal("selected control should be nil when initialized!")
				}
				nw, err := u.UpdateFunc(w)
				if err != nil {
					t.Fatal(err)
				}
				if nw.SelectedControl.IsEmpty() {
					t.Fatal("Nothing selected, 1st should be!")
				} else if nw.SelectedControl.Value != 1 {
					t.Fatal("Wrong control selected!")
				}
			},
		},
		{
			"1 selected, nextkey down",
			func() world.World {
				w := world.NewWorldEmpty()
				w = initworld_control(w)
				w.SelectedControl = w.ControlSelectableOrder["testroom"]
				w.Selectable[1] = component.Selectable{TargetEntityId: 4}
				w.ControlNumber[4] = component.ControlNumber{ValueStart: 1, ValueCurrent: 1, ValueMax: 3}

				w.UserInputProfile = world.UserInputProfile{ProfileType: world.ProfileTypeControl, KeyMoveSelectNext: "n"}
				w.KeyDown = "n"
				return w
			},
			func(w world.World, u game.UpdateFunc, t *testing.T) {
				Assert_IsNotEmptyUpdateFunc(w, u, t)
				nw, err := u.UpdateFunc(w)
				if err != nil {
					t.Fatal(err)
				}
				if nw.SelectedControl.IsEmpty() || nw.SelectedControl.Value != 1 {
					t.Fatal("Selected control changed unexpectedly")
				} else if nw.ControlNumber[4].ValueCurrent != 2 {
					t.Fatal("Value not incremented correctly")
				}
			},
		},
	}
	for _, tv := range tests {
		t.Run(tv.Name, func(t *testing.T) {
			s := control.ServiceControl{}
			w := tv.MakeWorld()
			res := s.GetUpdateFunc(w)
			tv.Assert(w, res, t)
		})
	}
}
func initworld_control(w world.World) world.World {
	n1 := world.Node[int]{}
	n2 := world.Node[int]{}
	n3 := world.Node[int]{}
	n1 = world.Node[int]{Next: &n2, Prev: &n3, Value: 1}
	n2 = world.Node[int]{Next: &n3, Prev: &n1, Value: 2}
	n3 = world.Node[int]{Next: &n1, Prev: &n2, Value: 3}
	w.ControlSelectableOrder = make(map[string]world.Node[int])
	w.ControlSelectableOrder["testroom"] = n1
	w.FocusedControl = n1
	w.SelectedControl = world.Node[int]{}
	return w
}
func Assert_IsEmptyUpdateFunc(w world.World, u game.UpdateFunc, t *testing.T) {
	if u.Order == 0 && u.UpdateFunc == nil {
		return
	}
	t.Fatal("Update func is not empty!")
}

func Assert_IsNotEmptyUpdateFunc(w world.World, u game.UpdateFunc, t *testing.T) {
	if u.Order != 0 && u.UpdateFunc != nil {
		return
	}
	t.Fatal("Update func is empty!")
}
