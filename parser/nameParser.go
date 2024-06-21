package parser

import (
	"errors"
	"runtime/trace"
	"vhdl/ast"
	"vhdl/token"
)

func (p *Parser) parseName() (ast.Name, error) {
	var name any
	if p.trace {
		defer un(trace(p, "Name"))
	}
    switch p.tok {
    case token.IDENT:
        p.parseNamePart()
        //Check if is selected_name, simple_name, indexed_name, slice_name, attribute_name
    case token.STRING:
        //Is operator_symbol
    case token.CHAR:
        //Is character_literal
    case default:
        return name, errors.New("invalid name")
    }

	return name, nil
}

func (p *Parser) parseNamePart() (ast.NamePart, error) {
    if p.trace {
        defer un(trace(p, "NamePart"))
    }
    var name any 
    switch p.tok {
    case token.DOT:
    case Token.LPAREN:
        p.parseFunctionCallOrIndexedNameOrSliceName()
    case token.APOS:
        p.parseAttributeName()
    default:
        return name, errors.New("Invalid namepart")
    }
}

func (p *Parser) parseSimpleName() (ast.SimpleName, error) {
	var simple_name ast.SimpleName
	if p.trace {
		defer un(trace(p, "SimpleName"))
	}

	identifier := ast.Identifier{Identifier: p.lit}
	simple_name = ast.SimpleName{Identifier: identifier}

	if p.expect(token.IDENT) == token.NoPos {
		return simple_name, errors.New("invalid identifier")
	}
	return simple_name, nil

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
	selected_name.Prefix = prefix

	if p.expect(token.DOT) == token.NoPos {
		return selected_name, errors.New("invalid selected name")
	}

	suffix, error := p.parseSuffix()
	if error != nil {
		return selected_name, errors.New("invalid suffix")
	}
	selected_name.Suffix = suffix

	return selected_name, nil
}

func (p *Parser) parsePrefix() (ast.Prefix, error) {
	var prefix ast.Prefix
	if p.trace {
		defer un(trace(p, "Prefix"))
	}
	name, error := p.recursiveParseName(nil)
	if error != nil {
		return prefix, errors.New("invalid name")
	}
	prefix = name

	return prefix, nil
}

func (p *Parser) recursiveParseName(name_option any) (ast.Name, error) {
	var name any
	if p.trace {
		defer un(trace(p, "Name"))
	}

	//Check if is selected_name
	if p.tok == token.IDENT && p.tok2 == token.DOT {
		//Call recursivePrefix function
		var selected_name ast.SelectedName
		if name_option == nil {
			name_option = ast.SelectedName{}
			//name_option.Prefix = p.lit
		}
		selected_name.Prefix = name_option
		p.next()
	}
	return name, nil

}

func (p *Parser) parseSuffix() (ast.Suffix, error) {
	var suffix ast.Suffix
	if p.trace {
		defer un(trace(p, "Suffix"))
	}
	switch p.tok {
	case token.IDENT:
		suffix = ast.SimpleName{Identifier: ast.Identifier{Identifier: p.lit}}
	case token.STRING:
		operator_string := ast.OperatorSymbol{Symbol: p.lit}
		suffix = operator_string
	case token.CHAR:
		character_literal := ast.CharacterLiteral{GraphicCharacter: ast.GraphicCharacter{Character: p.lit}}
		suffix = character_literal
	case token.ALL:
		all := ast.Keyword{Token: p.tok, Value: p.lit}
		suffix = all
	default:
		return suffix, errors.New("invalid suffix")
	}
	p.next()

	return suffix, nil
}
