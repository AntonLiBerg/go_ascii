package component

const (
	NamePosition   = "pos"
	NameASCII      = "ascii"
	NameImpassable = "impassable"
	NamePlayer     = "player"
	NameDoor       = "door"
)

type Position struct {
	X int
	Y int
}

type Ascii struct {
	Ascii rune
}

type Impassable struct{}

type Player struct{}

type Door struct {
	IsInteractable bool
}
