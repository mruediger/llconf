package promise

import (
	"io"
	"strings"
)

type ReadvarPromise struct {
	VarName Argument
	Exec Promise
}

type ReadvarWriter struct {
	writer io.Writer
	bytes []byte
}

func (w *ReadvarWriter) Write(p []byte) (n int, err error) {
	w.bytes = append(w.bytes, p...)
	return w.writer.Write(p)
}

func (p ReadvarPromise)	Desc(arguments []Constant) string {
	return ""
}

func (p ReadvarPromise) Eval(arguments []Constant, logger *Logger, vars *Variables) bool {
	bytes := []byte{}
	
	wrapped_stdout := ReadvarWriter{
		writer: logger.Stdout,
		bytes: bytes,
	}
	
	wrapped_logger := &Logger{
		Stdout: &wrapped_stdout,
		Stderr: logger.Stderr,
		Info: logger.Info,
		Changes: logger.Changes,
		Tests: logger.Tests,
	}

	result := p.Exec.Eval(arguments, wrapped_logger, vars)

	name  := p.VarName.GetValue(arguments, vars)
	value := string(wrapped_stdout.bytes)
		
	(*vars)[name] = strings.TrimSpace(value)
	
	return result
}

