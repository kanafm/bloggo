package bloggo

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type BuildRequest struct {
	EntryPoints []string
	OutDir      string
	Processors  []Processor
}

type Bloggable interface {
	Build(request BuildRequest) error
}

type Bloggo struct{}

func New() *Bloggo {
	return &Bloggo{}
}

func (b Bloggo) handle(path string, outdir string, processors []Processor) error {
	// select the processor
	var activeProcessor *Processor = nil
	for _, processor := range processors {
		if !processor.CanHandle(path) {
			continue
		}
		activeProcessor = &processor
		break
	}

	if activeProcessor == nil {
		return errors.New("No processors can handle " + path)
	}

	// create directory for output file
	ext := filepath.Ext(path)
	name := strings.TrimSuffix(filepath.Base(path), ext)

	absoluteOutDir, err := filepath.Abs(outdir)
	if err != nil {
		return err
	}

	err = os.MkdirAll(absoluteOutDir, 0o755)
	if err != nil {
		return err
	}

	outFile := fmt.Sprintf("%s/%s.html", absoluteOutDir, name)

	// handle it with processor
	err = (*activeProcessor).Handle(path, outFile)
	if err != nil {
		return err
	}

	fmt.Println("\t. -> ", outFile)
	return nil
}

func (b Bloggo) Build(request BuildRequest) error {
	for ei := range request.EntryPoints {
		e := request.EntryPoints[ei]

		err := filepath.WalkDir(e, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}

			//path is absolute path
			fmt.Println(". processing ", path)
			err = b.handle(path, request.OutDir, request.Processors)
			if err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return err
		}

	}

	return nil
}

type Processor interface {
	// CanHandle determines whether this processor can handle this file
	CanHandle(fileLocation string) bool

	// Handle processes file at file location, writing to output file
	Handle(fileLocation string, outputFile string) error
}

type NoOpProcessor struct{}

func (NoOpProcessor) CanHandle(fileLocation string) bool {
	return true
}

func (n NoOpProcessor) Handle(fileLocation string, outputFile string) error {
	if !n.CanHandle(fileLocation) {
		return errors.New("Cannot handle file: " + fileLocation)
	}

	data, err := os.ReadFile(fileLocation)
	if err != nil {
		return err
	}

	htmlified := fmt.Sprintf(
		"<div><pre style=\"font-family: ui-monospace, monospace; margin:0;\">%s</pre></div>",
		string(data))

	err = os.WriteFile(outputFile, []byte(htmlified), 0o644)
	if err != nil {
		return err
	}

	return nil
}
