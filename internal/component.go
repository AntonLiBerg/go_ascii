package component

const (
	NamePosition                = "pos"
	NameASCII                   = "ascii"
	NameImpassable              = "impassable"
	NamePlayer                  = "player"
	NameInteractable            = "interactable"
	InteractionTypeDoor         = "door"
	InteractionTypeHelm         = "helm"
	InteractionTypeCommandTable = "commandtable"
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
