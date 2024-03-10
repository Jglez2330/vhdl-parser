package token

type Token rune

const (
	ILLEGAL Token = iota
	EOF

	literal_begin
	IDENT   // main_main
	DEC     // 6.023E+24
	BASED   // 2#1111_1111# -> 255
	CHAR    // 'a'
	STRING  // "abc"
	BIT_STR // B"1111_1111_1111"
	COMMENT // --single /*multiple*/
	literal_end

	operator_beg
	//Operators and delimiters
	single_delimeter_beg
	CONCAT    // &
	APOS      // '
	LPAREN    // (
	RPAREN    // )
	MULT      // *
	PLUS      // +
	COMMA     // ,
	MINUS     // -
	DOT       // .
	DIV       // /
	COLON     // :
	SEMICOLON // ;
	LTH       // <
	EQL       // =
	GTH       // >
	BACKTICK  // `
	VLINE     // |
	LSQPAREN  // [
	RSQPAREN  // ]
	QUEST     // ?
	AT        // @
	single_delimeter_end
	compound_delimeter_beg
	ARROW      // =>
	EXP        //**
	VAR_ASSIGN // :=
	NEQ        // /=
	GEQ        // >=
	LEQ_SA     // <=
	BOX        // <>
	COND_CONV  // ??
	MEQ        // ?=
	MNEQ       // ?/=
	MLTH       // ?<
	MLEQ       // ?<=
	MGTH       // ?>
	MGEQ       // ?>=
	DOUBLE_LTH // <<
	DOUBLE_GTH // >>
	compound_delimeter_end
	operator_end

	oper_key_beg
	AND  //LAND
	OR   //LOR
	NAND //LNAND
	NOR  //LNOR
	XOR  //LXOR
	XNOR //LXNOR

	SLA //SLA
	SRL //SRL
	SLL //SLL
	SRA //SRA
	ROL //ROL
	ROR //ROR
	ABS //ABS
	NOT //NOT
	oper_key_end

	keyword_beg
	ACCESS        //ACCESS
	AFTER         //AFTER
	ALIAS         //ALIAS
	ALL           //ALL
	ARCHITECTURE  //ARCHITECTURE
	ARRAY         //ARRAY
	ASSERT        //ASSERT
	ASSUME        //ASSUME
	ATTRIBUTE     //ATTRIBUTE
	BEGIN         //BEGIN
	BLOCK         //BLOCK
	BODY          //BODY
	BUFFER        //BUFFER
	BUS           //BUS
	CASE          //CASE
	COMPONENT     //COMPONENT
	CONFIGURATION //CONFIGURATION
	CONSTANT      //CONSTANT
	CONTEXT       //CONTEXT
	COVER         //COVER
	DEFAULT       //DEFAULT
	DISCONNECT    //DISCONNECT
	DOWNTO        //DOWNTO
	ELSE          //ELSE
	ELSIF         //ELSIF
	END           //END
	ENTITY        //ENTITY
	EXIT          //EXIT
	FAIRNESS      //FAIRNESS
	FILE          //FILE
	FOR           //FOR
	FORCE         //FORCE
	FUNCTION      //FUNCTION
	GENERATE      //GENERATE
	GENERIC       //GENERIC
	GROUP         //GROUP
	GUARDED       //GUARDED
	IF            //IF
	IMPURE        //IMPURE
	IN            //IN
	INERTIAL      //INERTIAL
	INOUT         //INOUT
	IS            //IS
	LABEL         //LABEL
	LIBRARY       //LIBRARY
	LINKAGE       //LINKAGE
	LITERAL       //LITERAL
	LOOP          //LOOP
	MAP           //MAP
	MOD           //MOD
	NEW           //NEW
	NEXT          //NEXT
	NULL          //NULL
	OF            //OF
	ON            //ON
	OPEN          //OPEN
	OTHERS        //OTHERS
	OUT           //OUT
	PACKAGE       //PACKAGE
	PARAMETER     //PARAMETER
	PORT          //PORT
	POSTPONED     //POSTPONED
	PROCEDURE     //PROCEDURE
	PROCESS       //PROCESS
	PROPERTY      //PROPERTY
	PROTECTED     //PROTECTED
	PRIVATE       //PRIVATE
	PURE          //PURE
	RANGE         //RANGE
	RECORD        //RECORD
	REGISTER      //REGISTER
	REJECT        //REJECT
	RELEASE       //RELEASE
	REM           //REM
	REPORT        //REPORT
	RESTRICT      //RESTRICT
	RETURN        //RETURN
	SELECT        //SELECT
	SEQUENCE      //SEQUENCE
	SEVERITY      //SEVERITY
	SIGNAL        //SIGNAL
	SHARED        //SHARED
	STRONG        //STRONG
	SUBTYPE       //SUBTYPE
	THEN          //THEN
	TO            //TO
	TRANSPORT     //TRANSPORT
	TYPE          //TYPE
	UNAFFECTED    //UNAFFECTED
	UNITS         //UNITS
	UNTIL         //UNTIL
	USE           //USE
	VARIABLE      //VARIABLE
	VIEW          //VIEW
	VPKG          //VPKG
	VMODE         //VMODE
	VPROP         //VPROP
	VUNIT         //VUNIT
	WAIT          //WAIT
	WHEN          //WHEN
	WHILE         //WHILE
	WITH          //WITH
	//TODO: ADD PSL
	keyword_end
)

var tokens = [...]string{
	IDENT:   "IDENT",   // main_main
	DEC:     "DEC",     // 6.023E+24
	BASED:   "BASED",   // 2#1111_1111# -> 255
	CHAR:    "CHAR",    // 'a'
	STRING:  "STRING",  // "abc"
	BIT_STR: "BIT_STR", // B"1111_1111_1111"
	COMMENT: "COMMENT", // --single /*multiple*/

	CONCAT:     "&",   // &
	APOS:       "'",   // '
	LPAREN:     "(",   // (
	RPAREN:     ")",   // )
	MULT:       "*",   // *
	PLUS:       "+",   // +
	COMMA:      ",",   // ,
	MINUS:      "-",   // -
	DOT:        ".",   // .
	DIV:        "/",   // /
	COLON:      ":",   // :
	SEMICOLON:  ";",   // ;
	LTH:        "<",   // <
	EQL:        "=",   // =
	GTH:        ">",   // >
	BACKTICK:   "`",   // `
	VLINE:      "|",   // |
	LSQPAREN:   "[",   // [
	RSQPAREN:   "]",   // ]
	QUEST:      "?",   // ?
	AT:         "@",   // @
	ARROW:      "=>",  // =>
	EXP:        "**",  //**
	VAR_ASSIGN: ":=",  // :=
	NEQ:        "/=",  // /=
	GEQ:        ">=",  // >=
	LEQ_SA:     "<=",  // <=
	BOX:        "<>",  // <>
	COND_CONV:  "??",  // ??
	MEQ:        "?=",  // ?=
	MNEQ:       "?/=", // ?/=
	MLTH:       "?<",  // ?<
	MLEQ:       "?<=", // ?<=
	MGTH:       "?>",  // ?>
	MGEQ:       "?>=", // ?>=
	DOUBLE_LTH: "<<",  // <<
	DOUBLE_GTH: ">>",  // >>

	ABS:           "ABS",           //ABS
	ACCESS:        "ACCESS",        //ACCESS
	AFTER:         "AFTER",         //AFTER
	ALIAS:         "ALIAS",         //ALIAS
	ALL:           "ALL",           //ALL
	AND:           "AND",           //AND
	ARCHITECTURE:  "ARCHITECTURE",  //ARCHITECTURE
	ARRAY:         "ARRAY",         //ARRAY
	ASSERT:        "ASSERT",        //ASSERT
	ASSUME:        "ASSUME",        //ASSUME
	ATTRIBUTE:     "ATTRIBUTE",     //ATTRIBUTE
	BEGIN:         "BEGIN",         //BEGIN
	BLOCK:         "BLOCK",         //BLOCK
	BODY:          "BODY",          //BODY
	BUFFER:        "BUFFER",        //BUFFER
	BUS:           "BUS",           //BUS
	CASE:          "CASE",          //CASE
	COMPONENT:     "COMPONENT",     //COMPONENT
	CONFIGURATION: "CONFIGURATION", //CONFIGURATION
	CONSTANT:      "CONSTANT",      //CONSTANT
	CONTEXT:       "CONTEXT",       //CONTEXT
	COVER:         "COVER",         //COVER
	DEFAULT:       "DEFAULT",       //DEFAULT
	DISCONNECT:    "DISCONNECT",    //DISCONNECT
	DOWNTO:        "DOWNTO",        //DOWNTO
	ELSE:          "ELSE",          //ELSE
	ELSIF:         "ELSIF",         //ELSIF
	END:           "END",           //END
	ENTITY:        "ENTITY",        //ENTITY
	EXIT:          "EXIT",          //EXIT
	FAIRNESS:      "FAIRNESS",      //FAIRNESS
	FILE:          "FILE",          //FILE
	FOR:           "FOR",           //FOR
	FORCE:         "FORCE",         //FORCE
	FUNCTION:      "FUNCTION",      //FUNCTION
	GENERATE:      "GENERATE",      //GENERATE
	GENERIC:       "GENERIC",       //GENERIC
	GROUP:         "GROUP",         //GROUP
	GUARDED:       "GUARDED",       //GUARDED
	IF:            "IF",            //IF
	IMPURE:        "IMPURE",        //IMPURE
	IN:            "IN",            //IN
	INERTIAL:      "INERTIAL",      //INERTIAL
	INOUT:         "INOUT",         //INOUT
	IS:            "IS",            //IS
	LABEL:         "LABEL",         //LABEL
	LIBRARY:       "LIBRARY",       //LIBRARY
	LINKAGE:       "LINKAGE",       //LINKAGE
	LITERAL:       "LITERAL",       //LITERAL
	LOOP:          "LOOP",          //LOOP
	MAP:           "MAP",           //MAP
	MOD:           "MOD",           //MOD
	NAND:          "NAND",          //NAND
	NEW:           "NEW",           //NEW
	NEXT:          "NEXT",          //NEXT
	NOR:           "NOR",           //NOR
	NOT:           "NOT",           //NOT
	NULL:          "NULL",          //NULL
	OF:            "OF",            //OF
	ON:            "ON",            //ON
	OPEN:          "OPEN",          //OPEN
	OR:            "OR",            //OR
	OTHERS:        "OTHERS",        //OTHERS
	OUT:           "OUT",           //OUT
	PACKAGE:       "PACKAGE",       //PACKAGE
	PARAMETER:     "PARAMETER",     //PARAMETER
	PORT:          "PORT",          //PORT
	POSTPONED:     "POSTPONED",     //POSTPONED
	PROCEDURE:     "PROCEDURE",     //PROCEDURE
	PROCESS:       "PROCESS",       //PROCESS
	PROPERTY:      "PROPERTY",      //PROPERTY
	PROTECTED:     "PROTECTED",     //PROTECTED
	PRIVATE:       "PRIVATE",       //PRIVATE
	PURE:          "PURE",          //PURE
	RANGE:         "RANGE",         //RANGE
	RECORD:        "RECORD",        //RECORD
	REGISTER:      "REGISTER",      //REGISTER
	REJECT:        "REJECT",        //REJECT
	RELEASE:       "RELEASE",       //RELEASE
	REM:           "REM",           //REM
	REPORT:        "REPORT",        //REPORT
	RESTRICT:      "RESTRICT",      //RESTRICT
	RETURN:        "RETURN",        //RETURN
	ROL:           "ROL",           //ROL
	ROR:           "ROR",           //ROR
	SELECT:        "SELECT",        //SELECT
	SEQUENCE:      "SEQUENCE",      //SEQUENCE
	SEVERITY:      "SEVERITY",      //SEVERITY
	SIGNAL:        "SIGNAL",        //SIGNAL
	SHARED:        "SHARED",        //SHARED
	SLA:           "SLA",           //SLA
	SLL:           "SLL",           //SLL
	SRA:           "SRA",           //SRA
	SRL:           "SRL",           //SRL
	STRONG:        "STRONG",        //STRONG
	SUBTYPE:       "SUBTYPE",       //SUBTYPE
	THEN:          "THEN",          //THEN
	TO:            "TO",            //TO
	TRANSPORT:     "TRANSPORT",     //TRANSPORT
	TYPE:          "TYPE",          //TYPE
	UNAFFECTED:    "UNAFFECTED",    //UNAFFECTED
	UNITS:         "UNITS",         //UNITS
	UNTIL:         "UNTIL",         //UNTIL
	USE:           "USE",           //USE
	VARIABLE:      "VARIABLE",      //VARIABLE
	VIEW:          "VIEW",          //VIEW
	VPKG:          "VPKG",          //VPKG
	VMODE:         "VMODE",         //VMODE
	VPROP:         "VPROP",         //VPROP
	VUNIT:         "VUNIT",         //VUNIT
	WAIT:          "WAIT",          //WAIT
	WHEN:          "WHEN",          //WHEN
	WHILE:         "WHILE",         //WHILE
	WITH:          "WITH",          //WITH
	XNOR:          "XNOR",          //XNOR
	XOR:           "XOR",           //XOR

}

var keywords map[string]Token

func init() {
	keywords_size := (keyword_end - keyword_beg + 1) + (oper_key_end - oper_key_beg + 1)
	keywords = make(map[string]Token, keywords_size)
	for i := oper_key_beg + 1; i < oper_key_end; i++ {
		keywords[tokens[i]] = i
	}
	for i := keyword_beg + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
}

func Lookup(ident string) Token {
	tok, is_keyword := keywords[ident]
	if !is_keyword {
		return IDENT
	}
	return tok
}

func (tok Token) IsKeyword() bool {
	return oper_key_beg < tok && tok < oper_key_end || keyword_beg < tok && tok < keyword_end
}

func (tok Token) IsLiteral() bool {
	return literal_begin < tok && tok < literal_end
}

func (tok Token) IsOperator() bool {
	return oper_key_beg < tok && tok < oper_key_end
}

const (
	LowestPrecedence  = 0
	HighestPrecedence = 9
)

func (op Token) Precedence() int {
	switch op {
	case ABS, NOT:
		return 9
	case EXP:
		return 8
	case MULT, DIV, MOD, REM:
		return 7
	case PLUS, MINUS:
		return 6
	case CONCAT:
		return 5
	case SLL, SRL, SLA, SRA, ROL, ROR:
		return 4
	case EQL, NEQ, LTH, LEQ_SA, GTH, GEQ, MEQ, MNEQ, MLTH, MLEQ, MGTH, MGEQ:
		return 3
	case AND, OR, NAND, NOR, XOR, XNOR:
		return 2
	case COND_CONV:
		return 1
	default:
		return LowestPrecedence
	}
}

func (op Token) PrecedenceAlternativeOperator() int {
	switch op {
	//Unary operator
	case AND, OR, NAND, NOR, XOR, XNOR:
		return 9
	//Addition
	case PLUS, MINUS, CONCAT:
		return 5
	default:
		return op.Precedence()
	}
}

func (tok Token) String() string{
    return tokens[tok]
}
