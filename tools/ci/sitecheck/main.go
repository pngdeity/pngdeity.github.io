// Command sitecheck performs advisory integrity checks over a merged site tree.
//
// It resolves every local href/src/poster reference in HTML files against the
// files on disk (root-absolute, page-relative, and same-origin absolute URLs),
// and asserts a couple of site invariants. Findings are written to the GitHub
// step summary and emitted as workflow annotations; the process still exits
// non-zero so a wrapping step can decide whether to tolerate the failure.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/net/html"
)

const canonicalOrigin = "https://pngdeity.ru"

type refKind int

const (
	refSkipped refKind = iota
	refExternal
	refInternal
)

type resolution struct {
	Kind       refKind
	Candidates []string
	Relative   bool
}

type refLoc struct {
	Ref  string
	Line int
}

type finding struct {
	File    string
	Line    int
	Message string
}

func candidatesForPath(srcRoot, servedPath string) []string {
	p := servedPath
	if i := strings.IndexAny(p, "?#"); i >= 0 {
		p = p[:i]
	}
	p = strings.ReplaceAll(p, "\\", "/")
	trailingSlash := strings.HasSuffix(p, "/")

	cleaned := path.Clean("/" + strings.TrimPrefix(p, "/"))
	if cleaned == "/" || cleaned == "." {
		return []string{filepath.Join(srcRoot, "index.html")}
	}
	rel := strings.TrimPrefix(cleaned, "/")
	base := filepath.Join(srcRoot, filepath.FromSlash(rel))

	if trailingSlash {
		return []string{filepath.Join(base, "index.html")}
	}
	return []string{base, base + ".html", filepath.Join(base, "index.html")}
}

// resolve maps a raw reference from a page to candidate filesystem paths.
// originHost is the canonical host without a port (for example "pngdeity.ru").
func resolve(srcRoot, pageRelPosix, rawRef, originHost string) resolution {
	ref := strings.TrimSpace(rawRef)
	if ref == "" || strings.HasPrefix(ref, "#") {
		return resolution{Kind: refSkipped}
	}
	lower := strings.ToLower(ref)
	for _, prefix := range []string{"mailto:", "tel:", "data:", "javascript:", "blob:"} {
		if strings.HasPrefix(lower, prefix) {
			return resolution{Kind: refSkipped}
		}
	}
	if strings.HasPrefix(ref, "//") {
		ref = "https:" + ref
		lower = strings.ToLower(ref)
	}
	if strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://") {
		u, err := url.Parse(ref)
		if err != nil || u.Hostname() == "" {
			return resolution{Kind: refSkipped}
		}
		if !strings.EqualFold(u.Hostname(), originHost) {
			return resolution{Kind: refExternal}
		}
		return resolution{Kind: refInternal, Candidates: candidatesForPath(srcRoot, u.Path)}
	}

	p := ref
	if i := strings.IndexAny(p, "?#"); i >= 0 {
		p = p[:i]
	}
	if p == "" {
		return resolution{Kind: refSkipped}
	}
	if strings.HasPrefix(p, "/") {
		return resolution{Kind: refInternal, Candidates: candidatesForPath(srcRoot, p)}
	}

	joined := path.Join(path.Dir(pageRelPosix), p)
	return resolution{
		Kind:       refInternal,
		Relative:   true,
		Candidates: candidatesForPath(srcRoot, "/"+joined),
	}
}

func existsAny(candidates []string) bool {
	for _, c := range candidates {
		if fi, err := os.Stat(c); err == nil && !fi.IsDir() {
			return true
		}
	}
	return false
}

func lineAt(data []byte, offset int64) int {
	if offset < 0 {
		offset = 0
	}
	if offset > int64(len(data)) {
		offset = int64(len(data))
	}
	return 1 + bytes.Count(data[:offset], []byte("\n"))
}

func extractRefs(htmlPath string) ([]refLoc, error) {
	data, err := os.ReadFile(htmlPath)
	if err != nil {
		return nil, err
	}
	z := html.NewTokenizer(bytes.NewReader(data))
	var refs []refLoc
	cursor := 0
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			if err := z.Err(); err != nil && err != io.EOF {
				return refs, err
			}
			break
		}
		if tt != html.StartTagToken && tt != html.SelfClosingTagToken {
			continue
		}
		raw := z.Raw()
		line := 1
		if idx := bytes.Index(data[cursor:], raw); idx >= 0 {
			pos := cursor + idx
			line = lineAt(data, int64(pos))
			cursor = pos + len(raw)
		}
		token := z.Token()
		for _, attr := range token.Attr {
			if attr.Namespace != "" {
				continue
			}
			switch strings.ToLower(attr.Key) {
			case "href", "src", "poster":
				refs = append(refs, refLoc{Ref: attr.Val, Line: line})
			}
		}
	}
	return refs, nil
}

func checkCNAME(repoRoot, originHost string) (string, bool) {
	data, err := os.ReadFile(filepath.Join(repoRoot, "CNAME"))
	if err != nil {
		return "", false
	}
	return strings.TrimSpace(string(data)), true
}

func shortestCandidate(candidates []string) string {
	if len(candidates) == 0 {
		return "(no candidates)"
	}
	best := candidates[0]
	for _, c := range candidates[1:] {
		if len(c) < len(best) {
			best = c
		}
	}
	return best
}

func checkPage(srcRoot, pageRel string, refs []refLoc, originHost string) ([]finding, int) {
	var findings []finding
	internal := 0
	for _, r := range refs {
		res := resolve(srcRoot, pageRel, r.Ref, originHost)
		if res.Kind != refInternal {
			continue
		}
		internal++
		if pageRel == "404.html" && res.Relative {
			findings = append(findings, finding{
				File:    pageRel,
				Line:    r.Line,
				Message: fmt.Sprintf("relative reference %q breaks when this page is served from a nested path; use a root-absolute path", r.Ref),
			})
			continue
		}
		if !existsAny(res.Candidates) {
			findings = append(findings, finding{
				File:    pageRel,
				Line:    r.Line,
				Message: fmt.Sprintf("unresolved local reference %q (looked for %s)", r.Ref, shortestCandidate(res.Candidates)),
			})
		}
	}
	return findings, internal
}

// sanitize collapses control characters so values cannot break workflow
// annotations or the step summary when interpolated.
func sanitize(s string) string {
	return strings.NewReplacer("\r", " ", "\n", " ").Replace(s)
}

// runChecks walks srcRoot and returns findings plus counts of scanned HTML
// files and resolved local references.
func runChecks(srcRoot, repoRoot, originHost string) ([]finding, int, int, error) {
	var htmlFiles []string
	err := filepath.Walk(srcRoot, func(p string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !info.IsDir() && strings.EqualFold(filepath.Ext(p), ".html") {
			htmlFiles = append(htmlFiles, p)
		}
		return nil
	})
	if err != nil {
		return nil, 0, 0, err
	}
	sort.Strings(htmlFiles)

	var findings []finding
	refCount := 0
	for _, htmlPath := range htmlFiles {
		rel, err := filepath.Rel(srcRoot, htmlPath)
		if err != nil {
			rel = htmlPath
		}
		pageRel := filepath.ToSlash(rel)

		refs, err := extractRefs(htmlPath)
		if err != nil {
			findings = append(findings, finding{File: pageRel, Message: fmt.Sprintf("failed to parse HTML: %v", err)})
			continue
		}
		pageFindings, internal := checkPage(srcRoot, pageRel, refs, originHost)
		refCount += internal
		findings = append(findings, pageFindings...)
	}

	if cname, ok := checkCNAME(repoRoot, originHost); ok && cname != originHost {
		findings = append(findings, finding{
			File:    "CNAME",
			Line:    1,
			Message: fmt.Sprintf("CNAME is %q, expected %q", cname, originHost),
		})
	}

	sort.SliceStable(findings, func(i, j int) bool {
		if findings[i].File != findings[j].File {
			return findings[i].File < findings[j].File
		}
		return findings[i].Line < findings[j].Line
	})

	return findings, len(htmlFiles), refCount, nil
}

func writeSummary(findings []finding, htmlCount, refCount int) {
	summaryFile := os.Getenv("GITHUB_STEP_SUMMARY")
	if summaryFile == "" {
		return
	}
	f, err := os.OpenFile(summaryFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sitecheck: writing summary: %v\n", err)
		return
	}
	defer f.Close()

	fmt.Fprintln(f, "## Site Check (advisory)")
	fmt.Fprintf(f, "Scanned %d HTML files and resolved %d local references.\n\n", htmlCount, refCount)
	if len(findings) == 0 {
		fmt.Fprintln(f, "- No issues found.")
		return
	}
	fmt.Fprintf(f, "- %d issue(s) found:\n", len(findings))
	for _, fd := range findings {
		fmt.Fprintf(f, "  - `%s:%d` %s\n", sanitize(fd.File), fd.Line, sanitize(fd.Message))
	}
}

func main() {
	repo := flag.String("repo", ".", "repository root")
	srcFlag := flag.String("src", "", "merged site root (defaults to <repo>/src)")
	origin := flag.String("origin", canonicalOrigin, "canonical site origin")
	flag.Parse()

	repoRoot, err := filepath.Abs(*repo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sitecheck: invalid repo path: %v\n", err)
		os.Exit(1)
	}
	srcRoot := *srcFlag
	if srcRoot == "" {
		srcRoot = filepath.Join(repoRoot, "src")
	}
	srcRoot, err = filepath.Abs(srcRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sitecheck: invalid src path: %v\n", err)
		os.Exit(1)
	}
	if _, err := os.Stat(srcRoot); err != nil {
		fmt.Fprintf(os.Stderr, "sitecheck: missing site root: %s\n", srcRoot)
		os.Exit(1)
	}

	originURL, err := url.Parse(*origin)
	if err != nil || originURL.Hostname() == "" {
		fmt.Fprintf(os.Stderr, "sitecheck: invalid origin: %s\n", *origin)
		os.Exit(1)
	}

	findings, htmlCount, refCount, err := runChecks(srcRoot, repoRoot, originURL.Hostname())
	if err != nil {
		fmt.Fprintf(os.Stderr, "sitecheck: %v\n", err)
		os.Exit(1)
	}

	for _, fd := range findings {
		file := sanitize(fd.File)
		message := sanitize(fd.Message)
		fmt.Printf("::warning file=%s,line=%d::%s\n", file, fd.Line, message)
		fmt.Printf("WARN %s:%d %s\n", file, fd.Line, message)
	}

	writeSummary(findings, htmlCount, refCount)

	if len(findings) > 0 {
		fmt.Fprintf(os.Stderr, "sitecheck: %d issue(s) found\n", len(findings))
		os.Exit(1)
	}
	fmt.Printf("sitecheck: OK (%d HTML files, %d local references)\n", htmlCount, refCount)
}
