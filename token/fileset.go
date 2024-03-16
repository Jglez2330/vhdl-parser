package token

import (
	"fmt"
	"sort"
	"sync"
)

type FileSet struct {
	base  int
	mutex sync.RWMutex // protects the file set
	files []*File      // list of files in order

}

func NewFileSet() *FileSet {
	return &FileSet{base: 1}
}

func (s *FileSet) AddFile(filename string, base, size int) *File {
	newFile := &File{filename: filename, base: base, size: size}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if base < 0 {
		base = s.base
	}
	if base < s.base {
		panic(fmt.Sprintf("invalid base %d (should be >= %d)", base, s.base))
	}
	newFile.base = base
	if size < 0 {
		panic(fmt.Sprintf("invalid size %d (should be >= 0)", size))
	}
	// base >= s.base && size >= 0
	base += size + 1 // +1 because EOF also has a position
	if base < 0 {
		panic("token.Pos offset overflow (> 2G of source code in file set)")
	}
	s.base = base
	s.files = append(s.files, newFile)
	return newFile
}

func (s *FileSet) Base() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.base
}

func (s *FileSet) File(p Pos) (f *File) {
	f = nil
	if p == NoPos {
		return nil
	}
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	index := sort.Search(len(s.files), func(index int) bool { return s.files[index].base > int(p) }) - 1
	if index < 0 {
		return nil
	}
	if int(p) <= s.files[index].size+s.files[index].base {
		f = s.files[index]
		return f
	}
	return nil
}

func (s *FileSet) Iterate(f func(*File) bool) {
	// Iterate may be write or read
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, file := range s.files {
		if !f(file) || file == nil {
			break
		}
	}

}

func (s *FileSet) Position(p Pos) (pos Position) {
	return s.PositionFor(p, true)
}

func (s *FileSet) PositionFor(p Pos, adjusted bool) (pos Position) {

	if p == NoPos {
		return
	}
	file := s.File(p)

	if file == nil {
		return
	}

	return file.Position(p)
}

func (s *FileSet) Read(decode func(any) error) error {
	//TODO
	return nil
}

func (s *FileSet) RemoveFile(file *File) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if index := sort.Search(len(s.files), func(i int) bool { return s.files[i].base > file.base }) - 1; index >= 0 && file == s.files[index] {
		s.files = append(s.files[:index], s.files[index+1:]...)
	}
}

func (s *FileSet) Write(encode func(any) error) error {
	//TODO
	return nil
}
