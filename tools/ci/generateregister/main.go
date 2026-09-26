// Command generateregister renders hugo-src/data/registrations.toml into the
// register region of src/links.html.
//
// The static half has no server-side includes, so the wall of registrations
// cannot be templated at request time the way the blog is. Instead this tool
// reads the canonical data file and writes the rendered markup between the
// BEGIN/END markers in src/links.html, leaving everything outside them
// untouched. The data file stays the source of truth; the page is a view of it.
//
// Usage:
//
//	go run ./generateregister -repo ../..
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	beginMarker = "<!-- BEGIN REGISTER -->"
	endMarker   = "<!-- END REGISTER -->"
)

// Registration is one institution the handle is registered at.
//
// URL is where the cell points: the profile page where one is known and
// verified, otherwise the service root.
//
// Accounts is the number of separate accounts held there -- the piece's one
// datum. Since is the earliest creation date, retained for provenance rather
// than rendered on the wall.
type Registration struct {
	Name     string
	URL      string
	Accounts int
	Since    string
}

func parseRegistrations(path string) ([]Registration, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var out []Registration
	var cur Registration
	inEntry := false

	flush := func() {
		if inEntry && cur.Name != "" {
			out = append(out, cur)
		}
		cur = Registration{}
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if line == "[[registrations]]" {
			flush()
			inEntry = true
			continue
		}
		if !inEntry {
			continue
		}
		key, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		switch key {
		case "name":
			cur.Name = strings.Trim(value, `"`)
		case "url":
			cur.URL = strings.Trim(value, `"`)
		case "since":
			cur.Since = strings.Trim(value, `"`)
		case "accounts":
			n, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("accounts %q: %w", value, err)
			}
			cur.Accounts = n
		}
	}
	flush()
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// renderWall produces the register markup. The grid is uniform: every cell
// carries the handle and the institution, and the only variation is the
// account count, which is the encoded dimension. Cells with more than one
// account are marked so the long tail reads at a glance.
//
// The whole cell is the link target, so the wall stays free of inline link
// noise and the click area matches the visible cell.
func renderWall(regs []Registration) string {
	sort.Slice(regs, func(i, j int) bool {
		if regs[i].Accounts != regs[j].Accounts {
			return regs[i].Accounts > regs[j].Accounts
		}
		return strings.ToLower(regs[i].Name) < strings.ToLower(regs[j].Name)
	})

	var b strings.Builder
	b.WriteString(beginMarker)
	b.WriteString("\n\t\t<ul class=\"register\">\n")
	for _, r := range regs {
		class := "register__cell"
		if r.Accounts > 1 {
			class += " register__cell--multiple"
		}
		count := ""
		if r.Accounts > 1 {
			count = fmt.Sprintf("\n\t\t\t\t\t<span class=\"register__count\">%d</span>", r.Accounts)
		}
		fmt.Fprintf(&b,
			"\t\t\t<li class=\"%s\">\n\t\t\t\t<a class=\"register__link\" href=\"%s\">\n\t\t\t\t\t<span class=\"register__handle\">pngdeity</span>\n\t\t\t\t\t<span class=\"register__institution\">%s</span>%s\n\t\t\t\t</a>\n\t\t\t</li>\n",
			class, escapeAttr(r.URL), escape(r.Name), count)
	}
	b.WriteString("\t\t</ul>\n\t\t")
	b.WriteString(endMarker)
	return b.String()
}

func escape(s string) string {
	r := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;")
	return r.Replace(s)
}

func escapeAttr(s string) string {
	r := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;", "'", "&#39;")
	return r.Replace(s)
}

// replaceRegion swaps the text between the markers, preserving everything
// outside them so hand-maintained markup around the wall is never clobbered.
func replaceRegion(page, wall string) (string, error) {
	start := strings.Index(page, beginMarker)
	end := strings.Index(page, endMarker)
	if start < 0 || end < 0 || end < start {
		return "", fmt.Errorf("markers %s / %s not found in page", beginMarker, endMarker)
	}
	return page[:start] + wall + page[end+len(endMarker):], nil
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

	dataPath := filepath.Join(absRepo, "hugo-src", "data", "registrations.toml")
	pagePath := filepath.Join(absRepo, "src", "links.html")

	regs, err := parseRegistrations(dataPath)
	if err != nil {
		fmt.Printf("Reading %s: %v\n", dataPath, err)
		os.Exit(1)
	}
	if len(regs) == 0 {
		fmt.Printf("No registrations found in %s\n", dataPath)
		os.Exit(1)
	}

	page, err := os.ReadFile(pagePath)
	if err != nil {
		fmt.Printf("Reading %s: %v\n", pagePath, err)
		os.Exit(1)
	}

	updated, err := replaceRegion(string(page), renderWall(regs))
	if err != nil {
		fmt.Printf("Rendering %s: %v\n", pagePath, err)
		os.Exit(1)
	}

	if err := os.WriteFile(pagePath, []byte(updated), 0644); err != nil {
		fmt.Printf("Writing %s: %v\n", pagePath, err)
		os.Exit(1)
	}

	total := 0
	for _, r := range regs {
		total += r.Accounts
	}
	fmt.Printf("register: %d institutions, %d accounts\n", len(regs), total)
}
