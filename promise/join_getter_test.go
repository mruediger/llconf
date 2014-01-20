package promise

import (
	"testing"
)

func TestJoinGetter(t *testing.T) {
	args := []Argument{ArgGetter{0}, ArgGetter{1}, ArgGetter{2}}
	join := JoinArgument{args}
	result := join.GetValue([]Constant{Constant("hello"), Constant(" "), Constant("world")},
		&Variables{})
	equals(t, "hello world", result)

	desc := join.String()
	equals(t, "joinargs-> arg->0 + arg->1 + arg->2", desc)
}
