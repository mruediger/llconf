package promise

import "fmt"

type InDir struct {
	dir     Argument
	promise Promise
}

func (p InDir) Desc(arguments []Constant) string {
	return fmt.Sprintf("(indir %s %s)", p.dir, p.promise.Desc(arguments))
}

func (p InDir) Eval(arguments []Constant, ctx *Context, stack string) bool {
	copyied_ctx := *ctx
	copyied_ctx.InDir = p.dir.GetValue(arguments, &ctx.Vars)
	return p.promise.Eval(arguments, &copyied_ctx, stack)
}

func (p InDir) New(children []Promise, args []Argument) (Promise, error) {

	if len(args) != 1 {
		return nil, fmt.Errorf("(indir) needs exactly on argument, found %d", len(args))
	}

	if len(children) != 1 {
		return nil, fmt.Errorf("(indir) needs exactly on child promise, found %d", len(children))
	}

	return InDir{args[0], children[0]}, nil
}
