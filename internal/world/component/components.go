package component

type ComponentName string

const (
	C_POS        ComponentName = "pos"
	C_ASCII      ComponentName = "ascii"
	C_IMPASSABLE ComponentName = "impassable"
	C_MACHINE    ComponentName = "machine"
	C_PLAYER     ComponentName = "player"
	C_VISIBLE    ComponentName = "visible"
	C_WALKABLE   ComponentName = "walkable"
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

type Visible struct{}

type Walkable struct{}

type Machine struct {
	IsOn        bool
	MachineType MachineTypeName
}
type MachineTypeName string

const (
	MACHINENAME_RADIO MachineTypeName = "MACHINENAME_RADIO"
)
