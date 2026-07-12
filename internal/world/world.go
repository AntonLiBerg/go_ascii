package world

import (
	"fmt"
	cmp "go_ascii/internal/world/component"
	usr "go_ascii/internal/world/user"
)

type World struct {
	UserInputProfile usr.UserInputProfile
	StateUser        map[usr.UserState]bool
	UserInput        map[usr.Interaction]bool
	EByPos           map[cmp.Position]int
	HasChanged       bool
	IterationNr      int
	MenuChoices      MenuChoices
	//Entities
	NextEnt    int
	Entities   []int
	Pos        map[int]cmp.Position
	Ascii      map[int]cmp.Ascii
	Impassable map[int]cmp.Impassable
	Player     map[int]cmp.Player
	Visible    map[int]cmp.Visible
	Walkable   map[int]cmp.Walkable
	Machine    map[int]cmp.Machine
}

func NewWorldEmpty() World {
	return World{
		UserInputProfile: usr.NewUserInputProfileEmpty(),
		StateUser:        map[usr.UserState]bool{usr.S_playing: true},
		UserInput:        map[usr.Interaction]bool{},
		MenuChoices:      MenuChoices{Choices: []MenuChoice{}},
		NextEnt:          0,
		Entities:         []int{},
		Pos:              map[int]cmp.Position{},
		Ascii:            map[int]cmp.Ascii{},
		Impassable:       map[int]cmp.Impassable{},
		Player:           map[int]cmp.Player{},
		Visible:          map[int]cmp.Visible{},
		Walkable:         map[int]cmp.Walkable{},
		Machine:          map[int]cmp.Machine{},
		EByPos:           map[cmp.Position]int{},
		HasChanged:       false,
		IterationNr:      0,
	}
}
func NewWorld(aMap map[[2]int]rune, entities map[rune]string, components map[string]map[cmp.ComponentName][]string) (World, error) {
	world := NewWorldEmpty()
	for pos, ch := range aMap {
		eName := entities[ch]
		eComps := components[eName]
		err := world.AddEntity(pos, eComps)
		if err != nil {
			return world, err
		}
	}
	return world, nil
}
func (w *World) Clone() World {
	clone := World{
		UserInputProfile: w.UserInputProfile,
		StateUser:        make(map[usr.UserState]bool, len(w.StateUser)),
		UserInput:        make(map[usr.Interaction]bool, len(w.UserInput)),
		MenuChoices:      cloneMenuChoices(w.MenuChoices),
		NextEnt:          w.NextEnt,
		Entities:         append([]int(nil), w.Entities...),
		Pos:              make(map[int]cmp.Position, len(w.Pos)),
		Ascii:            make(map[int]cmp.Ascii, len(w.Ascii)),
		Impassable:       make(map[int]cmp.Impassable, len(w.Impassable)),
		Player:           make(map[int]cmp.Player, len(w.Player)),
		Visible:          make(map[int]cmp.Visible, len(w.Visible)),
		Walkable:         make(map[int]cmp.Walkable, len(w.Walkable)),
		Machine:          make(map[int]cmp.Machine, len(w.Machine)),
		EByPos:           make(map[cmp.Position]int, len(w.Pos)),
		HasChanged:       w.HasChanged,
		IterationNr:      w.IterationNr,
	}

	for state, active := range w.StateUser {
		clone.StateUser[state] = active
	}

	for key, value := range w.UserInput {
		clone.UserInput[key] = value
	}

	for id, pos := range w.Pos {
		clone.setPosition(id, pos)
	}

	for id, ascii := range w.Ascii {
		clone.setAscii(id, ascii)
	}

	for id, impassable := range w.Impassable {
		clone.setImpassable(id, impassable)
	}

	for id, player := range w.Player {
		clone.Player[id] = player
	}

	for id, visible := range w.Visible {
		clone.Visible[id] = visible
	}

	for id, walkable := range w.Walkable {
		clone.Walkable[id] = walkable
	}

	for id, machine := range w.Machine {
		clone.setMachine(id, machine)
	}

	return clone
}
func (w *World) ClearUserInput() {
	clear(w.UserInput)
}

func (w *World) SetKeyDown(interaction usr.Interaction) {
	w.UserInput[interaction] = true
}

func (w *World) SetUserState(state usr.UserState, isActive bool) {
	if w.StateUser == nil {
		w.StateUser = map[usr.UserState]bool{}
	}
	w.StateUser[state] = isActive
}

func (w World) HasUserState(state usr.UserState) bool {
	return w.StateUser[state]
}

func (w World) HasComponent(eID int, component cmp.ComponentName) bool {
	switch component {
	case cmp.C_POS:
		_, ok := w.Pos[eID]
		return ok
	case cmp.C_ASCII:
		_, ok := w.Ascii[eID]
		return ok
	case cmp.C_IMPASSABLE:
		_, ok := w.Impassable[eID]
		return ok
	case cmp.C_PLAYER:
		_, ok := w.Player[eID]
		return ok
	case cmp.C_VISIBLE:
		_, ok := w.Visible[eID]
		return ok
	case cmp.C_WALKABLE:
		_, ok := w.Walkable[eID]
		return ok
	case cmp.C_MACHINE:
		_, ok := w.Machine[eID]
		return ok
	default:
		return false
	}
}

func (w World) IsKeyDown(interaction usr.Interaction) bool {
	return w.UserInput[interaction]
}

func (w *World) MakeNewEntityId() int {
	w.Entities = append(w.Entities, w.NextEnt)
	w.NextEnt++
	return w.NextEnt - 1
}

func (w World) AddNewEntity() (World, int) {
	world := w.Clone()
	return world, world.MakeNewEntityId()
}

func (w World) AddUserInput(interaction usr.Interaction, isDown bool) World {
	world := w.Clone()
	world.UserInput[interaction] = isDown
	return world
}

func (w World) AddMenuChoices(menuChoices MenuChoices) World {
	world := w.Clone()
	world.MenuChoices = cloneMenuChoices(menuChoices)
	return world
}

func (w World) AddPosition(eID int, pos cmp.Position) World {
	world := w.Clone()
	world.setPosition(eID, pos)
	return world
}

func (w World) AddAscii(eID int, ascii cmp.Ascii) World {
	world := w.Clone()
	world.setAscii(eID, ascii)
	return world
}

func (w World) AddImpassable(eID int) World {
	world := w.Clone()
	world.setImpassable(eID, cmp.Impassable{})
	return world
}

func (w World) AddPlayer(eID int) World {
	world := w.Clone()
	world.Player[eID] = cmp.Player{}
	return world
}

func (w World) AddVisible(eID int) World {
	world := w.Clone()
	world.Visible[eID] = cmp.Visible{}
	return world
}

func (w World) AddWalkable(eID int) World {
	world := w.Clone()
	world.Walkable[eID] = cmp.Walkable{}
	return world
}

func (w World) AddMachine(eID int, machine cmp.Machine) World {
	world := w.Clone()
	world.setMachine(eID, machine)
	return world
}

func (w *World) AddEntity(pos [2]int, compWithVals map[cmp.ComponentName][]string) error {
	eId := w.MakeNewEntityId()
	for name, vals := range compWithVals {
		switch name {
		case cmp.C_POS:
			entityPos := cmp.Position{X: pos[0], Y: pos[1]}
			w.setPosition(eId, entityPos)
		case cmp.C_ASCII:
			if len(vals) != 1 || len(vals[0]) != 1 {
				return fmt.Errorf("Required values are incorrect for %s", cmp.C_ASCII)
			}
			w.setAscii(eId, cmp.Ascii{Ascii: []rune(vals[0])[0]})
		case cmp.C_IMPASSABLE:
			if len(vals) != 0 {
				return fmt.Errorf("Required values are incorrect for %s", cmp.C_IMPASSABLE)
			}
			w.setImpassable(eId, cmp.Impassable{})
		case cmp.C_PLAYER:
			if len(vals) != 0 {
				return fmt.Errorf("Required values are incorrect for %s", cmp.C_PLAYER)
			}
			w.Player[eId] = cmp.Player{}
		case cmp.C_VISIBLE:
			if len(vals) != 0 {
				return fmt.Errorf("Required values are incorrect for %s", cmp.C_VISIBLE)
			}
			w.Visible[eId] = cmp.Visible{}
		case cmp.C_WALKABLE:
			if len(vals) != 0 {
				return fmt.Errorf("Required values are incorrect for %s", cmp.C_WALKABLE)
			}
			w.Walkable[eId] = cmp.Walkable{}
		case cmp.C_MACHINE:
			if len(vals) != 1 {
				return fmt.Errorf("Required values are incorrect for %s", cmp.C_MACHINE)
			}
			machineType := cmp.MachineTypeName(vals[0])
			switch machineType {
			case cmp.MACHINENAME_RADIO:
				w.setMachine(eId, cmp.Machine{MachineType: machineType})
			default:
				return fmt.Errorf("Machine type does not exist %s", vals[0])
			}
		default:
			return fmt.Errorf("component does not exist %s", name)
		}
	}
	return nil
}

func (w *World) setPosition(eID int, pos cmp.Position) {
	if oldPos, ok := w.Pos[eID]; ok {
		delete(w.EByPos, oldPos)
	}
	w.Pos[eID] = pos
	w.EByPos[pos] = eID
}

func (w *World) setAscii(eID int, ascii cmp.Ascii) {
	w.Ascii[eID] = ascii
}

func (w *World) setImpassable(eID int, impassable cmp.Impassable) {
	w.Impassable[eID] = impassable
}

func (w *World) setMachine(eID int, machine cmp.Machine) {
	w.Machine[eID] = machine
}
