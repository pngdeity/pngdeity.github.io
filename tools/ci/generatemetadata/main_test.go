package main

import "testing"

func TestIsSitemapHTML(t *testing.T) {
	cases := []struct {
		rel  string
		want bool
	}{
		{"index.html", true},
		{"archive.html", true},
		{"colophon.html", true},
		{"404.html", false},
		{"blog/404.html", false},
		{"blog/index.html", true},
		{"blog/posts/puppet/index.html", true},
		{"misc/index.html", false},
		{"downloads/notes.html", false},
		{"style.css", false},
		{"blog/index.xml", false},
		{"", false},
	}
	for _, tc := range cases {
		if got := isSitemapHTML(tc.rel); got != tc.want {
			t.Errorf("isSitemapHTML(%q) = %v, want %v", tc.rel, got, tc.want)
		}
	}
}

func TestFilePathToURL(t *testing.T) {
	cases := []struct {
		rel  string
		want string
	}{
		{"index.html", canonicalOrigin + "/"},
		{"archive.html", canonicalOrigin + "/archive.html"},
		{"blog/index.html", canonicalOrigin + "/blog/"},
		{"blog/posts/puppet/index.html", canonicalOrigin + "/blog/posts/puppet/"},
		{"404.html", canonicalOrigin + "/404.html"},
	}
	for _, tc := range cases {
		if got := filePathToURL(tc.rel); got != tc.want {
			t.Errorf("filePathToURL(%q) = %q, want %q", tc.rel, got, tc.want)
		}
	}
}

func TestFallbackName(t *testing.T) {
	cases := []struct {
		rel  string
		want string
	}{
		{"archive-page.html", "Archive Page"},
		{"my_page.md", "My Page"},
	}
	for _, tc := range cases {
		if got := fallbackName(tc.rel); got != tc.want {
			t.Errorf("fallbackName(%q) = %q, want %q", tc.rel, got, tc.want)
		}
	}
}
