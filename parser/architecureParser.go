package parser

import (
	"errors"
	"vhdl/ast"
	"vhdl/token"
)

func (p *Parser) parseArchitectureBody() (ast.ArchitectureBody, error) {
	architecture := ast.ArchitectureBody{}
	if p.trace {
		defer un(trace(p, "ArchitectureBody"))
	}
	architecture.Pos = p.pos

	if p.expect(token.ARCHITECTURE) == token.NoPos {
		return architecture, errors.New("Expected ARCHITECTURE keyword")
	}

	architecture.Identifier.Identifier = p.lit
	if p.expect(token.IDENT) == token.NoPos {
		return architecture, errors.New("Expected IDENTIFIER")
	}

	if p.expect(token.OF) == token.NoPos {
		return architecture, errors.New("Expected OF keyword")
	}

	architecture.EntityName = ast.SimpleName{Identifier: ast.Identifier{Identifier: p.lit}}
	if p.expect(token.IDENT) == token.NoPos {
		return architecture, errors.New("Expected IDENTIFIER")
	}

	if p.expect(token.IS) == token.NoPos {
		return architecture, errors.New("Expected IS keyword")
	}

	//parse architecture declarative part
	architecture_declarative_part, error := p.parseArchitectureDeclarativePart()
	if error != nil {
		return architecture, errors.New("Error parsing architecture declarative part")
	}
	architecture.ArchitectureDeclarativePart = architecture_declarative_part

	if p.expect(token.BEGIN) == token.NoPos {
		return architecture, errors.New("Expected BEGIN keyword")
	}

	//parse architecture statement part
	architecture_statement_part, error := p.parseArchitectureStatementPart()
	if error != nil {
		return architecture, errors.New("Error parsing architecture statement part")
	}
	architecture.ArchitectureStatementPart = architecture_statement_part

	if p.expect(token.END) == token.NoPos {
		return architecture, errors.New("Expected END keyword")
	}

	if p.tok == token.ARCHITECTURE {
		//Consume the architecture keyword
		p.next()
	}

	if p.tok == token.IDENT {
		//Consume the identifier
		architecture.ArchitectureSimpleName = &ast.SimpleName{Identifier: ast.Identifier{Identifier: p.lit}}
		p.next()
		if architecture.Identifier.Identifier != architecture.ArchitectureSimpleName.Identifier.Identifier {
			return architecture, errors.New("Architecture identifier does not match")
		}
	}

	if p.expect(token.SEMICOLON) == token.NoPos {
		return architecture, errors.New("Expected SEMICOLON")
	}

	return architecture, nil

}

func (p *Parser) parseArchitectureDeclarativePart() (ast.ArchitectureDeclarativePart, error) {
	architecture_declarative_part := ast.ArchitectureDeclarativePart{}
	if p.trace {
		defer un(trace(p, "ArchitectureDeclarativePart"))
	}
	//Consume until the begin keyword
	for p.tok != token.BEGIN {
		p.next()
	}
	return architecture_declarative_part, nil
}

func (p *Parser) parseArchitectureStatementPart() (ast.ArchitectureStatementPart, error) {
	architecture_statement_part := ast.ArchitectureStatementPart{}
	if p.trace {
		defer un(trace(p, "ArchitectureStatementPart"))
	}
	//Consume until the end keyword
	for p.tok != token.END {
		p.next()
	}
	return architecture_statement_part, nil
}
