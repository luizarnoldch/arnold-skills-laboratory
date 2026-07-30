package skillmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var frontmatterRe = regexp.MustCompile(`(?s)^---\n(.*?)\n---\n?`)

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

// ReadRaw returns the full file contents.
func ReadRaw(path string) (string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(raw), nil
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

// WriteRaw writes full SKILL.md contents.
func WriteRaw(path, content string) error {
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	return os.WriteFile(path, []byte(content), 0o644)
}
