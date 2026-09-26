# TODO

Open items. These are deliberately not in `DESIGN.md`, which records only what
is settled. Priorities are unknown unless stated.

## CI guard for stylesheet token drift

`src/style.css` and `hugo-src/assets/ananke/css/theme.css` hold independent
copies of the `--pg-*` tokens, plus a `@media (prefers-color-scheme: dark)`
duplicate of the dark palette in each. Nothing enforces that they agree.

A `tools/ci` check could parse both files and fail when the token sets diverge.
Priority unknown. Until it exists, the sync rule is a convention only.

## Static-page drift after the homepage recomposition

The static pages do not match the homepage's conventions:

- `src/links.html` still centres its list and carries no `[*]` markers.
- `src/credits.html` uses a `.brand` rule of `1.5em` / `normal`, while
  `src/links.html` uses `1.5rem` / `400` / `-0.01em`.
- `src/links.html` and `src/credits.html` still declare `font-family: sans-serif`
  in their page style, while the homepage removed it so the page inherits the
  stack from `src/style.css`.

`DESIGN.md` describes the settled system; these pages should be brought into
line with it, or the description should be narrowed.

## Unreferenced image assets

`src/pngdeity_files/gladstone.jpg` and `gladstone_modified.jpg` are no longer
referenced by any page since the homepage moved to `gladstone-at-work.jpg`.
Either delete them or document why they are kept. (Deletion was offered and not
yet answered.)

## `content/develop-drafts` is behind `main`

The branch has not been re-synced since the earlier sync at `b500c59`. It should
be merged up with `main` before it accumulates more conflict surface.
