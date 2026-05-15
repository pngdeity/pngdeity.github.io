package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"golang.org/x/net/html"
)

var canonicalOrigin = "https://pngdeity.ru"

type URL struct {
	XMLName xml.Name `xml:"url"`
	Loc     string   `xml:"loc"`
	LastMod string   `xml:"lastmod"`
}

type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	XMLNS   string   `xml:"xmlns,attr"`
	URLs    []URL
}

func isSitemapHTML(relPath string) bool {
	if !strings.HasSuffix(relPath, ".html") {
		return false
	}
	if relPath == "404.html" {
		return false
	}
	parts := strings.SplitN(relPath, "/", 2)
	top := parts[0]
	return !strings.Contains(relPath, "/") || top == "blog" || top == "app"
}

func filePathToURL(relPath string) string {
	if relPath == "index.html" {
		return canonicalOrigin + "/"
	}
	if strings.HasSuffix(relPath, "/index.html") {
		return canonicalOrigin + "/" + strings.TrimSuffix(relPath, "index.html")
	}
	return canonicalOrigin + "/" + relPath
}

func getGitLastMod(repoRoot, sourcePath string) string {
	cmd := exec.Command("git", "log", "-1", "--format=%cI", "--", sourcePath)
	cmd.Dir = repoRoot
	out, err := cmd.Output()
	if err == nil {
		str := strings.TrimSpace(string(out))
		if str != "" {
			return str
		}
	}
	return ""
}

func getLastMod(repoRoot, srcRoot, filePath string) string {
	relPath, err := filepath.Rel(srcRoot, filePath)
	if err != nil {
		relPath = filePath
	}
	relPathPosix := filepath.ToSlash(relPath)
	sourcePath := filePath

	if strings.HasPrefix(relPathPosix, "blog/") && strings.HasSuffix(relPathPosix, ".html") {
		base := relPathPosix[len("blog/") : len(relPathPosix)-len(".html")]
		if strings.HasSuffix(base, "/index") {
			baseDir := strings.TrimSuffix(base, "index")
			cand1 := filepath.Join(repoRoot, "hugo-src", "content", baseDir, "_index.md")
			cand2 := filepath.Join(repoRoot, "hugo-src", "content", baseDir, "index.md")
			if _, err := os.Stat(cand1); err == nil {
				sourcePath = cand1
			} else if _, err := os.Stat(cand2); err == nil {
				sourcePath = cand2
			}
		} else {
			cand := filepath.Join(repoRoot, "hugo-src", "content", base+".md")
			if _, err := os.Stat(cand); err == nil {
				sourcePath = cand
			}
		}
	}

	gitDate := getGitLastMod(repoRoot, sourcePath)
	if gitDate != "" {
		return gitDate
	}

	stat, err := os.Stat(filePath)
	if err != nil {
		return time.Now().UTC().Format(time.RFC3339)
	}
	return stat.ModTime().UTC().Format(time.RFC3339)
}

func stripTags(s string) string {
	re := regexp.MustCompile(`\s+`)
	return strings.TrimSpace(re.ReplaceAllString(s, " "))
}

type Metadata struct {
	Title   string
	Heading string
}

func extractMetadata(filePath string) Metadata {
	file, err := os.Open(filePath)
	if err != nil {
		return Metadata{}
	}
	defer file.Close()

	z := html.NewTokenizer(file)
	var title, heading string
	inTitle, inH1 := false, false

	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			break
		}
		switch tt {
		case html.StartTagToken:
			t := z.Token()
			if t.Data == "title" && title == "" {
				inTitle = true
			} else if t.Data == "h1" && heading == "" {
				inH1 = true
			}
		case html.EndTagToken:
			t := z.Token()
			if t.Data == "title" {
				inTitle = false
			} else if t.Data == "h1" {
				inH1 = false
			}
		case html.TextToken:
			if inTitle {
				title += z.Token().Data
			} else if inH1 {
				heading += z.Token().Data
			}
		}
	}

	return Metadata{
		Title:   stripTags(title),
		Heading: stripTags(heading),
	}
}

func fallbackName(relPath string) string {
	base := strings.TrimSuffix(filepath.Base(relPath), filepath.Ext(relPath))
	re := regexp.MustCompile(`[-_]+`)
	base = re.ReplaceAllString(base, " ")
	if base == "" {
		return relPath
	}
	// Note: Title is deprecated, Title cases strings instead of unicode proper.
	// To fix later, but keeps same behavior.
	return strings.Title(base)
}

func writeFile(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func main() {
	var repoRoot string
	flag.StringVar(&repoRoot, "repo", ".", "Repository root directory")
	flag.Parse()

	absRepo, err := filepath.Abs(repoRoot)
	if err != nil {
		fmt.Println("Invalid repo path")
		os.Exit(1)
	}
	srcRoot := filepath.Join(absRepo, "src")

	if _, err := os.Stat(srcRoot); os.IsNotExist(err) {
		fmt.Printf("Missing source directory: %s\n", srcRoot)
		os.Exit(1)
	}

	var htmlFiles []string
	filepath.Walk(srcRoot, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			rel, _ := filepath.Rel(srcRoot, path)
			if isSitemapHTML(filepath.ToSlash(rel)) {
				htmlFiles = append(htmlFiles, path)
			}
		}
		return nil
	})

	sort.Strings(htmlFiles)

	robotsPath := filepath.Join(srcRoot, "robots.txt")
	robotsContent := fmt.Sprintf("User-agent: *\nAllow: /\nSitemap: %s/sitemap.xml\n\n", canonicalOrigin)
	writeFile(robotsPath, robotsContent)

	var urls []URL
	type Entry struct {
		Loc      string
		LastMod  string
		RelPath  string
		FullPath string
	}
	var entries []Entry

	for _, f := range htmlFiles {
		rel, _ := filepath.Rel(srcRoot, f)
		relPosix := filepath.ToSlash(rel)
		loc := filePathToURL(relPosix)
		lastMod := getLastMod(absRepo, srcRoot, f)
		urls = append(urls, URL{Loc: loc, LastMod: lastMod})
		entries = append(entries, Entry{Loc: loc, LastMod: lastMod, RelPath: relPosix, FullPath: f})
	}

	urlset := URLSet{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  urls,
	}
	outXML, _ := xml.MarshalIndent(urlset, "", "  ")
	sitemapContent := xml.Header + string(outXML) + "\n"
	writeFile(filepath.Join(srcRoot, "sitemap.xml"), sitemapContent)

	var rootPages, blogPages, appPages []Entry
	for _, e := range entries {
		if !strings.Contains(e.RelPath, "/") {
			rootPages = append(rootPages, e)
		} else if strings.HasPrefix(e.RelPath, "blog/") {
			blogPages = append(blogPages, e)
		} else if strings.HasPrefix(e.RelPath, "app/") {
			appPages = append(appPages, e)
		}
	}

	discoveredSections := map[string]bool{}
	for _, p := range blogPages {
		rest := strings.TrimPrefix(p.RelPath, "blog/")
		parts := strings.SplitN(rest, "/", 2)
		section := parts[0]
		if section != "" && section != "index.html" {
			discoveredSections[section] = true
		}
	}

	var sectionsList []string
	for k := range discoveredSections {
		sectionsList = append(sectionsList, k)
	}
	sort.Strings(sectionsList)

	var keyPages []string
	for i, e := range entries {
		if i >= 12 {
			break
		}
		md := extractMetadata(e.FullPath)
		name := md.Title
		if name == "" {
			name = md.Heading
		}
		if name == "" {
			name = fallbackName(e.RelPath)
		}
		keyPages = append(keyPages, fmt.Sprintf("- %s (%s)", name, e.Loc))
	}

	siteOverview := fmt.Sprintf("Site generated from static HTML plus Hugo blog content and an optional Blazor app section under /app/. Canonical host: %s/.", canonicalOrigin)
	appSection := "- /app/ section is currently not populated"
	if len(appPages) > 0 {
		appSection = "- Application assets/content under /app/"
	}
	mainSections := fmt.Sprintf("- Root static pages (landing and standalone HTML pages)\n- Blog content under /blog/\n%s", appSection)
	contentTypes := "- HTML pages\n- Downloadable static assets (documents, images, and related files)\n- Generated application files (when /app/ artifact is present)"
	attribution := "- License file present in repository: LICENSE.md (see project root)."
	contact := "- Security contact: mailto:security@pngdeity.ru"

	llmsContent := fmt.Sprintf("# Site Overview\n%s\n\n# Main Sections\n%s\n\n# Content Types\n%s\n\n# Attribution / Licensing\n%s\n\n# Contact\n%s\n\n", siteOverview, mainSections, contentTypes, attribution, contact)
	writeFile(filepath.Join(srcRoot, "llms.txt"), llmsContent)

	var rootPagesList []string
	for i, p := range rootPages {
		if i >= 50 {
			break
		}
		rootPagesList = append(rootPagesList, fmt.Sprintf("- %s -> %s", p.RelPath, p.Loc))
	}
	rootPagesStr := "- None found"
	if len(rootPagesList) > 0 {
		rootPagesStr = strings.Join(rootPagesList, "\n")
	}

	sectionsStr := "- No nested Hugo sections discovered"
	if len(sectionsList) > 0 {
		for i, s := range sectionsList {
			sectionsList[i] = "- " + s
		}
		sectionsStr = strings.Join(sectionsList, "\n")
	}

	keyPagesStr := "- No HTML pages discovered for sampling"
	if len(keyPages) > 0 {
		keyPagesStr = strings.Join(keyPages, "\n")
	}

	appStr := "- /app/ not detected or contains no HTML pages in this build output."
	if len(appPages) > 0 {
		appStr = fmt.Sprintf("- /app/ detected with %d HTML page(s) in this build output.", len(appPages))
	}

	llmsFullContent := fmt.Sprintf("# Site Overview\n%s\n\n# Main Sections\n%s\n\n# Content Types\n%s\n\n# Hugo Major Sections\n%s\n\n# Key Static Pages (src/)\n%s\n\n# Representative Page Samples\n%s\n\n# Application Section\n%s\n\n# Attribution / Licensing\n%s\n\n# Contact\n%s\n\n", siteOverview, mainSections, contentTypes, sectionsStr, rootPagesStr, keyPagesStr, appStr, attribution, contact)
	writeFile(filepath.Join(srcRoot, "llms-full.txt"), llmsFullContent)

	webfinger := map[string]interface{}{
		"subject": "acct:pngdeity@pngdeity.ru",
		"aliases": []string{"https://pngdeity.ru/"},
		"links": []map[string]string{
			{"rel": "self", "type": "application/activity+json", "href": "https://pngdeity.ru/"},
		},
	}
	wb, _ := json.MarshalIndent(webfinger, "", "  ")
	writeFile(filepath.Join(srcRoot, ".well-known", "webfinger"), string(wb)+"\n")
	writeFile(filepath.Join(srcRoot, ".well-known", "keybase.txt"), "# Keybase proof placeholder\n")

	expires := time.Now().UTC().AddDate(1, 0, 0).Format("2006-01-02")
	securityTxt := fmt.Sprintf("Contact: mailto:security@pngdeity.ru\nExpires: %s\nPolicy: %s/\n\n", expires, canonicalOrigin)
	writeFile(filepath.Join(srcRoot, ".well-known", "security.txt"), securityTxt)

	fmt.Printf("Generated metadata files in %s\n", srcRoot)
	fmt.Printf("Sitemap entries: %d\n", len(urls))
}
