package component

const (
	C_POS   ComponentName = "pos"
	C_ASCII ComponentName = "ascii"
)

type Position struct {
	X int
	Y int
}
type Ascii struct {
	Ascii rune
}

type ComponentName string
