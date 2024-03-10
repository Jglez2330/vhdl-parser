package token


type FileSet struct{

}

func NewFileSet() *FileSet{
    return &FileSet{}
}

func (s *FileSet) AddFile(filename string, base, size int) *File{

}
