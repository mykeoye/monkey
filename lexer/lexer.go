package lexer

import (
	"fmt"
	"monkey/token"
)

type Lexer struct {
	input   string // the stream of text to tokennize
	pos     int    // the current position of the input cursor (points to the current char)
	readPos int    // the postion of the read cursor (points to the next character to read)
	ch      byte   // the character under view
}

// Creates a new instace of the Lexer
func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// Read a character from the current input, only advancing if possible
func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0 // the null character signifying there's nothing more to read
	} else {
		l.ch = l.input[l.readPos]
	}
	l.pos = l.readPos
	l.readPos += 1 // advance the read cursor
}

// Peeks into the input sequence to find the next char to consume
func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.input) {
		return 0 // there's no character to read EOF
	} else {
		return l.input[l.readPos]
	}
}

// Returns the next token that the Lexer has eaten
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.eatWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			tok = l.makeTwoCharToken(token.EQ)
		} else {
			tok = newToken(token.ASSIGNMENT, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '(':
		tok = newToken(token.LEFT_PAREN, l.ch)
	case ')':
		tok = newToken(token.RIGHT_PAREN, l.ch)
	case '{':
		tok = newToken(token.LEFT_BRACE, l.ch)
	case '}':
		tok = newToken(token.RIGHT_BRACE, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '!':
		if l.peekChar() == '=' {
			tok = l.makeTwoCharToken(token.NOT_EQ)
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '<':
		tok = newToken(token.LESS_THAN, l.ch)
	case '>':
		tok = newToken(token.GREATER_THAN, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		// Read a full string of characters, especially useful for identifiers of varying lengths
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdentifier(tok.Literal)
			return tok // readIdentifier already advances to the next token
			// Reads a full string of digits
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			fmt.Printf("Literal is %q", l.ch)
			tok.Literal = ""
			tok.Type = token.ILLEGAL
		}
	}

	// Advance the cursor to the next char if any
	l.readChar()

	return tok
}

// Reads an identifier by repeatedly consuming characters until it reaches
// a non-letter.
func (l *Lexer) readIdentifier() string {
	pos := l.pos
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

// Function checks a char to see if its a letter. Currently only supports ASCII
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// Reads a numeric character one digit at at a time
func (l *Lexer) readNumber() string {
	pos := l.pos
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

// Checks if a char is a digit
func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// Eats whitespace characters including tabs, newlines and carriage return. Stops
// when it encounters a non-whitespace character
func (l *Lexer) eatWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// Makes a two character token by consuming the next character and advancing the cursor
func (l *Lexer) makeTwoCharToken(tokenType token.TokenType) token.Token {
	// Get the current character read by the lexer
	ch := l.ch

	// Advance the cursor to get the next character
	l.readChar()

	// Combine both characters to form the token's literal
	literal := string(ch) + string(l.ch)

	return token.Token{Type: tokenType, Literal: literal}
}
