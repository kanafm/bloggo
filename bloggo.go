package bloggo

import "github.com/kanafm/bloggo/core"

// Bloggo is a static site builder.
type Bloggo struct{}

var _ core.Bloggable = (*Bloggo)(nil)

// Creates a new Bloggo.
func New() *Bloggo {
	return &Bloggo{}
}

// Builds your site.
func (b Bloggo) Build(request core.BuildRequest) error {
	return core.Build(request)
}
