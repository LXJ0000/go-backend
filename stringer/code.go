package stringer

//go:generate stringer -type Code -linecomment
type Code int

const (
	CODE_OK Code = iota // success
	CODE_ERROR // fail
)
