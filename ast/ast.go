package ast

import (
	"vhdl/token"
)

//VHDL AST

//We are following the VHDL 2019 LRM
//https://standards.ieee.org/standard/1076-2019.html

type Node struct {
	Pos token.Pos
	End token.Pos
}

// 3 Design entities and configurations
type EntityDeclaration struct {
	Identifier            Identifier
	EntityHeader          EntityHeader
	EntityDeclarativePart EntityDeclarativePart
	EntityStatementPart   *EntityStatementPart
	EntitySimpleName      *EntitySimpleName
	Node
}

type EntityHeader struct {
	FormalGenericClause *FormalGenericClause
	FormalPortClause    *FormalPortClause
	Node
}

type EntityDeclarativePart struct {
	EntityDeclarativeItems *[]EntityDeclarativeItem
	Node
}

type EntityDeclarativeItem interface {
}

type EntityStatementPart struct {
	EntityStatements *[]EntityStatement
	Node
}

type EntityStatement interface{}

type ArchitectureBody struct {
	Identifier                Identifier
	EntityName                EntityName
	ArchitectureStatementPart ArchitectureStatementPart
	ArchitectureSimpleName    ArchitectureSimpleName
}
