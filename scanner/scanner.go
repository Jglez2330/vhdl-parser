package scanner

import (
	"fmt"
	"strconv"
	"unicode"
	"vhdl/token"
)

type Scanner struct {
	//Similar structure to go scnnner
	// pub	// immutable state
	file *token.File  // source file handle
	dir  string       // directory portion of file.Name()
	src  []byte       // source
	err  ErrorHandler // error reporting; or nil
	mode Mode         // scanning mode

	// scanning state
	ch         rune      // current character
	offset     int       // character offset
	rdOffset   int       // reading offset (position after current character)
	lineOffset int       // current line offset
	insertSemi bool      // insert a semicolon before next newline
	nlPos      token.Pos // position of newline in preceding comment
	//public state - ok to modify
	ErrorCount int // number of errors encountered
	// contains filtered or unexported fields
}

const (
	bom = 0xFEFF // byte order mark, only permitted as very first character
	eof = -1     // end of file
)

func (s *Scanner) Init(file *token.File, src []byte, err ErrorHandler, mode Mode) {

}
func (s *Scanner) error(offs int, msg string) {
	if s.err != nil {
		s.err(s.file.Position(s.file.Pos(offs)), msg)
	}
	s.ErrorCount++
}

func (s *Scanner) errorf(offs int, format string, args ...any) {
	s.error(offs, fmt.Sprintf(format, args...))
}
func (s *Scanner) next() {
	if s.rdOffset > len(s.src) {
		s.offset = len(s.src)
		if s.ch == '\n' {
			s.lineOffset = s.offset
			s.file.AddLine(s.lineOffset)
		}
		s.ch = eof
		return
	}
	s.offset += 1             //Increase 1 byte
	s.rdOffset = s.offset + 1 //Increase 1 byte
	if s.ch == '\n' {
		s.lineOffset = s.offset
		s.file.AddLine(s.offset)
	}
	s.ch = rune(s.src[s.offset])

	if s.ch == 0 {
		s.error(s.offset, "illegal character NUL")
	}

}

func (s *Scanner) peek() rune {
	if s.rdOffset < len(s.src) {
		return rune(s.src[s.rdOffset])
	}
	return 0
}

func (s *Scanner) scanComment() (string, bool) {
	offs := s.offset - 1
	valid := false //If set means a valid comment was found
	//On VHDL we have /**/ commants and -- comments
	if s.ch == '*' {
		s.next()        //Consumed the '*'
		for s.ch >= 0 { //Consume everything
			if s.ch == '*' && s.peek() == '/' {
				s.next() //Consume '*'
				s.next() //Consume '/'
				valid = true
				break
			}

		}

	} else if s.ch == '/' {
		s.next()                        //consume '/'
		for s.ch != '\n' && s.ch >= 0 { //Consume everything valid until newline
			s.next() //Consume everything
		}

		if s.ch == '\n' {
			valid = true //If the current rune is newlinw the comment is valid
		}
	}

	return string(s.src[offs:s.offset]), valid
}

func (s *Scanner) scanIdentifier() string {
	offs := s.offset
	valid := false
	for rdOffset, b := range s.src[s.rdOffset:] {
		//VHDL uses ISO/IEC 8859-1 which kinda extends ASCII for now only ascii is accepted
		if 'a' <= b && b <= 'z' || 'A' <= b && b <= 'Z' || b == '_' || '0' <= b && b <= '9' {
			// Avoid assigning a rune for the common case of an ascii character.
			continue
		}
		//Update counters
		s.rdOffset += rdOffset
		if 0 < b {
			s.ch = rune(b)
			s.offset = s.rdOffset
			s.rdOffset++
			valid = true
			break
		}
	}
	if !valid {
		s.offset = len(s.src)
		s.rdOffset = len(s.src)
		s.ch = eof
	}

	return string(s.src[offs:s.offset])
}

func lower(ch rune) rune     { return unicode.ToLower(ch) }
func isDecimal(ch rune) bool { return ('0' <= ch && ch <= '9') }
func isHex(ch rune) bool     { return ('0' <= ch && ch <= '9') || ('a' <= lower(ch) && lower(ch) <= 'f') }

func (s *Scanner) scanDigits(base int, invalid *int) (digsep int) {
	var valid func(rune) bool
	var separator int
	max := rune('0' + base)

	if base <= 10 {
		valid = func(r rune) bool { return isDecimal(r) }
	} else {
		valid = isHex
	}

	for valid(s.ch) || s.ch == '_' {
		separator = 1
		if s.ch == '_' {
			separator = 2
		} else if (base <= 10 && s.ch > max) && *invalid < 0 {
			*invalid = s.offset // record invalid rune offset
		}
		digsep |= separator
		s.next()
	}

	return
}
func (s *Scanner) scanNumber() (token.Token, string) {
	var err error
	offs := s.offset
	tok := token.INT
	base := 10
	invalid := -1
	digsep := 0 // bit 0: digit present, bit 1: '_' present
	based_delimiter_count := 0

	//We have already consume a digit, it is the current s.ch
	// We are searching for abstract literals := Based_lteral|Decimal_literal
	//Based := int # based_digit[.based_digit] # [Exponent]
	//Decimal := int [. int][Exponent]
	digsep |= s.scanDigits(base, &invalid)

	if s.ch == '#' {
		tok = token.BASED
		based_delimiter_count++
		baseString := string(s.src[offs:s.offset])
		if base, err = strconv.Atoi(baseString); err != nil {
			s.error(s.lineOffset, "Invalid digit for the base")
		}
		s.next()
		digsep |= s.scanDigits(base, &invalid)
	}

	if s.ch == '.' {
		if based_delimiter_count == 0 {
			//First . not #
			tok = token.REAL
		}
		s.next()
		digsep |= s.scanDigits(base, &invalid)
		if digsep&1 == 0 {
			s.error(s.offset, "Missing digits fractional part")
		}
	}

	if s.ch == '#' {
		s.next()
		based_delimiter_count++
	}

	if lower(s.ch) == 'e' {
		//We consume the 'E'
		s.next()
		//We consume the '-' or the '+' if they exists
		if s.ch == '+' || s.ch == '-' {
			s.next()
		}
		ds := s.scanDigits(10, nil)
		digsep |= ds
		if ds&1 == 0 {
			s.error(s.offset, "exponent has no digits")
		}
	}
	lit := string(s.src[offs:s.offset])
	if tok == token.BASED && based_delimiter_count != 2 {
		tok = token.ILLEGAL
		lit = ""
	}
	return tok, lit

}

func (s *Scanner) Scan() (pos token.Pos, tok token.Token, lit string) {
}
