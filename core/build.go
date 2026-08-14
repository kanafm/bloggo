package core

type Bloggable interface {
	Build(request BuildRequest) error
}
type BuildRequest struct {
	EntryPoints []string
	OutDir      string
	Processors  []Processor
	Router      Router
}
