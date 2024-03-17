package scanner

import(
    "vhdl/token"
)

type Error struct {
	Pos token.Position
	Msg string
}


func (e Error) Error() string{
    return ""
}


type ErrorHandler func(pos token.Position, msg string){

}


type ErrorList []*Error{

}

func (p *ErrorList) Add(pos token.Position, msg string){

}


func (p ErrorList) Err() error{

}


func (p ErrorList) Error() string{

}

func (p ErrorList) Len() int{

}

func (p ErrorList) Less(i, j int) bool{

}

func (p *ErrorList) Reset(){

}

func (p ErrorList) Sort(){

}

func (p ErrorList) Swap(i, j int){

}

func PrintError(w io.Writer, err error){

}
