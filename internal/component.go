package component

const (
	NamePosition            = "pos"
	NameASCII               = "ascii"
	NameImpassable          = "impassable"
	NamePlayer              = "player"
	NameInteractable        = "interactable"
	InteractionTypeDoor     = "door"
	InteractionTypeTerminal = "terminal"
	NameControlTypeNumber   = "controlnumber"
	NameControlTypeList     = "controllist"
	NameSelectable          = "selectable"
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
type ControlList struct {
	Current rune
	List    []rune
}
type ControlNumber struct {
	ValueStart   int
	ValueCurrent int
	ValueMax     int
}
type Selectable struct {
	UnfocusedASCII   rune
	FocusedASCII     rune
	SelectedASCII    rune
	TargetEntityId   int
	TargetEntityName string
}
