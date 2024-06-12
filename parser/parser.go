package parser

import (
	"errors"
	"fmt"
	"vhdl/ast"
	"vhdl/scanner"
	"vhdl/token"
)

type Mode uint

const (
	PackageClauseOnly    Mode             = 1 << iota // stop parsing after package clause
	ImportsOnly                                       // stop parsing after import declarations
	ParseComments                                     // parse comments and add them to AST
	Trace                                             // print a trace of parsed productions
	DeclarationErrors                                 // report declaration errors
	SpuriousErrors                                    // same as AllErrors, for backward-compatibility
	SkipObjectResolution                              // skip deprecated identifier resolution; see ParseFile
	AllErrors            = SpuriousErrors             // report all errors (not just the first 10 on different lines)
)

// The parser structure holds the parser's internal state.
type Parser struct {
	file    *token.File
	errors  scanner.ErrorList
	scanner scanner.Scanner
	mode    Mode
	trace   bool
	indent  int // indentation used for tracing output

	// Next token
	pos token.Pos   // token position
	tok token.Token // one token look-ahead
	lit string      // token literal

	exprLev int  // < 0: in control clause, >= 0: in expression
	inRhs   bool // if set, the parser is parsing a rhs expression

	// nestLev is used to track and limit the recursion depth
	// during parsing.
	nestLev int
}

func (p *Parser) Init(fset *token.FileSet, filename string, src []byte, mode Mode) {
	p.file = fset.AddFile(filename, -1, len(src))
	error_handler := func(pos token.Position, msg string) { p.errors.Add(pos, msg) }
	p.scanner.Init(p.file, src, error_handler, scanner.ScanComments)

	p.mode = mode
	p.trace = mode&Trace != 0 // for convenience (p.trace is used frequently)
	p.next()
}

// next advances to the next token.
// ----------------------------------------------------------------------------
// Parsing support

func (p *Parser) printTrace(a ...any) {
	const dots = ". . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . "
	const n = len(dots)
	pos := p.file.Position(p.pos)
	fmt.Printf("%5d:%3d: ", pos.Line, pos.Column)
	i := 2 * p.indent
	for i > n {
		fmt.Print(dots)
		i -= n
	}
	// i <= n
	fmt.Print(dots[0:i])
	fmt.Println(a...)
}

func trace(p *Parser, msg string) *Parser {
	p.printTrace(msg, "(")
	p.indent++
	return p
}

// Usage pattern: defer un(trace(p, "..."))
func un(p *Parser) {
	p.indent--
	p.printTrace(")")
}

// maxNestLev is the deepest we're willing to recurse during parsing
const maxNestLev int = 1e5

func incNestLev(p *Parser) *Parser {
	p.nestLev++
	if p.nestLev > maxNestLev {
		p.error(p.pos, "exceeded max nesting depth")
		panic(bailout{})
	}
	return p
}

type bailout struct {
	pos token.Pos
	msg string
}

// decNestLev is used to track nesting depth during parsing to prevent stack exhaustion.
// It is used along with incNestLev in a similar fashion to how un and trace are used.
func decNestLev(p *Parser) {
	p.nestLev--
}

func (p *Parser) next0() {
	if p.trace && p.pos.IsValid() {
		tok_value := p.tok.String()
		if p.tok.IsLiteral() {
			p.printTrace(tok_value, p.lit)
		} else if p.tok.IsKeyword() {
			p.printTrace("\"" + tok_value + "\"")
		} else {
			p.printTrace(tok_value)
		}
	}

	for {
		p.pos, p.tok, p.lit = p.scanner.Scan()
		if p.tok == token.COMMENT {
			if p.mode&ParseComments == 0 {
				continue
			}
		}
		break

	}
}

// Consume a comment and return it and the line on which it ends.
func (p *Parser) consumeComment() (comment *ast.Comment, endline int) {
	// /*-style comments may end on a different line than where they start.
	// Scan the comment for '\n' chars and adjust endline accordingly.
	endline = p.file.Line(p.pos)
	if p.lit[1] == '*' {
		// don't use range here - no need to decode Unicode code points
		for i := 0; i < len(p.lit); i++ {
			if p.lit[i] == '\n' {
				endline++
			}
		}
	}

	comment = &ast.Comment{Slash: p.pos, Text: p.lit}
	p.next0()

	return
}

func (p *Parser) next() {
	// prev := p.pos
	p.next0()

	if p.tok == token.COMMENT {
		//TODO parse comment
	}

}

func (p *Parser) ParseFile() (ast.File, error) {
	file := ast.File{}

	var designUnits []ast.DesignUnit
	for p.tok != token.EOF {
		if designUnit, error := p.ParseDesignUnit(); error != nil {
			designUnits = append(designUnits, designUnit)
		}
	}
	file.DesignUnits = designUnits
	return file, nil
}

func (p *Parser) ParseDesignUnit() (ast.DesignUnit, error) {
	var DesignUnit ast.DesignUnit
	//TODO parse context clause

	//Parse library unit
	if lu, error := p.parseLibraryUnit(); error != nil {
		DesignUnit.LibraryUnit = &lu
	}
	//if lu = p.parseLibraryUnit(); lu == nil {

	return DesignUnit, nil
}

func (p *Parser) parseLibraryUnit() (ast.LibraryUnit, error) {
	if p.trace {
		defer un(trace(p, "LibraryUnit"))
	}

    if p.isPrimaryUnit(p.tok) {
        return p.parsePrimaryUnit()
    }else if p.isSecondaryUnit(p.tok) {
        return p.parseSecondaryUnit()
    }else {
        p.errorExpected(p.pos, "expected primary or secondary unit, found %s", p.tok)
    }




	return ast.EntityDeclaration{}, nil

}


func (p *Parser) parsePrimaryUnit() (ast.PrimaryUnit, error) {

	switch p.tok {
	case token.ENTITY:
		entity, error := p.parseEntityDeclaration()
		if error != nil {
			return entity, errors.New("invalid entity declaration")
		}
		return entity, nil
	default:
		p.errorExpected(p.pos, "expected entity declaration, found %s", p.tok)
	}
	return ast.EntityDeclaration{}, nil

}

func (p* Parser) isPrimaryUnit(tok token.Token) bool {
    //token is primary unit if it is entity, package, configuration, package instatioation,context
    return tok == token.ENTITY || tok == token.PACKAGE || tok == token.CONFIGURATION || tok == token.CONTEXT
}

func (p* Parser) isSecondaryUnit(tok token.Token) bool {
    //token is secondary unit if it is architecture body or package body
    return tok == token.ARCHITECTURE 
}

func (p* Parser) parseSecondaryUnit() (ast.SecondaryUnit, error) {
    return ast.ArchitectureBody{}, nil
}
    
func (p *Parser) expect(tok token.Token) token.Pos {
	pos := token.NoPos
	if p.tok != tok {
		p.errorExpected(p.pos, "expected %s, found %s", tok, p.tok)
	} else {
		pos = p.pos
		p.next()
	}
	return pos
}

func (p *Parser) errorExpected(pos token.Pos, format string, args ...interface{}) {
	p.error(pos, format, args...)
}

func (p *Parser) error(pos token.Pos, format string, args ...interface{}) {
	p.errors.Add(p.file.Position(pos), fmt.Sprintf(format, args...))
}
