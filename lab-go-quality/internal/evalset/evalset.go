package evalset

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"
)

// Case is one output-quality eval case.
type Case struct {
	ID             int      `json:"id"`
	Name           string   `json:"name,omitempty"`
	Prompt         string   `json:"prompt"`
	ExpectedOutput string   `json:"expected_output"`
	Files          []string `json:"files,omitempty"`
	Assertions     []string `json:"assertions,omitempty"`
}

// Set is the root evals.json document.
type Set struct {
	SkillName string `json:"skill_name"`
	Evals     []Case `json:"evals"`
}

// Load reads and validates evals.json.
func Load(path string) (Set, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Set{}, err
	}
	var s Set
	if err := json.Unmarshal(raw, &s); err != nil {
		return Set{}, fmt.Errorf("%s: %w", path, err)
	}
	if err := s.Validate(); err != nil {
		return Set{}, fmt.Errorf("%s: %w", path, err)
	}
	return s, nil
}

// Validate checks required fields and unique ids.
func (s Set) Validate() error {
	if strings.TrimSpace(s.SkillName) == "" {
		return fmt.Errorf("skill_name is required")
	}
	if len(s.Evals) == 0 {
		return fmt.Errorf("evals must be non-empty")
	}
	seen := map[int]bool{}
	for i, c := range s.Evals {
		if c.ID < 1 {
			return fmt.Errorf("evals[%d]: id must be >= 1", i)
		}
		if seen[c.ID] {
			return fmt.Errorf("duplicate eval id %d", c.ID)
		}
		seen[c.ID] = true
		if strings.TrimSpace(c.Prompt) == "" {
			return fmt.Errorf("evals[%d]: prompt is required", i)
		}
		if strings.TrimSpace(c.ExpectedOutput) == "" {
			return fmt.Errorf("evals[%d]: expected_output is required", i)
		}
	}
	return nil
}

var nonSlug = regexp.MustCompile(`[^a-z0-9]+`)

// Slug returns eval-<name> or eval-<id>-<prompt-words>.
func Slug(c Case) string {
	if name := strings.TrimSpace(c.Name); name != "" {
		return "eval-" + slugify(name)
	}
	words := firstWords(c.Prompt, 4)
	if words == "" {
		return fmt.Sprintf("eval-%d", c.ID)
	}
	return fmt.Sprintf("eval-%d-%s", c.ID, words)
}

func firstWords(s string, n int) string {
	var parts []string
	var cur strings.Builder
	flush := func() {
		if cur.Len() == 0 {
			return
		}
		parts = append(parts, cur.String())
		cur.Reset()
	}
	for _, r := range strings.ToLower(s) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			cur.WriteRune(r)
			continue
		}
		flush()
		if len(parts) >= n {
			break
		}
	}
	flush()
	if len(parts) > n {
		parts = parts[:n]
	}
	return slugify(strings.Join(parts, "-"))
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = nonSlug.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		return "case"
	}
	if len(s) > 48 {
		s = strings.Trim(s[:48], "-")
	}
	return s
}

// Write writes evals.json with indentation.
func Write(path string, s Set) error {
	if err := s.Validate(); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	raw = append(raw, '\n')
	return os.WriteFile(path, raw, 0o644)
}
