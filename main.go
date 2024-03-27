package main

import (
	"fmt"
	"os"
	"vhdl/scanner"
	"vhdl/token"
)

func main() {
	var s scanner.Scanner
	src, _ := os.ReadFile("test/UART.vhd")
	fset := token.NewFileSet()
	file := fset.AddFile("UART.vhd", fset.Base(), len(src)) // register input "file"
	s.Init(file, src, nil /* no error handler */, 0)
	for i := 0; i < 3000; i++ {
		pos, tok, lit := s.Scan()
		fmt.Printf("%i %i %s \n", pos, tok, lit)
	}

}
