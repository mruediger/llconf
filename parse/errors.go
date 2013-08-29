package parse

type DuplicatePromise struct {
	nameOfPromise string
}

func (this DuplicatePromise) Error() string {
	return "duplicated Promise: " + this.nameOfPromise
}

type IllegalPromisePosition struct {
	nameOfPromise string
}

func (this IllegalPromisePosition) Error() string {
	return "("+ this.nameOfPromise + ") promise not allowed in primary position"
}

type NamedPromiseArgc struct {
	argc int
	nameOfPromise string
}

func (this NamedPromiseArgc) Error() string {
 	if this.argc < 1 {
		return "the named promise (" + this.nameOfPromise + ") needs to contain anoter promise"
	} else {
		return "the named promise (" + this.nameOfPromise + ") can only contain one promise"
	}
}

type MissingPromise struct {
	nameOfPromise string
}

func (this MissingPromise) Error() string {
	return "couldn't find promise (" + this.nameOfPromise + ")"
}

type UnknownGetterType struct {
	name string
}

func (this UnknownGetterType) Error() string {
	return "unknown getter type: " + this.name
}

type UnexpectedEOF struct {}

func (this UnexpectedEOF) Error() string {
	return "unexpected end of input"
}
