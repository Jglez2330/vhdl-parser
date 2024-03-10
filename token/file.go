package token

import "sync"



type File struct{
    filename string
    base int
    size int
    line_start[] int
    column_info[] columnInfo

    //This handler can be used by multiple go corutines a mutex is need or semaphore to control access
    mutex sync.Mutex

}

type columnInfo struct{
    Offset int
    Filename string
    Line int
    Column int
}

func (f *File) AddLine(offset int){
    f.mutex.Lock()
    defer f.mutex.Unlock()
    prevLineIndex := len(f.line_start) - 1
    if offset < f.size && (prevLineIndex < 0 || f.line_start[prevLineIndex] < offset){
        f.line_start = append(f.line_start, offset)
    }
}


func (f *File) AddLineColumnInfo(offset int, filename string, line, column int){
    f.mutex.Lock()
    defer f.mutex.Unlock()
    prevColumnInfo := len(f.column_info) - 1
    if offset < f.size && (prevColumnInfo < 0 || f.line_start[prevColumnInfo] < offset){
        f.column_info = append(f.column_info, columnInfo{offset, filename, line, column})
    }
}

func (f *File) Base() int{
    return f.base
}

func (f *File) Line(p Pos) int{

}
