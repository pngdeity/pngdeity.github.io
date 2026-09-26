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
  It is a self-contained Hugo project: templates live in `hugo-src/layouts/`,
  styles in `hugo-src/assets/css/`. **There is no `themes/` directory** — the
  site owns its own theme (ananke was removed).
- `tools/ci/` holds the Go helpers (`generatemetadata`, `validatemetadata`,
  `sitecheck`, `healthcheck`, `manageincident`, `resolverollback`).

## Navigation and links (single sources)

- Blog header nav: `hugo-src/layouts/_partials/site-navigation.html`
- Blog footer links + copyright: `hugo-src/data/footer.toml`
- Static pages (hand-maintained, mirror the footer): `src/index.html`,
  `src/links.html`, `src/404.html`

Use root-absolute hrefs (`/...`). Change link lists in `data/footer.toml` and the
static pages together.

## Design identity

`DESIGN.md` is the single source of truth for the visual identity: the pngdeity
brand, its motifs, the colour/type/layout rules, and the `Do's and Don'ts`. Read
it before changing anything visual. `TODO.md` holds what is still open.

The aesthetic position is **Deadpan Utility** -- Digital Brutalism and
post-internet surrealism. The interface is a stark, machined container that
presents absurd or heavy material with total seriousness. Three principles follow
from it and govern new work: the **Bowtie Principle** (bizarre elements presented
as entirely factual and ordinary), **Active Interrogation** (errors confront
rather than apologise; the 404 is the model), and **no gimmicks** (no simulated
terminals, no decorative system chrome -- the machined quality is real structure,
not theatre).

The framing is *not* editorial, tasteful, or midcentury-modern. That was
considered and rejected; see the `## Declined` section of `DESIGN.md` before
reasoning from elegance, restraint, or "a composed document".

Hard rules, in short:

- The brand is `pngdeity`, **always lowercase**, in titles and prose alike. Never
  capitalise it, never set it in display capitals, never decorate it with a logo,
  halo, or icon.
- **Do not invent a tagline.** The old site title `The End of the Internet` was
  an accident; it is retired and must not be replaced by another arbitrary
  phrase.
- The homepage stays brief: image, name, caption, five links. Nothing else.
- The header/footer band on the blog stays dark with light text in **both**
  themes.
- `[*]` is site identity, injected by CSS on the homepage and set in monospace.
- Ornament carries no information, so it is omitted. Do not add shadows, depth,
  gradients, or decoration to fill space.

## Theming (light/dark)

- One convention: a `data-theme` attribute on `<html>` (`light`/`dark`), a
  `pg-theme` `localStorage` key, and CSS custom properties (`--pg-bg`,
  `--pg-fg`, ...). `prefers-color-scheme` is the fallback when no choice is
  stored.
- Blog: vars live in `hugo-src/assets/css/theme.css`. The whole blog stylesheet
  is the concatenation of `css/utilities.css` + `css/_code.css` + `css/theme.css`
  by `hugo-src/layouts/_partials/site-style.html`, emitted as `/blog/css/site.css`
  (fingerprinted in production). There is no `custom_css` config and no
  `/ananke/` asset path. The no-flash pre-paint script is
  `hugo-src/layouts/_partials/head-additions.html`; the click handler is
  `hugo-src/assets/js/theme.js`, loaded by `site-scripts.html`.
- Static half: `src/style.css` holds the same vars, each page carries the
  same inline pre-paint snippet, and `src/theme.js` is a copy of the blog
  handler. Keep the key, var names, and snippet in sync by hand.
- The toggle is a fixed top-right circular button showing the theme it will
  switch to (moon in light, sun in dark); markup lives in
  `hugo-src/layouts/_partials/theme-toggle.html` and is duplicated into the
  static pages. `theme.js` keeps `aria-label`/`aria-pressed` in sync.
- `src/` and `hugo-src/` are two implementations of one convention; they drift
  unless edited together.

## Typography

Two registers, both recorded in `DESIGN.md`.

- **Chrome** (nav, section labels, metadata, share links, timestamps, toggle,
  `[*]` markers): the platform UI stack, declared once on `body` in
  `src/style.css` and applied to the blog via `.pg-site` in
  `hugo-src/assets/css/theme.css`. Auxiliary pages inherit it; do not add a
  page-level `font-family`.
- **Blog prose** (article body and its headings): **Bitter**, a self-hosted slab
  serif at `hugo-src/assets/fonts/bitter/` (400 normal, 700 normal, 400 italic
  woff2). Licensed OFL-1.1; `OFL.txt` ships beside the files and must travel with
  them. Applied by the reading-register rules at the end of
  `hugo-src/assets/css/theme.css` -- these need `article h1.athelas` and similar
  class specificity to beat the `athelas` utility in `utilities.css`.
- Headings are weight `600`, not `700`. Keep prose rules confined to prose:
  a serif leaking into the chrome collapses the two registers.

## Tags

`disableKinds` is not set; `[taxonomies] tag = 'tags'` in `hugo-src/hugo.toml`
enables tag pages under `/blog/tags/`. Posts carry `tags = [...]` in front
matter: `/blog/tags/` is rendered by `layouts/taxonomy.html` and each
`/blog/tags/<tag>/` by `layouts/list.html` (branching on `.Kind == "term"`).

## Generated site files

Some files under `src/` are produced and validated at deploy, not hand-written.
Do not edit them by hand; change `tools/ci/generatemetadata`.

- `src/robots.txt`, `src/sitemap.xml`
- `src/llms.txt`, `src/llms-full.txt` (LLM-facing summaries)
- `src/.well-known/ai-catalog.json` (a `did:web:pngdeity.ru` descriptor)
- `src/.well-known/keybase.txt` (a Keybase proof referring to the commit-signing
  key `0D762ADE0E00C8DF`)

`validatemetadata` runs before upload and fails the deploy if they are wrong.

## Workflows

Four workflows; `AGENTS.md` historically named only the first.

| Workflow | Trigger | Effect |
| --- | --- | --- |
| `build-deploy.yaml` | push to `main`, manual | builds, merges, deploys, health-checks |
| `validate-site.yaml` | pull request, manual | runs `sitecheck` **advisory** (`continue-on-error`) |
| `ci-tools.yaml` | push, pull request | `go fmt` / `vet` / `build` / `test` -- **blocking** |
| `auto-rollback-pages.yaml` | deploy-run failure, manual | republishes the last good Pages artifact |

Green CI is the acceptance bar before merging to `main`, which deploys on push.

## Blog posts

`hugo-src/content/posts/<slug>.md`; `draft = true` posts are excluded from the
build. `hugo-src/archetypes/default.md` is the template.

## Gotchas

- Hugo prepends the `/blog/` subpath to `[[menus.main]]` URLs, so the shared nav
  is a partial (`site-navigation.html`), not a menu.
- `enableGitInfo = true` requires a git repo; building a copied tree outside git
  fails. Pass `--enableGitInfo=false` for throwaway builds.
- The TOML formatter strips the indentation under `[markup]` in
  `hugo-src/hugo.toml` on edit; it is cosmetic, not a semantic change.
- `.serena/` is ignored agent-tooling state; do not commit it.
- `show_recent_posts` and the `[params.ananke.social.networks.*]` table (the
  share-link label/link/particle definitions) are declared in `hugo-src/hugo.toml`.
  They previously came from the theme's merged config; if they go missing the
  homepage loses its post cards and the share links render as an empty `#sharing`
  div. The `ananke` prefix in those keys is a legacy name, not a dependency.
- `hugo-src/assets/css/utilities.css` is a hand-maintained subset of Tachyons
  covering only the classes our templates emit. Adding markup that uses a new
  utility class requires adding the rule here; nothing generates it.
- Hugo strips CSS comments in production builds. Comments in `assets/css/` exist
  for agents reading the source, not for the served output.
- **There is no Node/Sass pipeline.** `.node-version`, `package.json` and
  `package-lock.json` were removed with the previous theme; the only consumer was
  its `css.Sass` call. Do not add `npm ci` steps or assume `node_modules` exists.
- The Hugo template APIs are current: `locale` (not `languageCode`),
  `.Language.Locale`, `.Language.Direction`, and `hugo.Data` (not `.Site.Data`).
  The build is warning-free; keep it that way rather than reintroducing the
  deprecated spellings from older theme code.

## Commands

```sh
hugo --source ./hugo-src --gc --minify   # build the blog
hugo server --source ./hugo-src          # preview the blog
cd tools/ci && go test ./...             # metadata tool unit tests

# sitecheck needs the merged tree; stage it first or it reports false positives:
rm -rf /tmp/stage && cp -r src /tmp/stage && mkdir -p /tmp/stage/blog
cp -r hugo-src/public/. /tmp/stage/blog/
cd tools/ci && go run ./sitecheck -repo ../.. -src /tmp/stage

# validate the design record (prints zero errors):
npx -p @google/design.md designmd lint DESIGN.md
cd tools/ci && go fmt ./... && go vet ./... && go build ./... && go test ./...
```

`sitecheck` is advisory and expects the merged tree (Hugo output staged at
`src/blog/`), which the deploy workflow does. Running it against a bare copy of
`src/` reports false unresolved refs.

Toolchain: **Hugo and Go only.** CI pins Hugo 0.154.4
(`.github/workflows/build-deploy.yaml`) and Go per `tools/ci/go.mod`. There is no
Node toolchain — no `package.json`, no `.node-version`, and no Sass/PostCSS
pipeline.

## Git

- Commits must be signed (`git commit -S`).
- `main` deploys on push; `content/develop-drafts` carries WIP posts.
