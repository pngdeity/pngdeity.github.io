Source for my personal website.

## Deployment metadata generation

The GitHub Pages deploy workflow (`.github/workflows/build-deploy.yaml`) generates and validates deployment metadata before uploading `src/`:

- `tools/ci/generatemetadata`
- `tools/ci/validatemetadata`

These Go applications create/validate `robots.txt`, `sitemap.xml`, `llms*.txt`, and `.well-known` files under `src/`.

## Auto-rollback for bad Pages deployments

The repository includes an automated rollback workflow for GitHub Pages incidents:

- **Primary deploy workflow:** `.github/workflows/build-deploy.yaml`
- **Rollback workflow:** `.github/workflows/auto-rollback-pages.yaml`

### Trigger criteria (what counts as a bad deployment)

A deployment is considered bad only when **post-deploy health verification fails** after a successful Pages deploy step.

Deterministic health checks validate HTTP 200 responses for:

- `/`
- `/blog/`
- `/app/`
- `/sitemap.xml`
- `/robots.txt`

Build or test failures alone do not trigger auto-rollback.

### What auto-rollback does

When a bad deployment is detected on `main`:

1. Finds the most recent known-good successful run of the deploy workflow before the failing run.
2. Promotes the exact immutable `github-pages` artifact from that prior successful run and directly redeploys it, entirely skipping any rebuilding of Hugo or Blazor.
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
