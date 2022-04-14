package constant

type Prior uint8

//go:generate stringer -type Prior
const (
	Super   Prior = 1<<8 - 1
	Admin   Prior = 1<<7 - 1
	Normal  Prior = 1<<6 - 1
	Visitor Prior = 1
)
