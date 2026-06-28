1. make changes...


2. build...
```bash
nix build
```

3. commit and push (or create a PR)
```bash
git commit -m ...

github-proxy remove kanafm/bloggo # or you can use a new terminal without https_proxy set. this is so that git pushes to real GitHub.
git push origin master
github-proxy add .
```
