// Command snippets pulls Go code out of the companion code repository by
// anchor comment, so the deck never contains pasted Go and cannot drift
// from the program that runs on stage.
//
// It reads the directory named by CRS_CODE_DIR, which is how the deck
// builds against a sibling checkout, or vendor/crs-code, where CI clones
// it. It scans slides.md for
// `<<< @/snippets/<name>.go` imports, finds the matching
// `// snippet:start <name>` region, strips the common indentation, and
// writes snippets/<name>.go. Everything that resolves is written; anything
// missing is reported and the exit code is non-zero, so CI fails on a
// stale anchor while a local preview still renders the other slides.
package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	vendorDir   = "vendor/crs-code"
	snippetsDir = "snippets"
	slidesFile  = "slides.md"
)

var (
	startRe     = regexp.MustCompile(`//\s*snippet:start\s+(\S+)`)
	endRe       = regexp.MustCompile(`//\s*snippet:end\s+(\S+)`)
	referenceRe = regexp.MustCompile(`<<<\s*@/snippets/([\w-]+)\.go`)
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "snippets:", err)
		os.Exit(1)
	}
}

func codeDir() string {
	if dir := os.Getenv("CRS_CODE_DIR"); dir != "" {
		return dir
	}
	return vendorDir
}

func run() error {
	referenced, err := referencedAnchors(slidesFile)
	if err != nil {
		return fmt.Errorf("scan %s: %w", slidesFile, err)
	}

	dir := codeDir()
	found, err := extractAnchors(dir)
	if err != nil {
		return fmt.Errorf("extract anchors from %s: %w", dir, err)
	}

	if err := os.MkdirAll(snippetsDir, 0o755); err != nil {
		return err
	}

	var missing []string
	names := make([]string, 0, len(referenced))
	for name := range referenced {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		path := filepath.Join(snippetsDir, name+".go")
		body, ok := found[name]
		if !ok {
			missing = append(missing, name)
			// Slidev's snippet importer throws on a missing file, which
			// kills the whole render. A placeholder keeps the rest of the
			// deck previewable while this command still fails.
			body = fmt.Sprintf("// snippet %q not found in %s\n// npm run snippets is failing on this\n", name, dir)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			return err
		}
		fmt.Printf("snippets: wrote %s\n", path)
	}

	if len(missing) > 0 {
		return fmt.Errorf("slides.md references snippets with no anchor in %s: %s", dir, strings.Join(missing, ", "))
	}
	return nil
}

func referencedAnchors(path string) (map[string]struct{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	names := make(map[string]struct{})
	for _, m := range referenceRe.FindAllStringSubmatch(string(data), -1) {
		names[m[1]] = struct{}{}
	}
	return names, nil
}

func extractAnchors(root string) (map[string]string, error) {
	results := make(map[string]string)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		return extractFromFile(path, results)
	})
	if err != nil {
		return nil, err
	}
	return results, nil
}

func extractFromFile(path string, results map[string]string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	scanner := bufio.NewScanner(f)
	var open string
	var lines []string

	for scanner.Scan() {
		line := scanner.Text()
		if m := startRe.FindStringSubmatch(line); m != nil {
			if open != "" {
				return fmt.Errorf("%s: snippet %q starts before %q ends", path, m[1], open)
			}
			open, lines = m[1], nil
			continue
		}
		if m := endRe.FindStringSubmatch(line); m != nil {
			if open != m[1] {
				return fmt.Errorf("%s: snippet:end %q does not match open snippet %q", path, m[1], open)
			}
			if _, dup := results[open]; dup {
				return fmt.Errorf("snippet %q defined more than once (also in %s)", open, path)
			}
			results[open] = dedent(lines)
			open = ""
			continue
		}
		if open != "" {
			lines = append(lines, line)
		}
	}
	if open != "" {
		return fmt.Errorf("%s: snippet %q never closed", path, open)
	}
	return scanner.Err()
}

// dedent strips the indentation shared by every non-blank line, so a
// snippet from inside a function reads as top-level code on the slide.
func dedent(lines []string) string {
	common := -1
	for _, l := range lines {
		if strings.TrimSpace(l) == "" {
			continue
		}
		n := len(l) - len(strings.TrimLeft(l, "\t "))
		if common < 0 || n < common {
			common = n
		}
	}
	if common < 0 {
		common = 0
	}
	out := make([]string, len(lines))
	for i, l := range lines {
		if len(l) >= common {
			out[i] = l[common:]
		} else {
			out[i] = strings.TrimLeft(l, "\t ")
		}
	}
	return strings.Join(out, "\n") + "\n"
}
