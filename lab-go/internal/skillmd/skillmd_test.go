package skillmd_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"skills-laboratory/lab-go/internal/skillmd"
)

const sample = `---
name: feature-expert
description: >
  Feature expert manager. Manage features.
---
# Body

Hello
`

func TestReadGetSetWrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "SKILL.md")
	if err := os.WriteFile(path, []byte(sample), 0o644); err != nil {
		t.Fatal(err)
	}

	parts, err := skillmd.Read(path)
	if err != nil {
		t.Fatal(err)
	}
	desc, err := skillmd.GetDescription(parts.Frontmatter)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(desc, "Feature expert manager") {
		t.Fatalf("desc=%q", desc)
	}

	parts.Frontmatter = skillmd.SetDescription(parts.Frontmatter, "New description for feature-expert triggers and exclusions clearly stated.")
	if err := skillmd.Write(path, parts); err != nil {
		t.Fatal(err)
	}

	parts2, err := skillmd.Read(path)
	if err != nil {
		t.Fatal(err)
	}
	desc2, err := skillmd.GetDescription(parts2.Frontmatter)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(desc2, "New description for feature-expert") {
		t.Fatalf("desc2=%q", desc2)
	}
	if !strings.Contains(parts2.Body, "# Body") {
		t.Fatalf("body lost: %q", parts2.Body)
	}
}

func TestClampDescription(t *testing.T) {
	short := "Use this skill when managing features."
	if got := skillmd.ClampDescription(short); got != short {
		t.Fatalf("short: %q", got)
	}
	var b strings.Builder
	for i := 0; i < skillmd.MaxDescriptionLen+50; i++ {
		b.WriteRune('á')
	}
	long := b.String()
	if !skillmd.DescriptionTooLong(long) {
		t.Fatal("expected too long")
	}
	clamped := skillmd.ClampDescription(long)
	if n := len([]rune(clamped)); n != skillmd.MaxDescriptionLen {
		t.Fatalf("rune count=%d want %d", n, skillmd.MaxDescriptionLen)
	}
	if skillmd.DescriptionTooLong(clamped) {
		t.Fatal("clamped still reported too long")
	}

	parts := skillmd.Parts{Frontmatter: "name: demo\ndescription: old\n", Body: "\n# x\n"}
	parts.Frontmatter = skillmd.SetDescription(parts.Frontmatter, long)
	desc, err := skillmd.GetDescription(parts.Frontmatter)
	if err != nil {
		t.Fatal(err)
	}
	if skillmd.DescriptionTooLong(desc) {
		t.Fatalf("SetDescription did not clamp: %d runes", len([]rune(desc)))
	}
}
