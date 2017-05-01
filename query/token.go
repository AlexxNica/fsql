package query

import (
	"fmt"
	"strings"
	"unicode"
)

// TokenType represents a Token's type.
type TokenType int8

const (
	// Unknown represents an unknown Token type.
	Unknown TokenType = iota
	// Select represents the SELECT clause.
	Select
	// From represents the FROM clause.
	From
	// Where represents the WHERE clause.
	Where
	// Or represents the OR condition concatenator (unimplemented).
	Or
	// And represents the AND condition concatenator (unimplemented).
	And
	// BeginsWith represents the BEGINSWITH comparator for string comparisons.
	BeginsWith
	// EndsWith represents the ENDSWITH comparator for string comparisons.
	EndsWith
	// Is represents the IS comparator for string comparisons.
	Is
	// Contains represents the CONTAINS comparator for string comparisons.
	Contains
	// Identifier represents the value for each Query.
	Identifier
	// OpenParen represents an open parenthesis.
	OpenParen
	// CloseParen represents a closed parenthesis.
	CloseParen
	// Comma represents a comma.
	Comma
	// Equals represents the `=` comparator for numeric comparisons.
	Equals
	// NotEquals represents the `<>` comparator for numeric comparisons.
	NotEquals
	// GreaterThanEquals represents the `>=` comparator for numeric comparisons.
	GreaterThanEquals
	// GreaterThan represents the `>` comparator for numeric comparisons.
	GreaterThan
	// LessThanEquals represents the `<=` comparator for numeric comparisons.
	LessThanEquals
	// LessThan represents the `<` comparator for numeric comparisons.
	LessThan
)

func (t TokenType) String() string {
	switch t {
	case Select:
		return "select"
	case From:
		return "from"
	case Where:
		return "where"
	case Or:
		return "or"
	case And:
		return "and"
	case BeginsWith:
		return "begins-with"
	case EndsWith:
		return "ends-with"
	case Is:
		return "is"
	case Contains:
		return "contains"
	case Identifier:
		return "identifier"
	case OpenParen:
		return "open-parentheses"
	case CloseParen:
		return "close-parentheses"
	case Comma:
		return "comma"
	case Equals:
		return "equal"
	case NotEquals:
		return "not-equal"
	case GreaterThanEquals:
		return "greater-than-or-equal"
	case GreaterThan:
		return "greater-than"
	case LessThanEquals:
		return "less-than-or-equal"
	case LessThan:
		return "less-than"
	default:
		return "unknown"
	}
}

// Token represents a single token.
type Token struct {
	Type TokenType
	Raw  string
}

func (t Token) String() string {
	return fmt.Sprintf("{type: %s, raw: \"%s\"}", t.Type.String(), t.Raw)
}

// Tokenizer represents a token worker.
type Tokenizer struct {
	input []rune
}

// NewTokenizer initializes a new Tokenizer.
func NewTokenizer(input string) *Tokenizer {
	return &Tokenizer{input: []rune(input)}
}

// All parses all tokens for this Tokenizer.
func (t *Tokenizer) All() []Token {
	tokens := []Token{}
	for tok := t.Next(); tok != nil; tok = t.Next() {
		tokens = append(tokens, *tok)
	}

	return tokens
}

// Next gets the next Token in this Tokenizer.
func (t *Tokenizer) Next() *Token {
	for {
		if !unicode.IsSpace(t.current()) {
			break
		}

		t.input = t.input[1:]
	}

	current := t.current()
	if current == -1 {
		return nil
	}

	switch current {
	case '(':
		t.input = t.input[1:]
		return &Token{Type: OpenParen, Raw: "("}

	case ')':
		t.input = t.input[1:]
		return &Token{Type: CloseParen, Raw: ")"}

	case ',':
		t.input = t.input[1:]
		return &Token{Type: Comma, Raw: ","}

	case '=':
		t.input = t.input[1:]
		return &Token{Type: Equals, Raw: "="}

	case '>':
		if t.peek() == '=' {
			t.input = t.input[2:]
			return &Token{Type: GreaterThanEquals, Raw: ">="}
		}

		t.input = t.input[1:]
		return &Token{Type: GreaterThan, Raw: ">"}

	case '<':
		if t.peek() == '=' {
			t.input = t.input[2:]
			return &Token{Type: LessThanEquals, Raw: ">="}
		}

		if t.peek() == '>' {
			t.input = t.input[2:]
			return &Token{Type: NotEquals, Raw: "<>"}
		}

		t.input = t.input[1:]
		return &Token{Type: LessThan, Raw: "<"}
	}

	if unicode.IsLetter(current) || unicode.IsDigit(current) ||
		current == '*' || current == '~' || current == '/' || current == '.' {
		word := t.readWord()

		tok := &Token{Raw: word}

		switch strings.ToUpper(word) {
		case "SELECT":
			tok.Type = Select
		case "FROM":
			tok.Type = From
		case "WHERE":
			tok.Type = Where
		case "OR":
			tok.Type = Or
		case "AND":
			tok.Type = And
		case "BEGINSWITH":
			tok.Type = BeginsWith
		case "ENDSWITH":
			tok.Type = EndsWith
		case "IS":
			tok.Type = Is
		case "CONTAINS":
			tok.Type = Contains
		default:
			tok.Type = Identifier
		}

		return tok
	}

	t.input = t.input[1:]
	return &Token{Type: Unknown, Raw: string([]rune{current})}
}

func (t *Tokenizer) current() rune {
	if len(t.input) == 0 {
		return -1
	}

	return t.input[0]
}

func (t *Tokenizer) peek() rune {
	if len(t.input) == 1 {
		return -1
	}

	return t.input[1]
}

func (t *Tokenizer) readWord() string {
	word := []rune{}

	for {
		r := t.current()

		if !(unicode.IsLetter(r) || r == '*' || r == '~' || r == '/' || r == '.' || unicode.IsDigit(r)) {
			return string(word)
		}

		word = append(word, r)
		t.input = t.input[1:]
	}
}