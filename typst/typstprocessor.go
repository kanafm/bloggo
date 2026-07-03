package typst

import (
	"path"

	"github.com/kanafm/bloggo/core"
)

type TypstProcessor struct {
	TypstClient *TypstClient
}

var _ core.Processor = (*TypstProcessor)(nil)

func (t TypstProcessor) CanHandle(inputFile string) bool {
	return path.Ext(inputFile) == ".typ" || path.Ext(inputFile) == ".typst"
}

func (t TypstProcessor) Handle(inputFile string, outputFile string) error {
	return t.TypstClient.CompileWithHtmlFeatures(inputFile, outputFile)
}

func NewProcessor() *TypstProcessor {
	return &TypstProcessor{
		TypstClient: NewClient(),
	}
}
