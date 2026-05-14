#!/usr/bin/env python3
import re
import subprocess
import sys
from pathlib import Path
from urllib.parse import urlparse

SRC_ROOT = Path.cwd() / "src"
REQUIRED_ROBOTS_LINES = [
    "User-agent: *",
    "Allow: /",
    "Sitemap: https://pngdeity.ru/sitemap.xml",
]
errors: list[str] = []


def fail(message: str) -> None:
    errors.append(message)


def ensure_exists(file_path: Path, label: str) -> bool:
    if not file_path.exists():
        fail(f"{label} is missing: {file_path}")
        return False
    return True


sitemap_path = SRC_ROOT / "sitemap.xml"
robots_path = SRC_ROOT / "robots.txt"
llms_path = SRC_ROOT / "llms.txt"
llms_full_path = SRC_ROOT / "llms-full.txt"
wk_webfinger = SRC_ROOT / ".well-known" / "webfinger"
wk_keybase = SRC_ROOT / ".well-known" / "keybase.txt"
wk_security = SRC_ROOT / ".well-known" / "security.txt"

if ensure_exists(sitemap_path, "sitemap.xml"):
    try:
        xmllint = subprocess.run(
            ["xmllint", "--noout", str(sitemap_path)],
            capture_output=True,
            text=True,
            check=False,
        )
        if xmllint.returncode != 0:
            stderr = (xmllint.stderr or "").strip()
            fail(f"sitemap.xml is not valid XML: {stderr or 'xmllint returned non-zero status'}")
    except Exception as exc:
        fail(f"xmllint failed to execute: {exc}")

    sitemap_content = sitemap_path.read_text(encoding="utf-8")
    loc_matches = [m.strip() for m in re.findall(r"<loc>([^<]+)</loc>", sitemap_content)]
    for loc in loc_matches:
        parsed = urlparse(loc)
        if not parsed.scheme or not parsed.netloc:
            fail(f"Invalid sitemap URL: {loc}")
            continue
        if parsed.scheme != "https":
            fail(f"Sitemap URL is not HTTPS: {loc}")
        if parsed.hostname != "pngdeity.ru":
            fail(f"Sitemap URL is not canonical domain pngdeity.ru: {loc}")

if ensure_exists(robots_path, "robots.txt"):
    robots = robots_path.read_text(encoding="utf-8")
    for required_line in REQUIRED_ROBOTS_LINES:
        if required_line not in robots:
            fail(f"robots.txt missing required line: {required_line}")

ensure_exists(wk_webfinger, ".well-known/webfinger")
ensure_exists(wk_keybase, ".well-known/keybase.txt")
ensure_exists(wk_security, ".well-known/security.txt")

if ensure_exists(llms_path, "llms.txt") and not llms_path.read_text(encoding="utf-8").strip():
    fail("llms.txt is empty")
if ensure_exists(llms_full_path, "llms-full.txt") and not llms_full_path.read_text(encoding="utf-8").strip():
    fail("llms-full.txt is empty")

if errors:
    print("Metadata validation failed:", file=sys.stderr)
    for error in errors:
        print(f"- {error}", file=sys.stderr)
    sys.exit(1)

print("Metadata validation passed.")
