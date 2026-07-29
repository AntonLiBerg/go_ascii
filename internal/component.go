package component

const (
	NamePosition   = "pos"
	NameASCII      = "ascii"
	NameImpassable = "impassable"
	NamePlayer     = "player"
	NameMachine    = "machine"
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

type Machine struct {
	IsOn bool
}
