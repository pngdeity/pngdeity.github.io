#!/usr/bin/env node
import fs from 'node:fs';
import path from 'node:path';
import { spawnSync } from 'node:child_process';

const SRC_ROOT = path.join(process.cwd(), 'src');
const REQUIRED_ROBOTS_LINES = [
  'User-agent: *',
  'Allow: /',
  'Sitemap: https://pngdeity.ru/sitemap.xml',
];

function fail(message) {
  errors.push(message);
}

function ensureExists(filePath, label) {
  if (!fs.existsSync(filePath)) {
    fail(`${label} is missing: ${filePath}`);
    return false;
  }
  return true;
}

const errors = [];

const sitemapPath = path.join(SRC_ROOT, 'sitemap.xml');
const robotsPath = path.join(SRC_ROOT, 'robots.txt');
const llmsPath = path.join(SRC_ROOT, 'llms.txt');
const llmsFullPath = path.join(SRC_ROOT, 'llms-full.txt');
const wkWebfinger = path.join(SRC_ROOT, '.well-known', 'webfinger');
const wkKeybase = path.join(SRC_ROOT, '.well-known', 'keybase.txt');
const wkSecurity = path.join(SRC_ROOT, '.well-known', 'security.txt');

if (ensureExists(sitemapPath, 'sitemap.xml')) {
  const xmlLint = spawnSync('xmllint', ['--noout', sitemapPath], { encoding: 'utf8' });
  if (xmlLint.error) {
    fail(`xmllint failed to execute: ${xmlLint.error.message}`);
  } else if (xmlLint.status !== 0) {
    const stderr = (xmlLint.stderr || '').trim();
    fail(`sitemap.xml is not valid XML: ${stderr || 'xmllint returned non-zero status'}`);
  }

  const sitemapContent = fs.readFileSync(sitemapPath, 'utf8');
  const locMatches = [...sitemapContent.matchAll(/<loc>([^<]+)<\/loc>/g)].map((m) => m[1].trim());
  for (const loc of locMatches) {
    try {
      const u = new URL(loc);
      if (u.protocol !== 'https:') {
        fail(`Sitemap URL is not HTTPS: ${loc}`);
      }
      if (u.hostname !== 'pngdeity.ru') {
        fail(`Sitemap URL is not canonical domain pngdeity.ru: ${loc}`);
      }
    } catch {
      fail(`Invalid sitemap URL: ${loc}`);
    }
  }
}

if (ensureExists(robotsPath, 'robots.txt')) {
  const robots = fs.readFileSync(robotsPath, 'utf8');
  for (const requiredLine of REQUIRED_ROBOTS_LINES) {
    if (!robots.includes(requiredLine)) {
      fail(`robots.txt missing required line: ${requiredLine}`);
    }
  }
}

ensureExists(wkWebfinger, '.well-known/webfinger');
ensureExists(wkKeybase, '.well-known/keybase.txt');
ensureExists(wkSecurity, '.well-known/security.txt');

if (ensureExists(llmsPath, 'llms.txt')) {
  if (!fs.readFileSync(llmsPath, 'utf8').trim()) {
    fail('llms.txt is empty');
  }
}
if (ensureExists(llmsFullPath, 'llms-full.txt')) {
  if (!fs.readFileSync(llmsFullPath, 'utf8').trim()) {
    fail('llms-full.txt is empty');
  }
}

if (errors.length > 0) {
  console.error('Metadata validation failed:');
  for (const error of errors) {
    console.error(`- ${error}`);
  }
  process.exit(1);
}

console.log('Metadata validation passed.');
