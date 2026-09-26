# TODO

Open items. These are deliberately not in `DESIGN.md`, which records only what
is settled. Priorities are unknown unless stated.

`DESIGN.md` also carries a `## Declined` section listing proposals that were
considered and rejected, so they are not silently re-proposed. Add to that
section rather than here when the answer is "no", and here when the answer is
"not yet".

## Blog chrome has no deliberate element

Raised while choosing the reading register: the blog chrome is a dark band, the
wordmark, five nav items, and the toggle -- nothing else. Every piece of visual
interest on the site lives in the homepage photograph, which does not appear
inside `/blog/`. The slab serif helped, but the underlying question is untouched.

The tension is real: `DESIGN.md` records "Elevation & Depth: None" and identity
carried by the photograph plus one derived red, so adding a chrome element means
adding a voice. Either give the chrome one deliberate element and record the
rationale, or accept the asceticism knowingly. Do not smuggle styling in through
the body font.

## CI guard for stylesheet token drift

`src/style.css` and `hugo-src/assets/css/theme.css` hold independent copies of
the `--pg-*` tokens, plus a `@media (prefers-color-scheme: dark)` duplicate of
the dark palette in each. Nothing enforces that they agree.

A `tools/ci` check could parse both files and fail when the token sets diverge.
Priority unknown. Until it exists, the sync rule is a convention only.

## Static-page drift after the homepage recomposition

The static pages do not match the homepage's conventions:

- `src/links.html` still centres its list and carries no `[*]` markers.
- `src/credits.html` uses a `.brand` rule of `1.5em` / `normal`, while
  `src/links.html` uses `1.5rem` / `400` / `-0.01em`.

The `font-family: sans-serif` overrides in both pages have been removed; the
typeface is now declared once in `src/style.css` and inherited everywhere.

The reading register is settled: blog prose is Bitter (see `DESIGN.md`). What
remains open here is only the static-page drift above.

`DESIGN.md` describes the settled system; these pages should be brought into
line with it, or the description should be narrowed.

## Image assets that are deliberately kept

`src/my-favorite-meme.jpg` and `src/pngdeity_files/extermination-of-evil-shoki.jpg`
are not referenced by any page and are kept on purpose. Do not propose them for
deletion or optimization.

`gladstone.jpg` and `gladstone_modified.jpg` were removed in `e6f77a0`, and
`favicon-16.png` with them (`favicon.ico` already embeds a 16x16 frame).

## `content/develop-drafts` has diverged from `main`

Not merely behind: the branch carries 13 content commits of its own (the Open
Source on Offense draft, Wikipedia editing, and revisions to `beets-config`,
`stealing-github-commits`, and `pressing-work`) while missing everything `main`
gained. It is a two-way merge, not a fast-forward resync.

All posts on the branch are `draft = true`, so merging it cannot publish
anything. Last synced with `main` at `b500c59`.
