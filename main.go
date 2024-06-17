package main

import (
	"os"
	"vhdl/parser"
	"vhdl/token"
)

func main() {
	var p parser.Parser
	src, _ := os.ReadFile("test/UART.vhd")
	fset := token.NewFileSet()
	//fset.AddFile("UART.vhd", fset.Base(), len(src)) // register input "file"
	p.Init(fset, "UART.vhd", src, 0)
	p.ParseFile()

	// s.Init(file, src, nil /* no error handler */, 0)
	// for i := 0; i < 3000; i++ {
	// 	pos, tok, lit := s.Scan()
	// 	fmt.Printf("%i %i %s \n", pos, tok, lit)
	//        if tok == token.EOF {
	//            break
	//        }
	// }

}
