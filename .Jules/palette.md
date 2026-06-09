## 2024-03-21 - Custom SVG Accessibility
**Learning:** Custom SVG data visualizations like contribution graphs aren't intrinsically accessible. While `<title>` provides a tooltip for mouse users, keyboard and screen reader users completely miss the data unless the individual data points (like `<rect>`) are given `tabindex="0"` and an explicit `aria-label`. Moreover, a global `role="img"` and `aria-label` on the parent `<svg>` is necessary to give users context.
**Action:** When building or reviewing custom SVG graphics that represent data, always ensure interactive or informative points have manual focus management (`tabindex="0"`, `:focus-visible` styles) and screen-reader readable properties (`aria-label`).

## 2024-03-21 - Basic HTML Accessibility Overlooked
**Learning:** Even simple, static HTML pages often lack fundamental accessibility attributes like the `lang` attribute on the `<html>` tag and `alt` text on primary images. Without a `lang` attribute, screen readers cannot properly determine the language and pronunciation rules, severely degrading the experience.
**Action:** Always verify that the root HTML structures (`<html>`, `<img>`) have the necessary accessibility attributes, even on simple landing pages that don't use complex frameworks.## 2026-06-09 - iframe Accessibility
**Learning:** Simple HTML files often omit accessibility attributes on embedded content like `iframe`s. Without a `title` attribute, screen readers cannot properly announce the purpose of the iframe to users before they navigate into it.
**Action:** When working with or reviewing pages with embedded content, ensure `iframe` tags always have a descriptive `title` attribute that explains the content or purpose of the embed.
