package skillmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"
)

// MaxDescriptionLen is the Agent Skills specification hard limit for description.
const MaxDescriptionLen = 1024

var frontmatterRe = regexp.MustCompile(`(?s)^---\n(.*?)\n---\n?`)

// ClampDescription trims and truncates to MaxDescriptionLen on a rune boundary.
func ClampDescription(s string) string {
	s = strings.TrimSpace(s)
	if utf8.RuneCountInString(s) <= MaxDescriptionLen {
		return s
	}
	runes := []rune(s)
	return string(runes[:MaxDescriptionLen])
}

// DescriptionTooLong reports whether s exceeds MaxDescriptionLen (after trim).
func DescriptionTooLong(s string) bool {
	return utf8.RuneCountInString(strings.TrimSpace(s)) > MaxDescriptionLen
}

// Parts holds SKILL.md frontmatter and body.
type Parts struct {
	Frontmatter string
	Body        string
}

// Read loads and splits a SKILL.md file.
func Read(path string) (Parts, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Parts{}, err
	}
	m := frontmatterRe.FindSubmatch(raw)
	if m == nil {
		return Parts{}, fmt.Errorf("%s: missing YAML frontmatter", path)
	}
	end := frontmatterRe.FindIndex(raw)
	return Parts{
		Frontmatter: string(m[1]),
		Body:        string(raw[end[1]:]),
	}, nil
}

// findDescriptionLoc returns [start, end) of the description field value block in frontmatter.
func findDescriptionLoc(frontmatter string) (start, end int, ok bool) {
	const key = "description:"
	idx := -1
	lines := strings.Split(frontmatter, "\n")
	offset := 0
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, key) || strings.HasPrefix(line, key) {
			idx = i
			start = offset
			break
		}
		offset += len(line) + 1
	}
	if idx < 0 {
		return 0, 0, false
	}
	// end at next top-level key (line matching ^[a-zA-Z_][\w-]*:) or EOF
	end = len(frontmatter)
	off := start
	for j := idx; j < len(lines); j++ {
		if j > idx {
			t := strings.TrimSpace(lines[j])
			if t != "" && !strings.HasPrefix(lines[j], " ") && !strings.HasPrefix(lines[j], "\t") {
				if isTopLevelKey(lines[j]) {
					end = off
					break
				}
			}
		}
		off += len(lines[j]) + 1
	}
	return start, end, true
}

func isTopLevelKey(line string) bool {
	if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {
		return false
	}
	for i, r := range line {
		if r == ':' {
			return i > 0
		}
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-') {
			if i == 0 {
				return false
			}
			return false
		}
	}
	return false
}

// GetDescription extracts the description field as a single paragraph.
func GetDescription(frontmatter string) (string, error) {
	start, end, ok := findDescriptionLoc(frontmatter)
	if !ok {
		return "", fmt.Errorf("frontmatter has no description field")
	}
	block := frontmatter[start:end]
	// strip "description:"
	rest := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(block), "description:"))
	if strings.HasPrefix(rest, ">") {
		var lines []string
		for _, line := range strings.Split(rest, "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, ">") {
				line = strings.TrimSpace(line[1:])
			}
			if line != "" {
				lines = append(lines, line)
			}
		}
		return strings.Join(lines, " "), nil
	}
	// first line only for inline
	first := strings.SplitN(rest, "\n", 2)[0]
	return strings.Trim(strings.TrimSpace(first), `"'`), nil
}

// SetDescription replaces or appends the description field (folded YAML).
// The value is clamped to MaxDescriptionLen.
func SetDescription(frontmatter, newDescription string) string {
	folded := strings.ReplaceAll(ClampDescription(newDescription), "\n", " ")
	block := fmt.Sprintf("description: >\n  %s\n", folded)
	start, end, ok := findDescriptionLoc(frontmatter)
	if !ok {
		return strings.TrimRight(frontmatter, "\n") + "\n" + block
	}
	rest := strings.TrimLeft(frontmatter[end:], "\n")
	return frontmatter[:start] + block + rest
}

// Write writes SKILL.md from frontmatter and body.
func Write(path string, parts Parts) error {
	body := parts.Body
	if body != "" && !strings.HasPrefix(body, "\n") {
		body = "\n" + body
	}
	content := fmt.Sprintf("---\n%s\n---%s", strings.TrimRight(parts.Frontmatter, "\n"), body)
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	return os.WriteFile(path, []byte(content), 0o644)
}
