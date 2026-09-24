# Agent notes

Personal site published at https://pngdeity.ru/ from this repo. Read `README.md`
for the human-facing architecture; this file is the short version an agent needs
to avoid re-discovering the layout.

## Layout

- `src/` is uploaded to GitHub Pages verbatim as the site root. Static pages
  have **no server-side includes**, so shared markup is either a Hugo partial
  (blog) or intentionally duplicated by hand (static pages).
- `hugo-src/` is the Hugo blog. It builds with `baseURL = /blog/` and its output
  is merged into `src/blog/` at deploy, so every Hugo page is under `/blog/`.
- `tools/ci/` holds the Go helpers (`generatemetadata`, `validatemetadata`,
  `sitecheck`, `healthcheck`, `manageincident`, `resolverollback`).

## Navigation and links (single sources)

- Blog header nav: `hugo-src/layouts/_partials/site-navigation.html`
- Blog footer links + copyright: `hugo-src/data/footer.toml`
- Static pages (hand-maintained, mirror the footer): `src/index.html`,
  `src/links.html`, `src/404.html`

Use root-absolute hrefs (`/...`). Change link lists in `data/footer.toml` and the
static pages together.

## Blog posts

`hugo-src/content/posts/<slug>.md`; `draft = true` posts are excluded from the
build. `hugo-src/archetypes/default.md` is the template.

## Gotchas

- Hugo prepends the `/blog/` subpath to `[[menus.main]]` URLs, so the shared nav
  is an overridden partial (`site-navigation.html`), not a menu.
- `enableGitInfo = true` requires a git repo; building a copied tree outside git
  fails. Pass `--enableGitInfo=false` for throwaway builds.
- The TOML formatter strips the indentation under `[markup]` in
  `hugo-src/hugo.toml` on edit; it is cosmetic, not a semantic change.
- `.serena/` is ignored agent-tooling state; do not commit it.

## Commands

```sh
hugo --source ./hugo-src --gc --minify   # build the blog
hugo server --source ./hugo-src          # preview the blog
cd tools/ci && go test ./...             # metadata tool unit tests
cd tools/ci && go run ./sitecheck -repo ../.. -src ../../src   # link check
```

`sitecheck` is advisory and expects the merged tree (Hugo output staged at
`src/blog/`), which the deploy workflow does.

Toolchain: CI pins Hugo 0.154.4 (`.github/workflows/build-deploy.yaml`), Node 24
(`.node-version`), and Go per `tools/ci/go.mod`.

## Git

- Commits must be signed (`git commit -S`).
- `main` deploys on push; `content/develop-drafts` carries WIP posts.
