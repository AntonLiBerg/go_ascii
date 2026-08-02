package component

const (
	NamePosition     = "pos"
	NameASCII        = "ascii"
	NameImpassable   = "impassable"
	NamePlayer       = "player"
	NameInteractable = "interactable"
	NameDoor         = "door"
	NameHelm         = "helm"
	NameCommandTable = "commandtable"
	NameBunkBed      = "bunkbed"
	NamePrisonBars   = "prisonbars"
	NameWall         = "wall"
	NameWindow       = "window"
)

type Position struct {
	Room string
	X    int
	Y    int
}

type Layer struct {
	Nr int
}

type Ascii struct {
	Ascii rune
}

type Impassable struct{}

type Player struct{}

type Interactable struct {
	InteractionType string
}

type Door struct{}

type Helm struct{}

type CommandTable struct{}

type BunkBed struct{}

type PrisonBars struct{}

type Wall struct{}

type Window struct{}
