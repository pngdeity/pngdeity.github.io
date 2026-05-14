#!/usr/bin/env node
import fs from 'node:fs';
import path from 'node:path';
import { execFileSync } from 'node:child_process';

const REPO_ROOT = process.cwd();
const SRC_ROOT = path.join(REPO_ROOT, 'src');
const CANONICAL_ORIGIN = 'https://pngdeity.ru';

function ensureDir(dirPath) {
  fs.mkdirSync(dirPath, { recursive: true });
}

function walkFiles(dirPath) {
  if (!fs.existsSync(dirPath)) return [];
  const out = [];
  for (const entry of fs.readdirSync(dirPath, { withFileTypes: true })) {
    const fullPath = path.join(dirPath, entry.name);
    if (entry.isDirectory()) {
      out.push(...walkFiles(fullPath));
    } else if (entry.isFile()) {
      out.push(fullPath);
    }
  }
  return out;
}

function toPosixPath(filePath) {
  return path.relative(SRC_ROOT, filePath).split(path.sep).join('/');
}

function isSitemapHtml(relPath) {
  if (!relPath.endsWith('.html')) return false;
  if (relPath === '404.html') return false;
  const top = relPath.split('/')[0];
  return !relPath.includes('/') || top === 'blog' || top === 'app';
}

function filePathToUrl(relPath) {
  if (relPath === 'index.html') return `${CANONICAL_ORIGIN}/`;
  if (relPath.endsWith('/index.html')) {
    return `${CANONICAL_ORIGIN}/${relPath.slice(0, -'index.html'.length)}`;
  }
  return `${CANONICAL_ORIGIN}/${relPath}`;
}

function xmlEscape(value) {
  return value
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&apos;');
}

function getLastMod(filePath) {
  try {
    const gitDate = execFileSync('git', ['log', '-1', '--format=%cI', '--', filePath], {
      cwd: REPO_ROOT,
      encoding: 'utf8',
      stdio: ['ignore', 'pipe', 'ignore'],
    }).trim();
    if (gitDate) return gitDate;
  } catch {
    // Fall back to mtime below.
  }
  const stat = fs.statSync(filePath);
  return stat.mtime.toISOString();
}

function stripTags(value) {
  return value
    .replace(/<[^>]*>/g, ' ')
    .replace(/\s+/g, ' ')
    .trim();
}

function extractMetadata(filePath) {
  try {
    const html = fs.readFileSync(filePath, 'utf8');
    const titleMatch = html.match(/<title[^>]*>([\s\S]*?)<\/title>/i);
    const h1Match = html.match(/<h1[^>]*>([\s\S]*?)<\/h1>/i);
    const title = titleMatch ? stripTags(titleMatch[1]) : '';
    const heading = h1Match ? stripTags(h1Match[1]) : '';
    return { title, heading };
  } catch {
    return { title: '', heading: '' };
  }
}

function fallbackName(relPath) {
  const base = path.basename(relPath, '.html').replace(/[-_]+/g, ' ');
  return base ? base.replace(/\b\w/g, (c) => c.toUpperCase()) : relPath;
}

function writeFile(filePath, content) {
  ensureDir(path.dirname(filePath));
  fs.writeFileSync(filePath, content, 'utf8');
}

if (!fs.existsSync(SRC_ROOT)) {
  throw new Error(`Missing source directory: ${SRC_ROOT}`);
}

const htmlFiles = walkFiles(SRC_ROOT)
  .map((fullPath) => ({ fullPath, relPath: toPosixPath(fullPath) }))
  .filter(({ relPath }) => isSitemapHtml(relPath))
  .sort((a, b) => a.relPath.localeCompare(b.relPath));

const robotsContent = [
  'User-agent: *',
  'Allow: /',
  `Sitemap: ${CANONICAL_ORIGIN}/sitemap.xml`,
  '',
].join('\n');
writeFile(path.join(SRC_ROOT, 'robots.txt'), robotsContent);

const urlEntries = htmlFiles.map(({ fullPath, relPath }) => {
  const url = filePathToUrl(relPath);
  return {
    loc: url,
    lastmod: getLastMod(fullPath),
    relPath,
    fullPath,
  };
});

const sitemapXml = [
  '<?xml version="1.0" encoding="UTF-8"?>',
  '<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">',
  ...urlEntries.map(
    ({ loc, lastmod }) =>
      `  <url><loc>${xmlEscape(loc)}</loc><lastmod>${xmlEscape(lastmod)}</lastmod></url>`,
  ),
  '</urlset>',
  '',
].join('\n');
writeFile(path.join(SRC_ROOT, 'sitemap.xml'), sitemapXml);

const rootPages = urlEntries.filter(({ relPath }) => !relPath.includes('/'));
const blogPages = urlEntries.filter(({ relPath }) => relPath.startsWith('blog/'));
const appPages = urlEntries.filter(({ relPath }) => relPath.startsWith('app/'));

const discoveredSections = new Set();
for (const { relPath } of blogPages) {
  const rest = relPath.slice('blog/'.length);
  const section = rest.split('/')[0];
  if (section && section !== 'index.html') discoveredSections.add(section);
}

const metadataSamples = urlEntries.slice(0, 200).map(({ fullPath, relPath }) => ({
  relPath,
  ...extractMetadata(fullPath),
}));

const keyPages = metadataSamples
  .slice(0, 12)
  .map(({ relPath, title, heading }) => {
    const label = title || heading || fallbackName(relPath);
    return `- ${label} (${filePathToUrl(relPath)})`;
  })
  .join('\n');

const siteOverview = `Site generated from static HTML plus Hugo blog content and an optional Blazor app section under /app/. Canonical host: ${CANONICAL_ORIGIN}/.`;
const mainSections = [
  '- Root static pages (landing and standalone HTML pages)',
  '- Blog content under /blog/',
  appPages.length > 0 ? '- Application assets/content under /app/' : '- /app/ section is currently not populated',
].join('\n');
const contentTypes = [
  '- HTML pages',
  '- Downloadable static assets (documents, images, and related files)',
  '- Generated application files (when /app/ artifact is present)',
].join('\n');
const attribution = '- License file present in repository: LICENSE.md (see project root).';
const contact = '- Security contact: mailto:security@pngdeity.ru';

const llmsContent = [
  '# Site Overview',
  siteOverview,
  '',
  '# Main Sections',
  mainSections,
  '',
  '# Content Types',
  contentTypes,
  '',
  '# Attribution / Licensing',
  attribution,
  '',
  '# Contact',
  contact,
  '',
].join('\n');
writeFile(path.join(SRC_ROOT, 'llms.txt'), llmsContent);

const fullRootPages = rootPages
  .slice(0, 50)
  .map(({ relPath }) => `- ${relPath} -> ${filePathToUrl(relPath)}`)
  .join('\n') || '- None found';
const fullSections = [...discoveredSections]
  .sort((a, b) => a.localeCompare(b))
  .map((s) => `- ${s}`)
  .join('\n') || '- No nested Hugo sections discovered';

const llmsFullContent = [
  '# Site Overview',
  siteOverview,
  '',
  '# Main Sections',
  mainSections,
  '',
  '# Content Types',
  contentTypes,
  '',
  '# Hugo Major Sections',
  fullSections,
  '',
  '# Key Static Pages (src/)',
  fullRootPages,
  '',
  '# Representative Page Samples',
  keyPages || '- No HTML pages discovered for sampling',
  '',
  '# Application Section',
  appPages.length > 0
    ? `- /app/ detected with ${appPages.length} HTML page(s) in this build output.`
    : '- /app/ not detected or contains no HTML pages in this build output.',
  '',
  '# Attribution / Licensing',
  attribution,
  '',
  '# Contact',
  contact,
  '',
].join('\n');
writeFile(path.join(SRC_ROOT, 'llms-full.txt'), llmsFullContent);

const webfinger = {
  subject: 'acct:pngdeity@pngdeity.ru',
  aliases: ['https://pngdeity.ru/'],
  links: [
    {
      rel: 'self',
      type: 'application/activity+json',
      href: 'https://pngdeity.ru/',
    },
  ],
};

writeFile(
  path.join(SRC_ROOT, '.well-known', 'webfinger'),
  `${JSON.stringify(webfinger, null, 2)}\n`,
);
writeFile(path.join(SRC_ROOT, '.well-known', 'keybase.txt'), '# Keybase proof placeholder\n');

const securityTxt = [
  'Contact: mailto:security@pngdeity.ru',
  'Expires: 2026-05-14',
  `Policy: ${CANONICAL_ORIGIN}/`,
  '',
].join('\n');
writeFile(path.join(SRC_ROOT, '.well-known', 'security.txt'), securityTxt);

console.log(`Generated metadata files in ${SRC_ROOT}`);
console.log(`Sitemap entries: ${urlEntries.length}`);
