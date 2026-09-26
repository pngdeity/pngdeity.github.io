---
name: pngdeity
version: "alpha"
description: "The visual identity of pngdeity.ru, a personal site with a static half and a Hugo blog half that share one theming convention. Deadpan Utility: Digital Brutalism and post-internet surrealism."
colors:
  primary: "#9e141d"
  light:
    bg: "#f4f4f4"
    fg: "#1a1a1a"
    fg-muted: "#5a5a5a"
    surface: "#ffffff"
    border: "#d0d0d0"
    link: "#9e141d"
    link-hover: "#bd3c44"
    header-bg: "#111111"
    header-fg: "#f4f4f4"
  dark:
    bg: "#111111"
    fg: "#e6e6e6"
    fg-muted: "#a8a8a8"
    surface: "#1b1b1b"
    border: "#333333"
    link: "#ff6f5e"
    link-hover: "#ff9c90"
    header-bg: "#000000"
    header-fg: "#f4f4f4"
typography:
  body:
    fontFamily: "system-ui, -apple-system, \"Segoe UI\", Roboto, \"Helvetica Neue\", Arial, sans-serif"
    fontSize: "18px"
    fontWeight: 400
    lineHeight: 1.6
  brand:
    fontSize: "1.5rem"
    fontWeight: 400
    letterSpacing: "-0.01em"
  mono:
    fontFamily: "ui-monospace, SFMono-Regular, Menlo, Consolas, monospace"
  serif:
    fontFamily: "Bitter, Georgia, serif"
    fontSize: "18px"
    fontWeight: 400
    lineHeight: 1.6
spacing:
  page-margin: "40px"
  page-max-width: "800px"
  grid-gap: "2.5rem"
  grid-gap-narrow: "1.75rem"
  breakpoint-two-column: "48em"
rounded:
  toggle: "1.125rem"
components:
  theme-toggle:
    backgroundColor: "{colors.light.surface}"
    textColor: "{colors.light.fg}"
    rounded: "50%"
    size: "2.25rem"
  skip-link:
    backgroundColor: "{colors.light.surface}"
    textColor: "{colors.light.link}"
  post-card:
    backgroundColor: "{colors.light.surface}"
    textColor: "{colors.light.fg-muted}"
    rounded: "0.25rem"
    padding: "2rem"
  tag-chip:
    backgroundColor: "{colors.light.surface}"
    textColor: "{colors.light.link}"
    rounded: "1.5rem"
    padding: "0.5rem 1rem"
---

## Overview

pngdeity is a personal site belonging to Nathan Somers. It publishes a spare
homepage, a short link hub, a 404 page, and a blog built with Hugo.

The design position is **Deadpan Utility**: a synthesis of **Digital Brutalism**
and **post-internet surrealism**. The interface is a stark, machined container
that presents absurd, confrontational, or heavy material with total, unwavering
seriousness. It does not coddle the visitor, and it does not decorate.

The site rejects modern Web 2.0 aesthetics, soft padding, and corporate
minimalism. What replaces them is not tastefulness but **mechanical exactness**:
strict semantic rigour and absolute alignment, with borders and spacing doing the
work that shadows and ornament do elsewhere. Ornament carries no information, so
it is omitted -- not because restraint is elegant, but because decoration would
be a lie about what the page is.

Three principles follow from the position and govern new work:

**The Bowtie Principle (deadpan surrealism).** Bizarre elements are presented as
entirely factual and ordinary. The homepage photograph -- a black cat in a red
polka-dot bow tie at an open ThinkPad -- is the founding instance and dictates
the rule: a monochrome environment punctuated by a single, vivid, absurd
disruption. The bow tie is that disruption, and the accent colour is derived from
it rather than chosen.

**Active Interrogation.** The site does not apologise for errors or missing
data. The 404 page is the model: rather than a helpful message, it delivers
Paalen's question -- "it has become the role of the painting to look at the
spectator and ask him: What do you represent?" A broken link becomes a
confrontation rather than a dead end.

**No gimmicks.** The aesthetic is not simulated. There are no fake terminal
windows, no green-on-black typing effects, no decorative "system" chrome. The
machined quality comes from real structure -- alignment, borders, and type --
not from theatre about machines.

This document is the single source of truth for the site's visual identity. It
records what is settled, and nothing else: open questions live in `TODO.md` so
that what is written here can be relied on. Where this document and the code
disagree, the disagreement is a bug in one of them, not a licence to pick a
favourite.

Two kinds of entry appear here, and they bind differently. Recorded *rationale*
explains why a decision was made, so that a future change extends it rather than
contradicting it. Recorded *prohibitions* are decisions that were already tested
and rejected -- re-proposing one requires new evidence, not a fresh preference.

The document is valid against the DESIGN.md standard and can be checked with
`npx -p @google/design.md designmd lint DESIGN.md`. It lints with **zero errors**;
the `orphaned-tokens` warnings it reports are structural artefacts and should
not be "fixed". That rule counts a colour as used only when a *component*
references it, and the page background, body foreground, header band, and border
are not components in this schema. Silencing them would mean inventing components
that do not exist. The component entries below describe UI objects that are
genuinely on the site.

The brand is the handle `pngdeity`, always lowercase. It is a username, and the
name is written the way its owner writes it -- never capitalised, never set in
display capitals, never decorated. Lowercase is not a style choice; it is
accuracy.

The name means what it looks like it means. PNG is a lossless image format and a
deity is a god: "I'm a digital god, or perhaps only an image of one." That pun is
deliberately left to do its own work. It is never illustrated, captioned, or
explained on the site.

The recurring motif is **the image that looks back**. Three existing things
express it and were not invented to fit:

- The favicon is a ring with a pupil and a catchlight -- an eye, and equally a
  portal. Its geometry is derived from the homepage photograph.
- The homepage photograph puts a living subject in frame and on display.
- The 404 page quotes Wolfgang R. Paalen, *Form and Sense*: "it has become the
  role of the painting to look at the spectator and ask him: What do you
  represent?"

New work may extend this motif. It may not contradict it, and it may not
decorate it.

## Colors

Two palettes, always selected through a `data-theme` attribute of `light` or
`dark` on the `<html>` element, with `prefers-color-scheme` as the fallback when
no choice is stored, and the choice persisted in `localStorage` under the
`pg-theme` key.

The tokens, in both palettes:

| Token | Light | Dark |
| --- | --- | --- |
| `--pg-bg` | `#f4f4f4` | `#111111` |
| `--pg-fg` | `#1a1a1a` | `#e6e6e6` |
| `--pg-fg-muted` | `#5a5a5a` | `#a8a8a8` |
| `--pg-surface` | `#ffffff` | `#1b1b1b` |
| `--pg-border` | `#d0d0d0` | `#333333` |
| `--pg-link` | `#9e141d` | `#ff6f5e` |
| `--pg-link-hover` | `#bd3c44` | `#ff9c90` |
| `--pg-header-bg` | `#111111` | `#000000` |
| `--pg-header-fg` | `#f4f4f4` | `#f4f4f4` |

`--pg-header-bg` and `--pg-header-fg` exist only on the blog half, where the
theme paints a dark header and footer band. **That band stays dark and its text
stays light in both themes**; the header is not restyled to follow the page
foreground. On the static half these two tokens are absent rather than
duplicated.

Accent colours are derived from the homepage photograph, never chosen by taste.
The bow tie is the one saturated colour in an otherwise grey and black image, and
its lit face measures at hue `356°` / saturation `87%`. Both accents are that red
with the *value* set for their background: dark mode **lightens** it so it glows on
`#111111`, light mode **darkens** it so it reads on `#f4f4f4`. The favicon follows
the same rule -- coral on its dark tile, `#9e141d` on its light tile. Colour choices
follow the photograph, not the reverse. If a future palette change is needed,
extend this rule rather than inventing a colour.

## Typography

Type is organised into two registers.

**Chrome** -- navigation, section labels, metadata, share links, timestamps, the
theme toggle, and the `[*]` markers -- uses the platform UI stack (`system-ui`,
`-apple-system`, `Segoe UI`, `Roboto`, `Helvetica Neue`, `Arial`, `sans-serif`)
at `18px` with `line-height: 1.6`. The same stack is declared once in the static
half (`src/style.css`) and applied to the blog body in the blog half via the
`.pg-site` class; the auxiliary pages inherit it rather than declaring their own.
Monospace is used for the `[*]` link markers only.

**Blog prose** -- article body copy and its headings -- uses **Bitter**, a slab
serif, self-hosted from `hugo-src/assets/fonts/bitter/` (three woff2 faces: 400,
700, and 400 italic), licensed under the SIL Open Font License 1.1. Body copy is
set at weight `400`; headings at `600`, since `700` is heavy at display sizes.
The register stops at prose: chrome keeps the system stack, so the two are
visibly distinct. The rules live at the end of `hugo-src/assets/css/theme.css`
so that they win the cascade, and heading selectors must carry class specificity
(`article h1.athelas`, not `article h1`) or the `athelas` utility in
`utilities.css` takes the font back.

Bitter was chosen on legibility rather than personality. In a comparison against
a transitional (Source Serif 4), an old-style (EB Garamond), and a book face
(Crimson Pro), it was the most readable at body size -- thick uniform strokes,
large x-height, low stroke contrast -- and the only candidate that held its own
against the dark header band. The sturdiness is the point: prose gets a serif so
it reads *as prose*, and legibility serves that better than elegance.

Headings set `line-height: 1.2`. The brand mark stays in the system stack at
`1.5rem` weight `400` with `letter-spacing: -0.01em`, in `--pg-fg`. It is plain
rather than unemphatic: a username is not a logo, so it is set as text and left
to sit beside the photograph without competing for attention. It is not
decorated, capitalised, or given a mark of its own.

Webfonts were once ruled out here on the grounds that a personal site
downloading a typeface "has misjudged its own proportions". That was never a
strong position; the real cost is a few hundred KB per visit and the obligation
to self-host and carry a licence. Where a second register earns its place, it is
worth paying -- as it is for blog prose. The default remains the system stack.

## Layout

Prose pages are single-column: `max-width: 800px`, centred with `40px` of
vertical margin, `line-height: 1.6`.

The homepage is the exception. At `48em` and wider it is a two-column grid with
the photograph on the left at `1.6fr` and the text stack on the right at `1fr`,
separated by `2.5rem` and vertically centred. Below the breakpoint it collapses
to one column with a `1.75rem` gap. The photograph leads at every width.

The reason is compositional: the photograph is asymmetric, with its subject
left-of-centre gazing rightward into the frame. Placed in a left column, that
gaze travels across the page and lands on the text. Centring the image instead
aims the subject's attention at nothing. The layout follows the photograph.

The homepage is deliberately brief -- an image, the name, a caption, and six
links -- because brevity funnels visitors into the content pages. It must not
grow a bio, a tagline, an "about" block, or a footer. When tempted to add
something to the homepage, the correct answer is usually a different page.

## Elevation & Depth

None. The site has no shadows, no overlays, and no layered surfaces. Depth is
ornament, and ornament carries no information: a border or a spacing change
states every boundary the site needs, so a shadow would be a decoration with no
job. This follows from the position rather than from restraint -- depth is not
*too much*, it is *beside the point*.

## Shapes

Restrained. The one distinct shape is the theme toggle: a fixed circular button
in the top-right corner, `2.25rem` across, with a `50%` radius. It shows the
theme it will switch *to* -- a moon in light mode, a sun in dark mode -- so the
control always describes its own action. Small radii elsewhere come from the
theme's utility classes, not from this design.

## Components

**The `[*]` marker.** Bracketed asterisks introduce navigation links. It is the
one piece of the site's earlier visual language that survived, and it is
retained as identity rather than utility. On the homepage it is injected by CSS
(`.portal a::before`) rather than written into the markup, set in monospace at
`0.9em` in `--pg-fg-muted` with `0.6em` of trailing space, and the list items use
a negative `text-indent` so the markers form a clean column of their own with
the link labels aligned to a second column. It reads as display type: an index,
not prose.

**The theme toggle.** A button showing the theme it will switch to, kept in sync
by `theme.js` through `aria-pressed` and `aria-label` ("Switch to dark theme" /
"Switch to light theme"). It is the only interactive chrome on the site.

**The post card.** The blog's list and term pages render each post as a card: a
surface background, the section label, the title, the summary, and a "read more"
link. Cards are separated by a bottom border rather than a shadow, consistent
with the no-depth rule.

**The tag chip.** Post tags render as pill-shaped chips in `--pg-link` on a
surface background, with `1.5rem` of radius. They are the only place the accent
appears as a fill rather than as text.

**The skip link.** Every page begins with a visually hidden "Skip to content"
link that reveals on focus and targets `#main-content`.

## Declined

Proposals that were considered and rejected. They are recorded so they are not
silently re-proposed; revisit only with new evidence.

**Replacing the colour tokens with a scarlet family.** An external design
proposal (`ethos.md`, working material, untracked) put forward
`--pg-link: #FF0033` / `#CC0029` / `#FF3333` / `#FF6666`. Declined: the accent is
a *sampled* value from the photograph (hue `356°`, saturation `87%`, a crimson),
and these are unchecked scarlet. Adopting them would undo the derivation rule
above.

**Absolute white and black.** The same proposal argued for `--pg-bg: #FFFFFF` /
`--pg-fg: #000000` and, in dark mode, `#000000` / `#E6E6E6`. Declined: the
calibrated `#f4f4f4` / `#1a1a1a` pair is deliberately softer for `800px` of
prose, and pure black on pure white reads as glare rather than as severity.

**A `#000000` light-mode border.** Argued as "maximum contrast, unyielding
borders". Declined: a 1px pure-black rule on white is a readability problem, not
a principle. Borders are `--pg-border`, which is tuned per theme.

**A single-layer theming block.** The proposal defined only
`html[data-theme="light"]` and `html[data-theme="dark"]`, with no
`@media (prefers-color-scheme: dark)`. Declined: the site ships three layers so
that a visitor who has never touched the toggle still follows their operating
system. Adopting the block verbatim would have broken that, and it omitted the
blog-only `--pg-header-bg` / `--pg-header-fg` pair as well.

**Stripping the homepage of all branding.** "Stripped of all branding and logos."
Declined: the homepage carries the `pngdeity` wordmark deliberately, and the
brand is the site's identity rather than decoration.

**A serif for prose plus a machined monospace for all structural chrome**
(Century Schoolbook / Charter / Clarendon; JetBrains Mono / Iosevka / Berkeley
Mono). The *direction* survived -- there are two registers -- but the specific
faces did not. Every candidate named is a webfont that is not a system font on
Linux or Windows, and the chrome register is the platform UI stack, not a
monospace. See Typography.

**Buttons are disallowed; links are plaintext only.** Declined as written: the
theme toggle is a button, deliberately, because it is a control with state
rather than navigation.

**Editorial or midcentury-modern framing.** The site was not designed as a
tasteful printed page, and that framing is not to be reintroduced. Declined
phrasings include "a composed document", "quiet and typographic", and treating
elegance as the goal. The position is deadpan utility, not print craft: the site
is stark by intent rather than refined by temperament, and where it is plain that
is because a username is not a logo, not because restraint is stylish. An
argument that reasons from taste is reasoning from the wrong premise.

**The 80/20 aluminium extrusion metaphor.** An earlier text described the layout
as resembling "industrial metrology tools or 80/20 aluminum extrusions". Declined
as overwrought. The mechanical-exactness principle itself is real and is recorded
in the Overview; the comparison is not, and must not be carried into design
writing.

## Do's and Don'ts

**Do**

- Write `pngdeity` lowercase, always, in every context including titles.
- Change a token in **both** stylesheets. The static half (`src/style.css`) and
  the blog half (`hugo-src/assets/css/theme.css`) hold independent copies
  of the `--pg-*` set; there is no shared source. Editing one alone is the single
  most likely way to damage this site.
- Take colour cues from the homepage photograph, as both accents do.
- Keep the homepage brief, and keep the header and footer band dark with light
  text in both themes.
- Treat the photograph as the page's one loud element: it is the bow tie of the
  composition, and everything else is the ground it sits on.
- Keep the reading register confined to blog prose. Chrome stays on the system
  stack; a serif in the navigation or metadata would collapse the two registers
  into one undifferentiated voice.
- Carry the font's licence file when adding or updating a self-hosted face, and
  keep the face subset to what the content needs.
- Argue from the principles in the Overview -- deadpan, mechanical exactness,
  interrogation -- rather than from taste. If a change can only be justified as
  "it looks nicer", it is not justified.

**Don't**

- Capitalise or decorate the brand: no `PNGDEITY`, no letterspaced capitals, no
  logo veneer, no halo, no lens flare, no icon after the word.
- Invent a tagline or slogan. A meaningless phrase is worse than no phrase; the
  previous site title `The End of the Internet` was an accident and is retired,
  not replaced with another one.
- Explain the `pngdeity` pun on the site. If it needs explaining, it is dead.
- Let the word `deity` drift toward grandeur. The name is a joke its owner is
  fond of; treating it solemnly turns a good joke into a bad brand.
- Add a webfont beyond the Bitter reading register, add shadows, or add depth.
- Set headings at weight `700` in the reading register; `600` is the ceiling.
- Add copy to the homepage.
- Simulate the aesthetic: no fake terminal windows, no typing effects, no
  decorative "system" chrome. The machined quality must come from real structure.
- Reason from elegance, taste, or editorial craft. That framing was declined.
