package world

import (
	"fmt"
	cmp "go_ascii/component"
	usr "go_ascii/user"
)

type World struct {
	UserInputProfile usr.UserInputProfile
	StateUser        usr.UserState
	UserInput        map[string]bool
	NextEnt          int
	Entities         []int
	Pos              map[int]cmp.Position
	Ascii            map[int]cmp.Ascii
	Impassable       map[int]cmp.Impassable
	Machine          map[int]cmp.Machine
	Tags             map[int]cmp.Tags
	EByTag           map[cmp.Tag]map[int]bool
	EByPos           map[cmp.Position]int
	HasChanged       bool
	IterationNr      int
}

func NewWorldEmpty() World {
	return World{
		UserInputProfile: usr.NewUserInputProfileEmpty(),
		StateUser:        usr.S_playing,
		UserInput:        map[string]bool{},
		NextEnt:          0,
		Entities:         []int{},
		Pos:              map[int]cmp.Position{},
		Ascii:            map[int]cmp.Ascii{},
		Impassable:       map[int]cmp.Impassable{},
		Machine:          map[int]cmp.Machine{},
		Tags:             map[int]cmp.Tags{},
		EByTag:           map[cmp.Tag]map[int]bool{},
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
		StateUser:        w.StateUser,
		UserInput:        make(map[string]bool, len(w.UserInput)),
		NextEnt:          w.NextEnt,
		Entities:         append([]int(nil), w.Entities...),
		Pos:              make(map[int]cmp.Position, len(w.Pos)),
		Ascii:            make(map[int]cmp.Ascii, len(w.Ascii)),
		Impassable:       make(map[int]cmp.Impassable, len(w.Impassable)),
		Machine:          make(map[int]cmp.Machine, len(w.Machine)),
		Tags:             make(map[int]cmp.Tags, len(w.Tags)),
		EByTag:           make(map[cmp.Tag]map[int]bool, len(w.EByTag)),
		EByPos:           make(map[cmp.Position]int, len(w.Pos)),
		HasChanged:       w.HasChanged,
		IterationNr:      w.IterationNr,
	}

	for key, value := range w.UserInput {
		clone.UserInput[key] = value
	}

	for id, pos := range w.Pos {
		clone.Pos[id] = pos
		clone.EByPos[pos] = id
	}

	for id, ascii := range w.Ascii {
		clone.Ascii[id] = ascii
	}

	for id, impassable := range w.Impassable {
		clone.Impassable[id] = impassable
	}

	for id, machine := range w.Machine {
		clone.Machine[id] = machine
	}

	for id, tags := range w.Tags {
		cloneTags := cmp.Tags{Vals: make(map[cmp.Tag]bool, len(tags.Vals))}
		for tag, ok := range tags.Vals {
			cloneTags.Vals[tag] = ok
			if !ok {
				continue
			}
			if clone.EByTag[tag] == nil {
				clone.EByTag[tag] = make(map[int]bool)
			}
			clone.EByTag[tag][id] = true
		}
		clone.Tags[id] = cloneTags
	}

	return clone
}
func (w *World) ClearUserInput() {
	clear(w.UserInput)
}

func (w *World) SetKeyDown(key string) {
	w.UserInput[key] = true
}

func (w World) IsKeyDown(key string) bool {
	return w.UserInput[key]
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

func (w World) AddUserInput(key string, isDown bool) World {
	world := w.Clone()
	world.UserInput[key] = isDown
	return world
}

func (w World) AddPosition(eID int, pos cmp.Position) World {
	world := w.Clone()
	if oldPos, ok := world.Pos[eID]; ok {
		delete(world.EByPos, oldPos)
	}
	world.Pos[eID] = pos
	world.EByPos[pos] = eID
	return world
}

func (w World) AddAscii(eID int, ascii cmp.Ascii) World {
	world := w.Clone()
	world.Ascii[eID] = ascii
	return world
}

func (w World) AddImpassable(eID int) World {
	world := w.Clone()
	world.Impassable[eID] = cmp.Impassable{}
	return world
}

func (w World) AddMachine(eID int, machine cmp.Machine) World {
	world := w.Clone()
	world.Machine[eID] = machine
	return world
}

func (w World) AddTags(eID int, tags cmp.Tags) World {
	world := w.Clone()
	world.removeTags(eID)
	world.Tags[eID] = cloneTags(tags)
	for tag, ok := range world.Tags[eID].Vals {
		if !ok {
			continue
		}
		if world.EByTag[tag] == nil {
			world.EByTag[tag] = map[int]bool{}
		}
		world.EByTag[tag][eID] = true
	}
	return world
}

func (w World) AddTag(eID int, tag cmp.Tag) World {
	world := w.Clone()
	tags := world.Tags[eID]
	if tags.Vals == nil {
		tags.Vals = map[cmp.Tag]bool{}
	}
	tags.Vals[tag] = true
	world.Tags[eID] = tags
	if world.EByTag[tag] == nil {
		world.EByTag[tag] = map[int]bool{}
	}
	world.EByTag[tag][eID] = true
	return world
}

func (w *World) AddEntity(pos [2]int, compWithVals map[cmp.ComponentName][]string) error {
	eId := w.MakeNewEntityId()
	for name, vals := range compWithVals {
		switch name {
		case cmp.C_POS:
			entityPos := cmp.Position{X: pos[0], Y: pos[1]}
			w.Pos[eId] = entityPos
			w.EByPos[entityPos] = eId
		case cmp.C_ASCII:
			if len(vals) != 1 || len(vals[0]) != 1 {
				return fmt.Errorf("Required values are incorrect for %s", cmp.C_ASCII)
			}
			w.Ascii[eId] = cmp.Ascii{Ascii: []rune(vals[0])[0]}
		case cmp.C_IMPASSABLE:
			w.Impassable[eId] = cmp.Impassable{}
		case cmp.C_MACHINE:
			if len(vals) != 1 {
				return fmt.Errorf("Required values are incorrect for %s", cmp.C_MACHINE)
			}
			machineType := cmp.MachineTypeName(vals[0])
			switch machineType {
			case cmp.MACHINENAME_RADIO:
				w.Machine[eId] = cmp.Machine{MachineType: machineType}
			default:
				return fmt.Errorf("Machine type does not exist %s", vals[0])
			}
		case cmp.C_TAGS:
			tags := cmp.Tags{Vals: make(map[cmp.Tag]bool, len(vals))}
			for _, value := range vals {
				tag := cmp.Tag(value)
				tags.Vals[tag] = true
				if w.EByTag[tag] == nil {
					w.EByTag[tag] = make(map[int]bool)
				}
				w.EByTag[tag][eId] = true
			}
			w.Tags[eId] = tags
		default:
			return fmt.Errorf("component does not exist %s", name)
		}
	}
	return nil
}

func (w *World) removeTags(eID int) {
	tags := w.Tags[eID]
	for tag, ok := range tags.Vals {
		if !ok {
			continue
		}
		delete(w.EByTag[tag], eID)
		if len(w.EByTag[tag]) == 0 {
			delete(w.EByTag, tag)
		}
	}
}

func cloneTags(tags cmp.Tags) cmp.Tags {
	clone := cmp.Tags{Vals: map[cmp.Tag]bool{}}
	for tag, ok := range tags.Vals {
		clone.Vals[tag] = ok
	}
	return clone
}
