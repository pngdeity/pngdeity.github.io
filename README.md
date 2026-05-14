Source for my personal website.

## Deployment metadata generation

The GitHub Pages deploy workflow (`.github/workflows/build-deploy.yaml`) generates and validates deployment metadata before uploading `src/`:

- `scripts/generate-site-metadata.mjs`
- `scripts/validate-site-metadata.mjs`

These Node.js scripts create/validate `robots.txt`, `sitemap.xml`, `llms*.txt`, and `.well-known` files under `src/`.
