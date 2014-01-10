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
	tokens     chan token.Token
	parenDepth int
	getterDepth int
}

type stateFn func(*Lexer) stateFn

const (
	leftPromise  = "("
	rightPromise = ")"
)

const eof = -1

func (l *Lexer) run() {
	for l.state = lexComment; l.state != nil; {
		l.state = l.state(l)
	}
}

func (l *Lexer) errorf(format string, args ...interface{}) stateFn {
	l.tokens <- token.Token{token.Error, token.Position{
		l.file,
		l.line(),
		l.start,
		l.pos}, fmt.Sprintf(format, args...)}
	return nil
}

func (l *Lexer) emit(tt token.Type) {
	token := token.Token{tt, token.Position{
		l.file,
		l.line(),
		l.start,
		l.pos}, l.input[l.start:l.pos]}
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

func (l *Lexer) backup() {
	l.pos -= l.width
}

// this solution is easier than handle next+backup
func (l *Lexer) line() int {
	return 1 + strings.Count(l.input[:l.pos], "\n")
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
			return lexInsideGetter
		case r == '"':
			l.backup()
			return lexArgument
		case unicode.IsSpace(r):
			// ignore
		default:
			return l.errorf("unexpected char inside promise: %q",r)
		}
	}
	return nil
}

func lexPromiseClosing(l *Lexer) stateFn {
	r := l.next()
	for unicode.IsSpace(r) {
		r = l.next()
	}


	if r != ')' {
		return l.errorf("unexpected char at end of promise: %q",r)
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

func lexInsideGetter(l *Lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == eof:
			return l.errorf("unexpected eof in getter")
		case r == '[':
			l.backup()
			return lexGetterOpening
		case r == ']':
			l.backup()
			return lexGetterClosing
		case unicode.IsSpace(r):
			//ignore
		default:
			return l.errorf("unexpected char inside getter: %q",r)
		}
	}
	return nil
}

func lexGetterOpening(l *Lexer) stateFn {
	l.next()
	l.emit(token.LeftGetter)
	l.getterDepth++
	return lexGetterType
}

func lexGetterType(l *Lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == eof:
			l.errorf("unclosed getter")
		case r == ':':
			l.backup()
			l.emit(token.GetterType)
			l.next()
			l.emit(token.GetterSeparator)
			return lexGetterValue
		case r == '[':
			l.backup()
			l.emit(token.GetterType)
			return lexInsideGetter
		}
	}
	return l.errorf("didn't found getter separator")
}

func lexGetterValue(l *Lexer) stateFn {
	for r := ' '; r != eof; r = l.next() {
		if r == ']' {
			l.backup()
			l.emit(token.GetterValue)
			return lexGetterClosing
		}
	}

	return l.errorf("couldn't parse getter value")
}

func lexGetterClosing(l *Lexer) stateFn {
	l.next()
	l.emit(token.RightGetter)
	l.getterDepth--
	if l.getterDepth == 0 {
		return lexInsidePromise
	} else {
		return lexInsideGetter
	}
}

func isValidPromiseNameRune(r rune) bool {
	return r == '-' || r == '_' || r == ' ' || unicode.IsLetter(r) || unicode.IsDigit(r)
}


func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
