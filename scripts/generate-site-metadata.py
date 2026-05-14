#!/usr/bin/env python3
import json
import os
import re
import subprocess
from datetime import datetime, timezone
from pathlib import Path

REPO_ROOT = Path.cwd()
SRC_ROOT = REPO_ROOT / "src"
CANONICAL_ORIGIN = "https://pngdeity.ru"


def walk_files(dir_path: Path) -> list[Path]:
    if not dir_path.exists():
        return []
    files: list[Path] = []
    for root, _, names in os.walk(dir_path):
        root_path = Path(root)
        for name in names:
            full_path = root_path / name
            if full_path.is_file():
                files.append(full_path)
    return files


def to_posix_path(file_path: Path) -> str:
    return file_path.relative_to(SRC_ROOT).as_posix()


def is_sitemap_html(rel_path: str) -> bool:
    if not rel_path.endswith(".html"):
        return False
    if rel_path == "404.html":
        return False
    top = rel_path.split("/", 1)[0]
    return "/" not in rel_path or top in {"blog", "app"}


def file_path_to_url(rel_path: str) -> str:
    if rel_path == "index.html":
        return f"{CANONICAL_ORIGIN}/"
    if rel_path.endswith("/index.html"):
        return f"{CANONICAL_ORIGIN}/{rel_path[: -len('index.html')]}"
    return f"{CANONICAL_ORIGIN}/{rel_path}"


def xml_escape(value: str) -> str:
    return (
        value.replace("&", "&amp;")
        .replace("<", "&lt;")
        .replace(">", "&gt;")
        .replace('"', "&quot;")
        .replace("'", "&apos;")
    )


def get_lastmod(file_path: Path) -> str:
    try:
        git_date = subprocess.run(
            ["git", "log", "-1", "--format=%cI", "--", str(file_path)],
            cwd=REPO_ROOT,
            capture_output=True,
            text=True,
            check=False,
        ).stdout.strip()
        if git_date:
            return git_date
    except Exception:
        pass
    mtime = file_path.stat().st_mtime
    return (
        datetime.fromtimestamp(mtime, tz=timezone.utc)
        .isoformat(timespec="milliseconds")
        .replace("+00:00", "Z")
    )


def strip_tags(value: str) -> str:
    return re.sub(r"\s+", " ", re.sub(r"<[^>]*>", " ", value)).strip()


def extract_metadata(file_path: Path) -> dict[str, str]:
    try:
        html = file_path.read_text(encoding="utf-8")
        title_match = re.search(r"<title[^>]*>([\s\S]*?)</title>", html, flags=re.IGNORECASE)
        h1_match = re.search(r"<h1[^>]*>([\s\S]*?)</h1>", html, flags=re.IGNORECASE)
        title = strip_tags(title_match.group(1)) if title_match else ""
        heading = strip_tags(h1_match.group(1)) if h1_match else ""
        return {"title": title, "heading": heading}
    except Exception:
        return {"title": "", "heading": ""}


def fallback_name(rel_path: str) -> str:
    base = re.sub(r"[-_]+", " ", Path(rel_path).stem)
    return re.sub(r"\b\w", lambda m: m.group(0).upper(), base) if base else rel_path


def write_file(file_path: Path, content: str) -> None:
    file_path.parent.mkdir(parents=True, exist_ok=True)
    file_path.write_text(content, encoding="utf-8")


if not SRC_ROOT.exists():
    raise RuntimeError(f"Missing source directory: {SRC_ROOT}")

html_files = [
    {"full_path": full_path, "rel_path": to_posix_path(full_path)}
    for full_path in walk_files(SRC_ROOT)
]
html_files = [entry for entry in html_files if is_sitemap_html(entry["rel_path"])]
html_files.sort(key=lambda e: e["rel_path"])

robots_content = "\n".join(
    ["User-agent: *", "Allow: /", f"Sitemap: {CANONICAL_ORIGIN}/sitemap.xml", ""]
)
write_file(SRC_ROOT / "robots.txt", robots_content)

url_entries: list[dict[str, str | Path]] = []
for entry in html_files:
    full_path = entry["full_path"]
    rel_path = entry["rel_path"]
    url_entries.append(
        {
            "loc": file_path_to_url(rel_path),
            "lastmod": get_lastmod(full_path),
            "rel_path": rel_path,
            "full_path": full_path,
        }
    )

sitemap_lines = [
    '<?xml version="1.0" encoding="UTF-8"?>',
    '<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">',
]
for entry in url_entries:
    sitemap_lines.append(
        f"  <url><loc>{xml_escape(entry['loc'])}</loc>"
        f"<lastmod>{xml_escape(entry['lastmod'])}</lastmod></url>"
    )
sitemap_lines.extend(["</urlset>", ""])
write_file(SRC_ROOT / "sitemap.xml", "\n".join(sitemap_lines))

root_pages = [e for e in url_entries if "/" not in e["rel_path"]]
blog_pages = [e for e in url_entries if e["rel_path"].startswith("blog/")]
app_pages = [e for e in url_entries if e["rel_path"].startswith("app/")]

discovered_sections: set[str] = set()
for page in blog_pages:
    rest = page["rel_path"][len("blog/") :]
    section = rest.split("/", 1)[0]
    if section and section != "index.html":
        discovered_sections.add(section)

metadata_samples = []
for entry in url_entries[:200]:
    metadata_samples.append(
        {
            "rel_path": entry["rel_path"],
            **extract_metadata(entry["full_path"]),
        }
    )

key_pages = "\n".join(
    [
        f"- {(sample['title'] or sample['heading'] or fallback_name(sample['rel_path']))} "
        f"({file_path_to_url(sample['rel_path'])})"
        for sample in metadata_samples[:12]
    ]
)

site_overview = (
    "Site generated from static HTML plus Hugo blog content and an optional Blazor app "
    f"section under /app/. Canonical host: {CANONICAL_ORIGIN}/."
)
main_sections = "\n".join(
    [
        "- Root static pages (landing and standalone HTML pages)",
        "- Blog content under /blog/",
        "- Application assets/content under /app/"
        if len(app_pages) > 0
        else "- /app/ section is currently not populated",
    ]
)
content_types = "\n".join(
    [
        "- HTML pages",
        "- Downloadable static assets (documents, images, and related files)",
        "- Generated application files (when /app/ artifact is present)",
    ]
)
attribution = "- License file present in repository: LICENSE.md (see project root)."
contact = "- Security contact: mailto:security@pngdeity.ru"

llms_content = "\n".join(
    [
        "# Site Overview",
        site_overview,
        "",
        "# Main Sections",
        main_sections,
        "",
        "# Content Types",
        content_types,
        "",
        "# Attribution / Licensing",
        attribution,
        "",
        "# Contact",
        contact,
        "",
    ]
)
write_file(SRC_ROOT / "llms.txt", llms_content)

full_root_pages = (
    "\n".join([f"- {page['rel_path']} -> {file_path_to_url(page['rel_path'])}" for page in root_pages[:50]])
    or "- None found"
)
full_sections = (
    "\n".join([f"- {section}" for section in sorted(discovered_sections)]) or "- No nested Hugo sections discovered"
)

llms_full_content = "\n".join(
    [
        "# Site Overview",
        site_overview,
        "",
        "# Main Sections",
        main_sections,
        "",
        "# Content Types",
        content_types,
        "",
        "# Hugo Major Sections",
        full_sections,
        "",
        "# Key Static Pages (src/)",
        full_root_pages,
        "",
        "# Representative Page Samples",
        key_pages or "- No HTML pages discovered for sampling",
        "",
        "# Application Section",
        f"- /app/ detected with {len(app_pages)} HTML page(s) in this build output."
        if len(app_pages) > 0
        else "- /app/ not detected or contains no HTML pages in this build output.",
        "",
        "# Attribution / Licensing",
        attribution,
        "",
        "# Contact",
        contact,
        "",
    ]
)
write_file(SRC_ROOT / "llms-full.txt", llms_full_content)

webfinger = {
    "subject": "acct:pngdeity@pngdeity.ru",
    "aliases": ["https://pngdeity.ru/"],
    "links": [
        {
            "rel": "self",
            "type": "application/activity+json",
            "href": "https://pngdeity.ru/",
        }
    ],
}

write_file(SRC_ROOT / ".well-known" / "webfinger", f"{json.dumps(webfinger, indent=2)}\n")
write_file(SRC_ROOT / ".well-known" / "keybase.txt", "# Keybase proof placeholder\n")

security_txt = "\n".join(
    [
        "Contact: mailto:security@pngdeity.ru",
        "Expires: 2026-05-14",
        f"Policy: {CANONICAL_ORIGIN}/",
        "",
    ]
)
write_file(SRC_ROOT / ".well-known" / "security.txt", security_txt)

print(f"Generated metadata files in {SRC_ROOT}")
print(f"Sitemap entries: {len(url_entries)}")
