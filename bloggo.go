package bloggo

import (
	"fmt"
	"io"
	"os"
)

type Bloggo struct {
	out io.Writer
}

func Of() *Bloggo {
	return &Bloggo{out: os.Stdout}
}

func OfWriter(w io.Writer) *Bloggo {
	return &Bloggo{out: w}
}

func (w *Bloggo) Writeln(text string) *Bloggo {
	fmt.Fprintln(w.out, text)
	return w
}

func (w *Bloggo) Write(text string) *Bloggo {
	fmt.Fprint(w.out, text)
	return w
}

func (w *Bloggo) Writef(format string, args ...any) *Bloggo {
	fmt.Fprintf(w.out, format, args...)
	return w
}

func (w *Bloggo) SayHello() *Bloggo {
	fmt.Fprintln(os.Stdout, "Hello, World!")
	return w
}
