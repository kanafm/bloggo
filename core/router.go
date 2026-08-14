package core

import (
	"path/filepath"
	"strings"
)

// Router takes an input path (relative to its entry point) and generates an output path (relative to OutDir).
type Router interface {
	Route(inputPath string) (outputPath string, ok bool)
}

// MirrorRouter preserves an input's relative path and changes its extension to .html.
type MirrorRouter struct{}

var _ Router = MirrorRouter{}

func (MirrorRouter) Route(inputPath string) (string, bool) {
	extension := filepath.Ext(inputPath)
	return strings.TrimSuffix(inputPath, extension) + ".html", true
}
