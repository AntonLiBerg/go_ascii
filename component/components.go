package component

type ComponentName string

const (
	C_POS        ComponentName = "pos"
	C_ASCII      ComponentName = "ascii"
	C_TAGS       ComponentName = "tags"
	C_IMPASSABLE ComponentName = "impassable"
	C_MACHINE    ComponentName = "machine"
)

type Position struct {
	X int
	Y int
}

type Ascii struct {
	Ascii rune
}

type Tag string

const (
	TAG_PLAYER Tag = "player"
)

type Tags struct {
	Vals map[Tag]bool
}

type Impassable struct{}

type Machine struct {
	IsOn        bool
	MachineType MachineTypeName
}
type MachineTypeName string

const (
	MACHINENAME_RADIO MachineTypeName = "MACHINENAME_RADIO"
)
