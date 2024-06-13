package parser

import (
	"errors"
	"vhdl/ast"
	"vhdl/token"
)


func (p *Parser) parseName() (ast.Name, error) {
    var name ast.Name
    if p.trace {
        defer un(trace(p, "Name"))
    }
    //All the names start with an identifier
    if p.tok != token.IDENT {
        return name, errors.New("invalid name")
    }
    name_start := p.lit
    p.next()
    //If the name is a selected name



    return name, nil
}


func (p *Parser) parseSelectedName() (ast.SelectedName, error) {
    var selected_name ast.SelectedName
    if p.trace {
        defer un(trace(p, "SelectedName"))
    }

    prefix, error := p.parsePrefix()
    if error != nil {
        return selected_name, errors.New("invalid prefix")
    }

    if p.expect(token.DOT) == token.NoPos {
        return selected_name, errors.New("invalid selected name")
    }

    suffix, error := p.parseSuffix()
    if error != nil {
        return selected_name, errors.New("invalid suffix")
    }

    return selected_name, nil
}

func (p *Parser) parsePrefix() (ast.Prefix, error) {
    var prefix ast.Prefix
    if p.trace {
        defer un(trace(p, "Prefix"))
    }

    name, error := p.parseName()
    if error != nil {
        return prefix, errors.New("invalid name")
    }
    prefix = name

    return prefix, nil
}

