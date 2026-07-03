package typst

import (
	"cmp"
	"os"
	"os/exec"
)

type TypstClient struct {
	Exe string
}

func NewClient() *TypstClient {
	return &TypstClient{
		Exe: cmp.Or(os.Getenv("BLOGGO_TYPST_EXE"), "typst"),
	}
}

func NewClientWithExe(exe string) *TypstClient {
	return &TypstClient{
		Exe: exe,
	}
}

func (t TypstClient) Compile(inputFile string, outputFile string, rest ...string) error {
	args := []string{"compile", inputFile, outputFile}
	args = append(args, rest...)
	cmd := exec.Command(t.Exe, args...)
	return cmd.Run()
}

func (t TypstClient) CompileWithHtmlFeatures(inputFile string, outputFile string) error {
	return t.Compile(inputFile, outputFile, "--features", "html")
}
