package main

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	URLs    []URL    `xml:"url"`
}

type URL struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: validatemetadata <src_dir>")
		os.Exit(1)
	}
	srcDir := os.Args[1]

	errors := []string{}
	fail := func(msg string) {
		errors = append(errors, msg)
	}

	ensureExists := func(filename string) bool {
		path := filepath.Join(srcDir, filename)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fail(fmt.Sprintf("Missing required file: %s", filename))
			return false
		}
		return true
	}

	sitemapPath := filepath.Join(srcDir, "sitemap.xml")
	if ensureExists("sitemap.xml") {
		data, err := os.ReadFile(sitemapPath)
		if err != nil {
			fail(fmt.Sprintf("Failed to read sitemap.xml: %v", err))
		} else {
			var urlset URLSet
			if err := xml.Unmarshal(data, &urlset); err != nil {
				fail(fmt.Sprintf("sitemap.xml is not valid XML: %v", err))
			} else {
				for _, u := range urlset.URLs {
					parsed, err := url.Parse(u.Loc)
					if err != nil || parsed.Scheme == "" || parsed.Host == "" {
						fail(fmt.Sprintf("Invalid sitemap URL: %s", u.Loc))
						continue
					}
					if parsed.Scheme != "https" {
						fail(fmt.Sprintf("Sitemap URL is not HTTPS: %s", u.Loc))
					}
					if parsed.Hostname() != "pngdeity.ru" {
						fail(fmt.Sprintf("Sitemap URL is not canonical domain pngdeity.ru: %s", u.Loc))
					}
				}
			}
		}
	}

	robotsPath := filepath.Join(srcDir, "robots.txt")
	if ensureExists("robots.txt") {
		data, _ := os.ReadFile(robotsPath)
		content := string(data)
		requiredLines := []string{
			"User-agent: *",
			"Allow: /",
			"Sitemap: https://pngdeity.ru/sitemap.xml",
		}
		for _, line := range requiredLines {
			if !strings.Contains(content, line) {
				fail(fmt.Sprintf("robots.txt missing required line: %s", line))
			}
		}
	}

	ensureExists(".well-known/webfinger")
	ensureExists(".well-known/keybase.txt")
	ensureExists(".well-known/security.txt")

	checkNotEmpty := func(filename string) {
		path := filepath.Join(srcDir, filename)
		if _, err := os.Stat(path); err == nil {
			data, _ := os.ReadFile(path)
			if strings.TrimSpace(string(data)) == "" {
				fail(fmt.Sprintf("%s is empty", filename))
			}
		}
	}

	checkNotEmpty("llms.txt")
	checkNotEmpty("llms-full.txt")

	if len(errors) > 0 {
		fmt.Fprintln(os.Stderr, "Metadata validation failed:")
		for _, err := range errors {
			fmt.Fprintf(os.Stderr, "- %s\n", err)
		}
		os.Exit(1)
	}

	fmt.Println("Metadata validation passed.")
}
