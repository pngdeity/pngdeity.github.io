Source for my personal website, published at https://pngdeity.ru/.

## Site structure

The published site is assembled from two sources:

- **`src/`** — the static root, uploaded to GitHub Pages verbatim. `index.html`
  (`/`), `links.html` (`/links.html`), `404.html` (served for the whole domain),
  `style.css`, plus `pngdeity_files/` (résumé, vCard, images) and `downloads/`
  (files linked from posts).
- **`hugo-src/`** — the Hugo blog. It builds with
  `baseURL = 'https://pngdeity.ru/blog/'` and is merged into `src/blog/` at
  deploy, so every Hugo page lives under `/blog/`. Posts are in
  `content/posts/`; the shared nav/footer are
  `layouts/_partials/site-navigation.html` and `site-footer.html`, with footer
  links in `data/footer.toml`. `themes/ananke/` is a vendored submodule.
- **`tools/ci/`** — Go helpers used by CI (`generatemetadata`,
  `validatemetadata`, `sitecheck`, `healthcheck`, `manageincident`,
  `resolverrollback`).

Static pages are copied as-is with no server-side includes, so shared markup is
either a Hugo partial (blog) or intentionally duplicated by hand (static pages).
Use root-absolute links (`/...`) everywhere.

## Editing navigation and links

| Surface | Source |
| --- | --- |
| Blog header nav | `hugo-src/layouts/_partials/site-navigation.html` |
| Blog footer links + copyright | `hugo-src/data/footer.toml` |
| Home page | `src/index.html` |
| Link hub | `src/links.html` |
| 404 page | `src/404.html` |

The static lists mirror `data/footer.toml` by hand; update both when links change.

## Design and identity

`DESIGN.md` records the site's visual identity in the
[DESIGN.md format](https://github.com/google-labs-code/design.md): design tokens
in YAML front matter, and the reasoning behind them in prose.

The short version:

- The brand is the handle **`pngdeity`**, always lowercase. It is a username, not
  a title, and it is never capitalised or decorated.
- The motif is **the image that looks back** — the eye/portal favicon, the
  homepage photograph, and the Paalen quotation on the 404 page.
- Colour, type, and layout rules live in the document, along with explicit
  Do's and Don'ts.

`TODO.md` lists what is not yet settled, including the light-mode accent colour
and a possible CI check that the two stylesheets' tokens stay in sync.

## Theming (light and dark)

The site follows the OS color scheme by default and offers a manual toggle that
persists in `localStorage` under the `pg-theme` key.

- The convention is a `data-theme` attribute (`light`/`dark`) on `<html>` plus
  CSS custom properties (`--pg-bg`, `--pg-fg`, `--pg-link`, ...).
- Blog styles: `hugo-src/assets/ananke/css/theme.css`, registered through
  `params.custom_css`. The pre-paint script is
  `hugo-src/layouts/_partials/head-additions.html`; the toggle behavior is
  `hugo-src/assets/js/theme.js`.
- Static pages: the same variables live in `src/style.css`, with an inline
  pre-paint snippet and a copy of the handler at `src/theme.js`.

Because `src/` cannot include the blog's partials, these are two hand-synced
implementations of one convention — change them together.

## Tags

`hugo-src/hugo.toml` sets `[taxonomies] tag = 'tags'`, so posts carry
`tags = [...]` front matter and are browsable at `/blog/tags/` (each tag gets a
term page). This is the topical counterpart to the chronological `/blog/` list.

## Search-engine and social sharing

Single post pages render text-only share links (Reddit and X) from the theme's
share partial. Networks are set in `hugo-src/hugo.toml` under
`[params.ananke.social.share]`. The share text is the post's front-matter
`description` (falling back to a truncated summary), via the override at
`hugo-src/layouts/_partials/func/social/getShareLink.html` — set `description`
on every post you want shared well.

## Adding a blog post

Copy `hugo-src/archetypes/default.md` to `hugo-src/content/posts/<slug>.md` and
set `draft = false` to publish. Drafts are excluded from the build and `/blog/`.

## Local build and validation

```sh
hugo --source ./hugo-src --gc --minify   # build the blog
hugo server --source ./hugo-src          # preview the blog
cd tools/ci && go test ./...             # metadata tool unit tests
cd tools/ci && go run ./sitecheck -repo ../.. -src ../../src   # link check
```

`sitecheck` expects the merged tree (Hugo output staged as `src/blog/`), which
the deploy workflow performs.

## Toolchain

- **Hugo** (`extended`) — CI pins `HUGO_VERSION` in
  `.github/workflows/build-deploy.yaml`; local builds with the system Hugo work.
- **Go** — version in `tools/ci/go.mod`; used only by the CI helpers.
- **Node** — `.node-version` and `package.json` (theme asset pipeline).

## Deployment

`.github/workflows/build-deploy.yaml` runs two jobs:

1. `build_hugo` — builds `hugo-src/` and uploads the result.
2. `merge_and_deploy` — stages the Hugo output at `src/blog/`, runs the advisory
   site check and the metadata generator/validator, creates `src/.nojekyll`,
   uploads `src/` to GitHub Pages, then health-checks `/`, `/blog/`,
   `/sitemap.xml`, and `/robots.txt`.

`CNAME` pins the canonical origin to `pngdeity.ru`.

## License

GPL-3.0 — see `LICENSE.md`.

## Deployment metadata generation

The GitHub Pages deploy workflow (`.github/workflows/build-deploy.yaml`) generates and validates deployment metadata before uploading `src/`:

- `tools/ci/generatemetadata`
- `tools/ci/validatemetadata`

These Go applications create/validate `robots.txt`, `sitemap.xml`, `llms*.txt`, and `.well-known` files under `src/`.

## Advisory site checks

`tools/ci/sitecheck` resolves every local `href`/`src`/`poster` reference in the merged `src/` tree (including the Hugo output staged under `src/blog/`), follows same-origin absolute URLs back to files on disk, and asserts that `CNAME` matches the canonical origin. It is **advisory**: findings are reported as workflow annotations and in the step summary, but `continue-on-error` keeps the run green.

- Runs on pull requests via `.github/workflows/validate-site.yaml`.
- Runs before the Pages upload in `.github/workflows/build-deploy.yaml`.

Unit tests for the metadata tools run in `.github/workflows/ci-tools.yaml` and are blocking.

## Auto-rollback for bad Pages deployments

The repository includes an automated rollback workflow for GitHub Pages incidents:

- **Primary deploy workflow:** `.github/workflows/build-deploy.yaml`
- **Rollback workflow:** `.github/workflows/auto-rollback-pages.yaml`

### Trigger criteria (what counts as a bad deployment)

A deployment is considered bad only when **post-deploy health verification fails** after a successful Pages deploy step.

Deterministic health checks validate HTTP 200 responses for:

- `/`
- `/blog/`
- `/sitemap.xml`
- `/robots.txt`

Build or test failures alone do not trigger auto-rollback.

### What auto-rollback does

When a bad deployment is detected on `main`:

1. Finds the most recent known-good successful run of the deploy workflow before the failing run.
2. Promotes the exact immutable `github-pages` artifact from that prior successful run and directly redeploys it, entirely skipping any rebuilding of Hugo.
3. Runs post-rollback health verification checks using a robust Go-based healthchecker.
4. Publishes incident and rollback details in workflow summaries and creates/updates a tracking GitHub issue.

### Safety controls

- Rollback only runs for failed deploy runs caused by the `Post-deploy Health Verification` step.
- One auto-rollback attempt per incident run (skips repeated run attempts).
- Never rolls back to the same failing commit.
- Skips rollback when no known-good target can be resolved.

### Manual rollback

Use **Actions → Auto Rollback Bad Pages Deployment → Run workflow** and provide:

- `rollback_sha` (required): commit SHA to redeploy
- `bad_sha` (optional): failing commit SHA for context in summaries/issues

### Temporarily disable auto-rollback

Set repository variable `AUTO_ROLLBACK_ENABLED` to `false`.

- This disables automated rollback from failed deploy workflow runs.
- Manual rollback via workflow dispatch remains available.
