package scanner

//A mode value is a set of flags (or 0). They control scanner behavior.
type Mode uint

const (
	ScanComments Mode = 1 << iota // return comments as COMMENT tokens

)
