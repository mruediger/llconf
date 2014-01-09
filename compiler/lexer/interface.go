package lexer

import (
	"github.com/mruediger/llconf/compiler/token"
)

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
