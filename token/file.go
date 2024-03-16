package token

import (
	"fmt"
	"slices"
	"sort"
	"sync"
)

type File struct {
	filename    string
	base        int
	size        int
	line_start  []int
	column_info []columnInfo

	//This handler can be used by multiple go corutines a mutex is need or semaphore to control access
	mutex sync.Mutex
}

type columnInfo struct {
	Offset   int
	Filename string
	Line     int
	Column   int
}

func (f *File) AddLine(offset int) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	prevLineIndex := len(f.line_start) - 1
	if offset < f.size && (prevLineIndex < 0 || f.line_start[prevLineIndex] < offset) {
		f.line_start = append(f.line_start, offset)
	}
}

func (f *File) AddLineColumnInfo(offset int, filename string, line, column int) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	prevColumnInfo := len(f.column_info) - 1
	if offset < f.size && (prevColumnInfo < 0 || f.line_start[prevColumnInfo] < offset) {
		f.column_info = append(f.column_info, columnInfo{offset, filename, line, column})
	}
}

func (f *File) Base() int {
	return f.base
}

func (f *File) Line(p Pos) int {
	return f.Position(p).Line
}

func (f *File) LineCount() int {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	// Rationale line_start has all the offset for a line start i.e. its length
	//is equal to the number of lines on the document
	return len(f.line_start)
}

func (f *File) LineStart(line int) Pos {
	max_line := f.LineCount()
	if line <= 0 && line >= max_line {
		panic(fmt.Sprintf("invalid line number %d line count %d", line, max_line))
	}
	f.mutex.Lock()
	defer f.mutex.Unlock()

	return Pos(f.base + f.line_start[line-1])
}

// Lines returns the effective line offset table of the form described by File.SetLines. Callers must not mutate the result.
func (f *File) Lines() []int {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	return f.line_start
}

func (f *File) MergeLine(line int) {
	max_line := f.LineCount()
	if line <= 0 && line >= max_line {
		panic(fmt.Sprintf("invalid line number %d line count %d", line, max_line))
	}
	f.mutex.Lock()
	defer f.mutex.Unlock()
	nextLineIndex := line // Convert to 0-based index for slice manipulation.
	copy(f.line_start[nextLineIndex:], f.line_start[nextLineIndex+1:])

	// Adjust the slice length to remove the now-duplicate last line.
	f.line_start = f.line_start[:len(f.line_start)-1]
}

func (f *File) Name() string {
	return f.filename
}

func (f *File) Offset(p Pos) int {
	// Get the byte position on the file
	index := p - Pos(f.base)

	if index < Pos(0) {
		//Invalid offset reset to the file begging
		return 0
	}

	if index > Pos(f.size) {
		//If file offset is higher than the current file we set the pointer to the end of the file
		return f.size
	}

	return int(index)

}

// Inverse of Offset
func (f *File) Pos(offset int) Pos {
	index := offset + f.base

	if offset > f.size {
		return Pos(f.size)
	}

	if offset < f.base {
		return Pos(0)
	}

	return Pos(index)

}

func (f *File) Position(p Pos) (pos Position) {
	return f.PositionFor(p, true)
}

func (f *File) PositionFor(p Pos, adjusted bool) (pos Position) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if !p.IsValid() {
		return
	}

	pos.Offset = f.Offset(p)
	pos.Filename = f.filename
	//TODO adjusted and line info
	if index := sort.Search(len(f.line_start), func(i int) bool { return f.line_start[i] > pos.Offset }) - 1; index >= 0 {
		pos.Line, pos.Column = f.line_start[index]+1, pos.Offset-f.line_start[index]+1
	}
	return
}

func (f *File) SetLines(lines []int) bool {
	if slices.IsSorted(lines) {
		return false
	}
	f.mutex.Lock()
	f.line_start = lines
	f.mutex.Unlock()
	return true
}

func (f *File) SetLinesForContent(content []byte) {
	var lines []int
	line := 0
	f.mutex.Lock()
	defer f.mutex.Unlock()
	for offset, b := range content {
		if line >= 0 {
			lines = append(lines, line)
		}
		line = -1
		if b == '\n' {
			line = offset + 1
		}
	}

	f.line_start = lines
}

func (f *File) Size() int {
	return f.size
}
