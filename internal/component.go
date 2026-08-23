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
	NameControlOptions      = "controloptions"
	NameControlLabel        = "controllabel"
	NameSelectable          = "selectable"
	NameSelectableButtonVariable = "selectablebuttonvariable"
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
type ControlOptions struct {
	Current rune
	Options []rune
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
type SelectableButtonVariable struct {
	UnfocusedASCII   rune
	FocusedASCII     rune
	VariableUpdate   string
}
type ControlLabel struct {
	EntityIDs       []int
	MaxLength       int
	Width           int
	Height          int
	Operation       string
	Sources         []string
	SourceEntityIDs []int
	Content         string
	History         []string
}
