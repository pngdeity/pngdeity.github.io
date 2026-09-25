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

## Theming (light/dark)

- One convention: a `data-theme` attribute on `<html>` (`light`/`dark`), a
  `pg-theme` `localStorage` key, and CSS custom properties (`--pg-bg`,
  `--pg-fg`, ...). `prefers-color-scheme` is the fallback when no choice is
  stored.
- Blog: vars live in `hugo-src/assets/ananke/css/theme.css`, registered via
  `params.custom_css` in `hugo-src/hugo.toml` (the ananke pipeline only finds
  custom CSS under `assets/ananke/css/`, not `assets/css/`). The no-flash
  pre-paint script is `hugo-src/layouts/_partials/head-additions.html`; the
  click handler is `hugo-src/assets/js/theme.js`, loaded by the
  `site-scripts.html` override.
- Static half: `src/style.css` holds the same vars, each page carries the
  same inline pre-paint snippet, and `src/theme.js` is a copy of the blog
  handler. Keep the key, var names, and snippet in sync by hand.
- The toggle is a fixed top-right circular button showing the theme it will
  switch to (moon in light, sun in dark); markup lives in
  `hugo-src/layouts/_partials/theme-toggle.html` and is duplicated into the
  static pages. `theme.js` keeps `aria-label`/`aria-pressed` in sync.
- `src/` and `hugo-src/` are two implementations of one convention; they drift
  unless edited together.

## Tags

`disableKinds` is not set; `[taxonomies] tag = 'tags'` in `hugo-src/hugo.toml`
enables tag pages under `/blog/tags/`. Posts carry `tags = [...]` in front
matter; `hugo-src/layouts...` uses the theme defaults, restyled by theme.css.

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
