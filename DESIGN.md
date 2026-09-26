---
name: pngdeity
version: "alpha"
description: The visual identity of pngdeity.ru, a personal site with a static half and a Hugo blog half that share one theming convention.
colors:
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
    fontFamily: "-apple-system, BlinkMacSystemFont, \"Segoe UI\", Roboto, Helvetica, Arial, sans-serif"
    fontSize: "18px"
    fontWeight: 400
    lineHeight: 1.6
  brand:
    fontSize: "1.5rem"
    fontWeight: 400
    letterSpacing: "-0.01em"
  mono:
    fontFamily: "ui-monospace, SFMono-Regular, Menlo, Consolas, monospace"
spacing:
  page-margin: "40px"
  page-max-width: "800px"
  grid-gap: "2.5rem"
  grid-gap-narrow: "1.75rem"
  breakpoint-two-column: "48em"
rounded:
  toggle: "50%"
components:
  theme-toggle:
    backgroundColor: "{colors.light.surface}"
    textColor: "{colors.light.fg}"
    rounded: "{rounded.toggle}"
    size: "2.25rem"
  skip-link:
    backgroundColor: "{colors.light.surface}"
    textColor: "{colors.light.link}"
---

## Overview

pngdeity is a personal site belonging to Nathan Somers. It publishes a spare
homepage, a short link hub, a 404 page, and a blog built with Hugo. The design
intent is quiet and typographic: a page should read as a composed document, not
as a dressed-up template.

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

The body face is the platform UI stack (`-apple-system`, `BlinkMacSystemFont`,
`Segoe UI`, `Roboto`, `Helvetica`, `Arial`, `sans-serif`) at `18px` with
`line-height: 1.6`. No webfont is loaded. The site must not acquire one: a
personal site that downloads a typeface to say six words has misjudged its own
proportions.

Headings set `line-height: 1.2`. The brand mark is `1.5rem` at weight `400` with
`letter-spacing: -0.01em`, in `--pg-fg`. It is deliberately unemphatic -- the
photograph is the loud element on the page, and the name does not compete with
it.

Monospace type is used for the `[*]` link markers only. See Components.

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

None. The site has no shadows, no overlays, and no layered surfaces. Depth would
compete with the photograph and add nothing that a border or a spacing change
cannot say.

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

**The skip link.** Every page begins with a visually hidden "Skip to content"
link that reveals on focus and targets `#main-content`.

## Do's and Don'ts

**Do**

- Write `pngdeity` lowercase, always, in every context including titles.
- Change a token in **both** stylesheets. The static half (`src/style.css`) and
  the blog half (`hugo-src/assets/ananke/css/theme.css`) hold independent copies
  of the `--pg-*` set; there is no shared source. Editing one alone is the single
  most likely way to damage this site.
- Take colour cues from the homepage photograph, as both accents do.
- Keep the homepage brief, and keep the header and footer band dark with light
  text in both themes.
- Treat the photograph as the page's one loud element and let type stay quiet.

**Don't**

- Capitalise or decorate the brand: no `PNGDEITY`, no letterspaced capitals, no
  logo veneer, no halo, no lens flare, no icon after the word.
- Invent a tagline or slogan. A meaningless phrase is worse than no phrase; the
  previous site title `The End of the Internet` was an accident and is retired,
  not replaced with another one.
- Explain the `pngdeity` pun on the site. If it needs explaining, it is dead.
- Let the word `deity` drift toward grandeur. The name is a joke its owner is
  fond of; treating it solemnly turns a good joke into a bad brand.
- Load a webfont, add shadows, or add depth.
- Add copy to the homepage.
