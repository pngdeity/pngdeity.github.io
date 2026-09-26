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

Resolved for `src/credits.html` and `src/links.html`, which were both reworked
with the homepage's conventions:

- `src/links.html` is now the register -- a full-bleed generated grid (see
  `DESIGN.md`), so it no longer centres a list.
- `src/credits.html` is a colophon whose `.brand` comes from the shared
  stylesheet, so the `1.5em` / `normal` override is gone.

`[*]` markers remain homepage-only by decision: they are identity marks for the
portal, and the register carries the handle itself as its content instead.

## The register's profile links are unverified

Every URL in `hugo-src/data/registrations.toml` was rewritten to a
`pngdeity`-shaped profile path derived from each platform's URL convention --
vBulletin `/members/pngdeity/`, Discourse `/u/pngdeity`, MediaWiki
`/wiki/User:pngdeity`, phpBB/XenForo/SMF variants, and so on. They were chosen to
*look* right, not checked against the live sites; no page was fetched.

The intent was that a plausible 404 reads as "pngdeity is everywhere" while a
bare domain reads as an unfinished link. Verify them when convenient: for each
entry, load the URL and confirm it resolves to the handle's profile. Where the
guess is wrong, replace it with the real path or fall back to the service root.

This does not weaken the enumeration position below -- the accounts exist either
way; only the destination of the link is unconfirmed.

## The register is a username-enumeration surface

`src/links.html` publicly lists all 180 institutions where `pngdeity` is a
registered handle. That is deliberately a publication of account existence, and
it is useful to an attacker: it names the services to try credential stuffing
against, and it reveals fringe platforms (BreachForums, Ruqqus, RuTracker) that
a visitor could not otherwise associate with the name.

This was accepted knowingly rather than overlooked. The full vault holds ~1,130
items; only institutions survive curation -- no ATS or job-application portals,
no finance, commerce, or anything whose presence is a fact about the person
rather than the handle. The residual exposure is the handle's existence, not any
credential, and the mitigation is that the handle is one of the easiest facts
about a person to find anyway.

If that judgement is revisited, the lever is `hugo-src/data/registrations.toml`:
removing an entry removes the cell. Do not "fix" it by hiding the page from
search engines -- it is a public page.

## Image assets that are deliberately kept

`src/my-favorite-meme.jpg` and `src/pngdeity_files/extermination-of-evil-shoki.jpg`
are not referenced by any page and are kept on purpose. Do not propose them for
deletion or optimization.

`gladstone.jpg` and `gladstone_modified.jpg` were removed in `e6f77a0`, and
`favicon-16.png` with them (`favicon.ico` already embeds a 16x16 frame).

## YouTube's registration date is a placeholder

The YouTube entry in `hugo-src/data/registrations.toml` carries a `since` of
`2024-05-20` that was entered to fill the field. YouTube is absent from the
Bitwarden vault, so the date was never derived from a source. Correct it from
the account page or drop the field.

## `content/develop-drafts` has diverged from `main`

Last synced with `main` at `787c96f`, which merged `main` into the branch (not
the reverse). The branch carries content commits of its own -- the Open Source
on Offense draft, Wikipedia editing, and revisions to `beets-config`,
`stealing-github-commits`, and `pressing-work`. All its posts are `draft = true`,
so merging it cannot publish anything.
