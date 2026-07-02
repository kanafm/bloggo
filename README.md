1. make changes...


2. build...
```bash
nix build
```

(development with a separate consumer package)
* checkout a branch (`git checkout -b feature/blah`)
* commit (`git commit -m ...` or `git commit --amend`)
* build library (`nix build`)
* make remote aware of it.
  * if you use github proxy: `github-proxy update kanafm/bloggo`
  * if you use github.com: `git push origin feature/blah`
* now cd into consumer package
* update your go.mod: `go get github.com/kanafm/bloggo@feature/blah`
* and build your consumer package (like `nix build`)

(pushing to remote)
* commit and push (or create a PR)
```bash
git commit -m ...

github-proxy remove kanafm/bloggo # or you can use a new terminal without https_proxy set. this is so that git pushes to real GitHub.
git push origin master
github-proxy add .
```
