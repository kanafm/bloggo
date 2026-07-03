package bloggo

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	. "github.com/kanafm/bloggo/core"
)

// Bloggo is a static site builder.
type Bloggo struct{}

var _ Bloggable = (*Bloggo)(nil)

// Creates a new Bloggo.
func New() *Bloggo {
	return &Bloggo{}
}

// Builds your site.
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
