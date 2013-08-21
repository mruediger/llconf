package llconf

type Unmarshaler interface {
	UnmarshalLLC([]rune) error
}

