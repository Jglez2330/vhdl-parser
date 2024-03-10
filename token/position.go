package token

import "strconv"

// This is based on https://pkg.go.dev/go/token#File
type Postion struct {
	Filename string // filename, if any
	Offset   int    // offset, starting at 0
	Line     int    // line number, starting at 1
	Column   int    // column number, starting at 1 (byte count)
}

// Position describes an arbitrary source position including the file, line, and column location. A Position is valid if the line number is > 0.
func (pos *Postion) IsValid() bool {
	return pos.Line > 0
}

func (pos *Postion) String() string {
	result := ""
	if result = pos.Filename; pos.IsValid() {
		line := strconv.FormatInt(int64(pos.Line), 10)
		if pos.Filename != "" {
			result += ":"
		}
		result += line
		if pos.Column > 0 {
			result += ":" + strconv.FormatInt(int64(pos.Column), 10)
		}

	}

	if result == "" {
		result = "-"
	}
	return result
}
