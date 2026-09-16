package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCandidatesForPath(t *testing.T) {
	src := "/site"
	cases := []struct {
		name   string
		served string
		want   []string
	}{
		{
			name:   "file",
			served: "/downloads/a.pdf",
			want: []string{
				filepath.Join(src, "downloads", "a.pdf"),
				filepath.Join(src, "downloads", "a.pdf.html"),
				filepath.Join(src, "downloads", "a.pdf", "index.html"),
			},
		},
		{
			name:   "directory",
			served: "/blog/",
			want:   []string{filepath.Join(src, "blog", "index.html")},
		},
		{
			name:   "root slash",
			served: "/",
			want:   []string{filepath.Join(src, "index.html")},
		},
		{
			name:   "empty",
			served: "",
			want:   []string{filepath.Join(src, "index.html")},
		},
		{
			name:   "query stripped",
			served: "/style.css?v=2",
			want: []string{
				filepath.Join(src, "style.css"),
				filepath.Join(src, "style.css.html"),
				filepath.Join(src, "style.css", "index.html"),
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := candidatesForPath(src, tc.served)
			if len(got) != len(tc.want) {
				t.Fatalf("candidatesForPath(%q) = %v, want %v", tc.served, got, tc.want)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Errorf("candidate[%d] = %q, want %q", i, got[i], tc.want[i])
				}
			}
		})
	}
}

func TestResolve(t *testing.T) {
	const src = "/site"
	const host = "pngdeity.ru"
	cases := []struct {
		name     string
		page     string
		ref      string
		kind     refKind
		relative bool
		want     string
	}{
		{name: "external", page: "index.html", ref: "https://example.com/a", kind: refExternal},
		{name: "same origin absolute", page: "blog/posts/p/index.html", ref: "https://pngdeity.ru/downloads/a.pdf", kind: refInternal, want: filepath.Join(src, "downloads", "a.pdf")},
		{name: "same origin root", page: "index.html", ref: "https://pngdeity.ru/", kind: refInternal, want: filepath.Join(src, "index.html")},
		{name: "protocol relative same host", page: "index.html", ref: "//pngdeity.ru/blog/", kind: refInternal, want: filepath.Join(src, "blog", "index.html")},
		{name: "protocol relative external", page: "index.html", ref: "//example.com/x", kind: refExternal},
		{name: "root absolute", page: "index.html", ref: "/style.css", kind: refInternal, want: filepath.Join(src, "style.css")},
		{name: "relative root page", page: "404.html", ref: "style.css", kind: refInternal, relative: true, want: filepath.Join(src, "style.css")},
		{name: "relative nested", page: "blog/posts/p/index.html", ref: "../img/a.png", kind: refInternal, relative: true, want: filepath.Join(src, "blog", "posts", "img", "a.png")},
		{name: "fragment", page: "index.html", ref: "#top", kind: refSkipped},
		{name: "mailto", page: "index.html", ref: "mailto:a@b.c", kind: refSkipped},
		{name: "empty", page: "index.html", ref: "   ", kind: refSkipped},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := resolve(src, tc.page, tc.ref, host)
			if got.Kind != tc.kind {
				t.Fatalf("resolve(%q).Kind = %v, want %v", tc.ref, got.Kind, tc.kind)
			}
			if got.Relative != tc.relative {
				t.Errorf("resolve(%q).Relative = %v, want %v", tc.ref, got.Relative, tc.relative)
			}
			if tc.want != "" {
				if !contains(got.Candidates, tc.want) {
					t.Errorf("resolve(%q).Candidates = %v, want to contain %q", tc.ref, got.Candidates, tc.want)
				}
			}
		})
	}
}

func TestSanitize(t *testing.T) {
	got := sanitize("a\r\nb\tc")
	if got != "a  b\tc" {
		t.Errorf("sanitize() = %q, want %q", got, "a  b\tc")
	}
}

func contains(values []string, want string) bool {
	for _, v := range values {
		if v == want {
			return true
		}
	}
	return false
}

func writeFile(t *testing.T, dir, rel, content string) {
	t.Helper()
	p := filepath.Join(dir, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("write %s: %v", rel, err)
	}
}

func TestRunChecks(t *testing.T) {
	repo := t.TempDir()
	src := filepath.Join(repo, "src")

	writeFile(t, src, "index.html", `<html><head><link rel="stylesheet" href="/style.css"></head><body><a href="blog/">blog</a></body></html>`)
	writeFile(t, src, "style.css", "body{}")
	writeFile(t, src, "blog/index.html", `<a href="https://pngdeity.ru/downloads/a.pdf">pdf</a>`)
	writeFile(t, src, "downloads/a.pdf", "pdf")
	writeFile(t, src, "404.html", `<a href="style.css">relative</a>`)
	writeFile(t, src, "missing.html", `<a href="/downloads/absent.pdf">gone</a>`)
	writeFile(t, repo, "CNAME", "pngdeity.ru\n")

	findings, htmlCount, refCount, err := runChecks(src, repo, "pngdeity.ru")
	if err != nil {
		t.Fatalf("runChecks: %v", err)
	}
	if htmlCount != 4 {
		t.Errorf("htmlCount = %d, want 4", htmlCount)
	}
	if refCount == 0 {
		t.Errorf("refCount = 0, want > 0")
	}

	var messages []string
	for _, f := range findings {
		messages = append(messages, f.File+": "+f.Message)
	}
	joined := strings.Join(messages, "\n")

	if !strings.Contains(joined, "404.html") || !strings.Contains(joined, "relative reference") {
		t.Errorf("expected a 404 relative-reference finding, got:\n%s", joined)
	}
	if !strings.Contains(joined, "absent.pdf") {
		t.Errorf("expected an unresolved reference finding, got:\n%s", joined)
	}
	if len(findings) != 2 {
		t.Errorf("findings = %d, want 2:\n%s", len(findings), joined)
	}

	// A wrong CNAME is reported.
	writeFile(t, repo, "CNAME", "example.com\n")
	findings, _, _, err = runChecks(src, repo, "pngdeity.ru")
	if err != nil {
		t.Fatalf("runChecks: %v", err)
	}
	foundCNAME := false
	for _, f := range findings {
		if f.File == "CNAME" {
			foundCNAME = true
		}
	}
	if !foundCNAME {
		t.Errorf("expected a CNAME finding for a mismatched CNAME")
	}
}
