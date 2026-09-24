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
