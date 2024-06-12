package parser

import (
	"errors"
	"vhdl/ast"
	"vhdl/token"
)

func (p *Parser) parseEntityDeclaration() (ast.EntityDeclaration, error) {
	var entity ast.EntityDeclaration
	if p.trace {
		defer un(trace(p, "EntityDeclaration"))
	}
	entity.Pos = p.pos
	if p.expect(token.ENTITY) == token.NoPos {
		return entity, errors.New("Expected ENTITY keyword")
	}
	entity.Identifier.Identifier = p.lit
	if p.expect(token.IDENT) == token.NoPos {
		return entity, errors.New("Expected IDENTIFIER")
	}
	if p.expect(token.IS) == token.NoPos {
		return entity, errors.New("Expected IS keyword")
	}

	entity_header, error := p.parseEntityHeader()
	if error != nil {
		return entity, errors.New("Error parsing entity header")
	}
	entity.EntityHeader = entity_header

	entity_declaritve_part, error := p.parseEntityDeclarativePart()
	if error != nil {
		return entity, errors.New("Error parsing entity declarative part")
	}
	entity.EntityDeclarativePart = entity_declaritve_part

	if p.tok == token.BEGIN {
		//Consume the begin keyword
		p.next()
		//Parse the entity statement part
		entity_statements, error := p.parseEntityStatementPart()
		entity.EntityStatementPart = &ast.EntityStatementPart{EntityStatements: &entity_statements}
		if error != nil {
			return entity, errors.New("Error parsing entity statement part")
		}
	}

	if p.expect(token.END) == token.NoPos {
		return entity, errors.New("Expected END keyword")
	}

	if p.tok == token.ENTITY {
		//Consume the entity keyword
		p.next()
	}

	if p.tok == token.IDENT {
		//Consume the identifier
		entity.EntitySimpleName = &ast.SimpleName{Identifier: ast.Identifier{Identifier: p.lit}}
		if entity.Identifier.Identifier != entity.EntitySimpleName.Identifier.Identifier {
			p.errorExpected(p.pos, "Expected %s, found %s", entity.Identifier.Identifier, entity.EntitySimpleName.Identifier.Identifier)
		}
		p.next()
	}

	if p.expect(token.SEMICOLON) == token.NoPos {
		return entity, errors.New("Expected SEMICOLON")
	}

	return entity, nil
}

func (p *Parser) parseEntityHeader() (ast.EntityHeader, error) {
	var entityHeader ast.EntityHeader
	if p.trace {
		defer un(trace(p, "EntityHeader"))
	}
	//TODO add later
	//If there is a generic clause or port clause consume it

	if p.tok == token.GENERIC {
		for p.tok != token.SEMICOLON {
			p.next()
		}
	}

	if p.tok == token.PORT {
		for p.tok != token.SEMICOLON {
			p.next()
		}
	}

	return entityHeader, nil
}

func (p *Parser) parseEntityDeclarativePart() (ast.EntityDeclarativePart, error) {
	//TODO parse entity declarative part
	//Right now we are going to to skip until the begin or end keyword
	for p.tok != token.BEGIN && p.tok != token.END {
		p.next()
	}

	return ast.EntityDeclarativePart{}, nil
}

func (p *Parser) parseEntityStatementPart() ([]ast.EntityStatement, error) {
	//TODO parse entity statement part
	//Right now we are going to to skip until the end keyword
	for p.tok != token.END {
		p.next()
	}

	return []ast.EntityStatement{}, nil

}
