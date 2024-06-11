package scanner

//TODO make scanner more readle

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
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
	ch         rune // current character
	offset     int  // character offset
	rdOffset   int  // reading offset (position after current character)
	lineOffset int  // current line offset
	//public state - ok to modify
	ErrorCount int // number of errors encountered
	// contains filtered or unexported fields
}

const (
	bom = 0xFEFF // byte order mark, only permitted as very first character
	eof = -1     // end of file
)

func (s *Scanner) Init(file *token.File, src []byte, err ErrorHandler, mode Mode) {
	if file.Size() != len(src) {
		panic(fmt.Sprintf("file size (%d) does not match src len (%d)", file.Size(), len(src)))
	}
	s.file = file
	s.dir, _ = filepath.Split(file.Name())
	s.src = src
	s.err = err
	s.mode = mode

	s.ch = ' '
	s.offset = 0
	s.rdOffset = 0
	s.lineOffset = 0
	s.ErrorCount = 0

	s.next()
	if s.ch == bom {
		s.next() // ignore BOM at file beginning
	}
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
	if s.rdOffset >= len(s.src) {
		s.offset = len(s.src)
		if s.ch == '\n' {
			s.lineOffset = s.offset
			s.file.AddLine(s.lineOffset)
		}
		s.ch = eof
		return
	}
	if s.ch == '\n' {
		s.lineOffset = s.offset
		s.file.AddLine(s.offset)
	}
	s.ch = rune(s.src[s.rdOffset])

	s.offset = s.rdOffset //Increase 1 byte
	s.rdOffset += 1       //Increase 1 byte
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

func (s *Scanner) scanMultipleLineComment() (string, bool) {
	offs := s.offset - 1
	valid := false //If set means a valid comment was found
	//On VHDL we have /**/ commants and -- comments
	if s.ch == '*' {
		//Consumed the '*'
		for s.ch >= 0 { //Consume everything
			s.next()
			if s.ch == '*' && s.peek() == '/' {
				s.next() //Consume '*'
				s.next() //Consume '/'
				valid = true
				break
			}

		}

	}
	return string(s.src[offs:s.offset]), valid
}

func (s *Scanner) scanSingleLineComment() (string, bool) {
	offs := s.offset - 1
	valid := false //If set means a valid comment was found
	//On VHDL we have /**/ commants and -- comments
	if s.ch == '-' {
		s.next()                        //consume '-'
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
		}
		break
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

func (s *Scanner) scanRune() string {
	offs := s.offset - 1 // Consumed "'"
	//VHDL char is 'graphic_character'
	if s.ch == '\n' {
		s.error(s.offset, "Breaking space character")
	}
	s.next() //Consume graphic_character
	s.next() //Consume '
	return string(s.src[offs:s.offset])
}

func (s *Scanner) scanString() string {
	offs := s.offset - 1 //Already consumed "
	for {
		if s.ch == '\n' {
			s.error(offs, "string literal not terminated")
			break
		}
		s.next()
		if s.ch == '"' {
			//Consume '"'
			s.next()
			break
		}
	}
	return string(s.src[offs:s.offset])

}

func (s *Scanner) skipWhitespace() {
	for s.ch == ' ' || s.ch == '\t' || s.ch == '\n' || s.ch == '\r' {
		s.next()
	}
}

func isBaseSpecifierPrefix(ch rune) bool {
	//B|O|X|UB|UO|UX|SB|SO|SX|D
	return lower(ch) == 'b' || lower(ch) == 'o' || lower(ch) == 'u' || lower(ch) == 'x' || lower(ch) == 's' || lower(ch) == 'd'
}

func (s *Scanner) scanBitStringPrefix() bool {
	is_valid_bit_string_prefix := true
	has_sign := false

	if lower(s.ch) == 'u' || lower(s.ch) == 's' {
		//Consume 'u' or 's'
		s.next()
		has_sign = true
	}

	if lower(s.ch) == 'd' && has_sign {
		is_valid_bit_string_prefix = false
	}

	//Consume X. O, B or D
	s.next()

	return is_valid_bit_string_prefix
}

func (s *Scanner) Scan() (pos token.Pos, tok token.Token, lit string) {

	s.skipWhitespace()
	pos = s.file.Pos(s.offset)

	switch ch := s.ch; {

	//--------------------------------------Letter-------------------------------------------------------
	//We detected a graphic_character it could be a identifier or a keyword or a bit strign
	case unicode.IsLetter(ch):
		lit = s.scanIdentifier()
		//If the len is greater than 1, we have a keyword, identifier or bit_string so we are going to match the literal to the token to see
		//if it is a keyword
		if len(lit) > 1 {
			tok = token.Lookup(lit)
			switch tok {
			case token.IDENT:
				//Ugly but needed to capture bit_str with number length
				switch strings.ToLower(lit) {
				case "sb", "so", "sx", "ub", "uo", "ux":
					if s.ch == '"' {
						s.next()
						bit_str := s.scanString()
						lit = lit + bit_str
						tok = token.BIT_STR
					}

				}
			}
		} else {
			tok = token.IDENT
			switch strings.ToLower(lit) {
			case "x", "o", "b", "d":
				if s.ch == '"' {
					s.next()
					bit_str := s.scanString()
					lit = lit + bit_str
					tok = token.BIT_STR
				}

			}
		}
	//---------------------------------------------------------------------------------------------

	//------------------Digits---------------------------------------------------------------------------
	case isDecimal(ch):
		// We found a string literal or an abstract literal
		tok, lit = s.scanNumber()
		//Right now we have a abstract literal to check if it as string we need the following
		// that the lit is a integer, not float or based. It is followed by: B|O|X|UB|UO|UX|SB|SO|SX|D
		if tok == token.INT && isBaseSpecifierPrefix(s.ch) {
			offs := s.offset
			if s.scanBitStringPrefix() {
				s.next()
				s.scanString()
				tok = token.BIT_STR
				lit = lit + string(s.src[offs:s.offset])
			} else {
				tok = token.ILLEGAL
				lit = lit + string(s.src[offs:s.offset])
			}
		}

		//---------------------------------------------------------------------------------------------
		//-------------Default Other symbols--------------------------------------------------------------------------------
	default:
		s.next()
		switch ch {
		case eof:
			tok = token.EOF
		case '\n':
			// we only reach here if s.insertSemi was
			// set in the first place and exited early
			// from s.skipWhitespace()
			return pos, token.ILLEGAL, "\n"
		case '"':
			tok = token.STRING
			lit = s.scanString()
		case '\'':
			//TODO add only ' recognition
			//peek must be ' only 1 char inside ' '
			if s.peek() != '\'' {
				tok = token.APOS
			} else {
				tok = token.CHAR
				lit = s.scanRune()
			}
		case ':':
			tok = token.COLON
			if s.ch == '=' {
				s.next()
				tok = token.VAR_ASSIGN
			}
		case '.':
			tok = token.DOT
		case ',':
			tok = token.COMMA
		case ';':
			tok = token.SEMICOLON
			lit = ";"
		case '(':
			tok = token.LPAREN
		case ')':
			tok = token.RPAREN
		case '[':
			tok = token.LSQPAREN
		case ']':
			tok = token.RSQPAREN
		case '+':
			tok = token.PLUS
		case '-':
			tok = token.MINUS
			if s.ch == '-' {
				//Comment
				comment, valid := s.scanSingleLineComment()
				tok = token.COMMENT
				if !valid {
					tok = token.ILLEGAL
				}
				lit = comment
			}

		case '&':
			tok = token.CONCAT
		case '?':
			tok = token.COND_CONV
		case '=':
			tok = token.EQL
			if s.ch == '>' {
				tok = token.ARROW
				s.next()
			}
		case '*':
			tok = token.MULT
			if s.ch == '*' {
				tok = token.EXP
				s.next()
			}
		case '/':
			tok = token.DIV
			if s.ch == '*' {
				//Comment
				comment, valid := s.scanMultipleLineComment()
				tok = token.COMMENT
				if !valid {
					tok = token.ILLEGAL
				}
				lit = comment
			} else if s.ch == '=' {
				tok = token.NEQ
				s.next()
			}
		case '<':
			tok = token.LTH
			if s.ch == '=' {
				tok = token.LEQ_SA
				s.next()
			} else if s.ch == '>' {
				tok = token.BOX
				s.next()
			}
		case '>':
			tok = token.GTH
			if s.ch == '=' {
				tok = token.GEQ
				s.next()
			}
		case '|':
			tok = token.VLINE
			//TODO: Do rest of ? operators/delimeters

		default:
			// next reports unexpected BOMs - don't repeat
			if ch != bom {
				// Report an informative error for U+201[CD] quotation
				// marks, which are easily introduced via copy and paste.
				if ch == '“' || ch == '”' {
					s.errorf(s.file.Offset(pos), "curly quotation mark %q (use neutral %q)", ch, '"')
				} else {
					s.errorf(s.file.Offset(pos), "illegal character %#U", ch)
				}
			}
			tok = token.ILLEGAL
			lit = string(ch)

		}

	}

	return
}
