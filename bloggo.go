package bloggo

import (
	"fmt"
	"io"
	"os"
)

type Writer struct {
	out io.Writer
}

func New() *Writer {
	return &Writer{out: os.Stdout}
}

func NewWithWriter(w io.Writer) *Writer {
	return &Writer{out: w}
}

func (w *Writer) Writeln(text string) {
	fmt.Fprintln(w.out, text)
}

func (w *Writer) Write(text string) {
	fmt.Fprint(w.out, text)
}

func (w *Writer) Writef(format string, args ...any) {
	fmt.Fprintf(w.out, format, args...)
}
