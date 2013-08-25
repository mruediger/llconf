package promise


type DummyPromise struct {
	StringValue string
	EvalValue bool
}

func (p DummyPromise) String() string {
	return p.StringValue
}

func (p DummyPromise) Eval() bool {
	return p.EvalValue
}
