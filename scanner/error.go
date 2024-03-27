package scanner

import (
	"fmt"
	"io"
	"sort"
	"vhdl/token"
)

type Error struct {
	Pos token.Position
	Msg string
}

func (e Error) Error() string {
	return ""
}

type ErrorHandler func(pos token.Position, msg string)

type ErrorList []*Error

func (p *ErrorList) Add(pos token.Position, msg string) {
	*p = append(*p, &Error{pos, msg})
}

func (p ErrorList) Err() error {
	if len(p) == 0 {
		return nil
	}
	return p
}

func (p ErrorList) Error() string {
	switch len(p) {
	case 0:
		return "no errors"
	case 1:
		return p[0].Error()
	}
	return fmt.Sprintf("%s (and %d more errors)", p[0], len(p)-1)
}

func (p ErrorList) Len() int {
	return len(p)
}

func (p ErrorList) Less(i, j int) bool {
	// Extract references to error positions for easier readability
	e, f := &p[i].Pos, &p[j].Pos

	// Compare filenames first
	if e.Filename != f.Filename {
		return e.Filename < f.Filename
	}

	// Compare line numbers
	if e.Line != f.Line {
		return e.Line < f.Line
	}

	// Compare column numbers
	if e.Column != f.Column {
		return e.Column < f.Column
	}

	// If all other comparisons are equal, compare error messages
	return p[i].Msg < p[j].Msg
}

func (p *ErrorList) Reset() {
	*p = (*p)[0:0]
}

func (p ErrorList) Sort() {
	sort.Sort(p)
}

func (p ErrorList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func PrintError(w io.Writer, err error) {
	if list, ok := err.(ErrorList); ok {
		for _, e := range list {
			fmt.Fprintf(w, "%s\n", e)
		}
	} else if err != nil {
		fmt.Fprintf(w, "%s\n", err)
	}
}
