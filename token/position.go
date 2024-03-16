package token

import "strconv"

// This is based on https://pkg.go.dev/go/token#File
type Position struct {
	Filename string // filename, if any
	Offset   int    // offset, starting at 0
	Line     int    // line number, starting at 1
	Column   int    // column number, starting at 1 (byte count)
}

// Position describes an arbitrary source position including the file, line, and column location. A Position is valid if the line number is > 0.
func (pos *Position) IsValid() bool {
	return pos.Line > 0
}

func (pos *Position) String() string {
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

// Pos is a compact encoding of a source position within a file set. It can be converted into a Position for a more convenient, 
// but much larger, representation. The Pos value for a given file is a number in the range [base, base+size], where base and size 
// are specified when a file is added to the file set. The difference between a Pos value and the corresponding file base corresponds 
// to the byte offset of that position (represented by the Pos value) from the beginning of the file. Thus, the file base offset is the 
// Pos value representing the first byte in the file.
type Pos int

const NoPos Pos = 0

func (p Pos) IsValid() bool{
    return p != NoPos
}
