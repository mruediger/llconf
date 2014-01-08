package lexer

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/mruediger/llconf/compiler/token"
)

type Lexer struct {
	file       string
	input      string
	state      stateFn
	start      int
	pos        int
	width      int
	line       int
	tokens     chan token.Token
	parenDepth int
}

type stateFn func(*Lexer) stateFn

const (
	leftPromise  = "("
	rightPromise = ")"
)

const eof = -1

func Lex(file, input string) *Lexer {
	l := &Lexer{
		file: file,
		input: input,
		tokens: make(chan token.Token),
	}
	go l.run()
	return l
}

func (l *Lexer) NextToken() token.Token {
	token := <- l.tokens
	return token
}

func (l *Lexer) run() {
	for l.state = lexComment; l.state != nil; {
		l.state = l.state(l)
	}
}

func (l *Lexer) errorf(format string, args ...interface{}) stateFn {
	l.tokens <- token.Token{token.Error, token.Position{l.file, l.line, l.start, l.pos}, fmt.Sprintf(format, args...)}
	return nil
}

func (l *Lexer) emit(tt token.Type) {
	token := token.Token{tt, token.Position{l.file, l.line, l.start, l.pos}, l.input[l.start:l.pos]}
	token.Val = strings.TrimSpace(token.Val)
	l.tokens <- token
	l.start = l.pos
}

func (l *Lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		return eof
	}
	r,w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width

	return r
}

func (l *Lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *Lexer) backup() {
	l.pos -= l.width
}

func lexComment(l *Lexer) stateFn {
	switch r := l.next(); {
	case r == eof:
		l.emit(token.EOF)
		return nil;
	case r == '(':
		l.backup()
		return lexPromiseOpening
	default:
		l.emit(token.Comment)
	}

	return lexComment
}


func lexPromiseOpening(l *Lexer) stateFn {
	l.next()
	l.parenDepth++
	l.emit(token.LeftPromise)
	return lexPromiseName
}

func lexPromiseName(l *Lexer) stateFn {
	for {
		r := l.next()
		if r == eof {
			return l.errorf("unexpected eof in promise")
		}

		if !isValidPromiseNameRune(r) {
			l.backup()
			l.emit(token.PromiseName)
			return lexInsidePromise
		}
	}
	return nil
}

func lexInsidePromise(l *Lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == eof:
			return l.errorf("unexpected eof in promise")
		case r == '(':
			l.backup()
			return lexPromiseOpening
		case r == ')':
			l.backup()
			return lexPromiseClosing
		case r == '[':
			l.backup()
			return lexGetter
		case r == '"':
			l.backup()
			return lexArgument
		case unicode.IsSpace(r):
			// ignore
		default:
			return l.errorf("unexpected char inside promise: %#U",r)
		}
	}
	return nil
}

func lexPromiseClosing(l *Lexer) stateFn {
	for unicode.IsSpace(l.peek()) {
		l.next()
	}
	r := l.next()

	if r != ')' {
		return l.errorf("unexpected char at end of promise: %#U",r)
	}

	l.emit(token.RightPromise)
	l.parenDepth--

	if l.parenDepth == 0 {
		return lexComment
	} else {
		return lexInsidePromise
	}
}

func lexArgument(l *Lexer) stateFn {
	l.next()
	l.emit(token.LeftArg)

	for {
		switch r := l.next(); {
		case r == eof:
			return l.errorf("unexpected eof in argument")
		case r == '"':
			l.backup()
			l.emit(token.Argument)
			l.next()
			l.emit(token.RightArg)
			return lexInsidePromise
		}
	}
	return nil;
}

func lexGetter(l *Lexer) stateFn {
	l.next()
	l.emit(token.LeftGetter)
	return lexGetterType
}

func lexGetterType(l *Lexer) stateFn {
	for r := ' '; r != eof; r = l.next(){
		if r == ':' {
			l.backup()
			l.emit(token.GetterType)
			l.next()
			l.emit(token.GetterSeparator)
			return lexGetterValue
		}
	}
	return l.errorf("didn't found getter separator")
}

func lexGetterValue(l *Lexer) stateFn {
	for r := ' '; r != eof; r = l.next() {
		if r == ']' {
			l.backup()
			l.emit(token.GetterValue)
			l.next()
			l.emit(token.RightGetter)
			return lexInsidePromise
		}
	}

	return l.errorf("couldn't parse getter value")
}

func isValidPromiseNameRune(r rune) bool {
	return r == '-' || r == '_' || r == ' ' || unicode.IsLetter(r) || unicode.IsDigit(r)
}


func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
