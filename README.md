

## usage

bloggo is a zero-dependency *library* and it takes input directory/ies and outputs a site.

```go
b := bloggo.New()

err := b.Build(BuildRequest{
    EntryPoints: []string{"./docs"},
    OutDir:      "./site",
    Processors:  []Processor{typst.NewProcessor(), NoOpProcessor{}},
})
```

where,
* EntryPoints: the set of entrypoints (directories)
* OutDir: where to write your site
* Processors: processors to compile every file in your entrypoints into HTML

## roadmap


* [x] Processors - supports Typst as well as provide-your-own
* [ ] Templating - headers, footers, navigation
* [ ] Invert responsibility for routing logic (i.e. `docs/hello/1.md` maps to `docs/hello-1.md`). Currently everything generates a flat list of HTML (like `1.html`) and gives consumer no flexibility.
* [ ] Bring your own styling/CSS

## development
1. make changes...


2. build...
```bash
nix build
```

3. enjoy