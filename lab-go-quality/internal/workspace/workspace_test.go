package workspace_test

import (
	"os"
	"path/filepath"
	"testing"

	"skills-laboratory/lab-go-quality/internal/workspace"
)

func TestNextIterationAndRunDirs(t *testing.T) {
	base := t.TempDir()
	a, err := workspace.NextIterationDir(base)
	if err != nil {
		t.Fatal(err)
	}
	b, err := workspace.NextIterationDir(base)
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(a) != "iteration-1" || filepath.Base(b) != "iteration-2" {
		t.Fatalf("got %s %s", a, b)
	}
	cfg := filepath.Join(a, "eval-x", "with_skill")
	r1, err := workspace.NextRunDir(cfg)
	if err != nil {
		t.Fatal(err)
	}
	r2, err := workspace.NextRunDir(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(r1) != "run_001" || filepath.Base(r2) != "run_002" {
		t.Fatalf("runs: %s %s", r1, r2)
	}
	if _, err := os.Stat(filepath.Join(r1, "outputs")); err != nil {
		t.Fatal(err)
	}
}

func TestCopyTree(t *testing.T) {
	src := t.TempDir()
	dst := filepath.Join(t.TempDir(), "copy")
	_ = os.WriteFile(filepath.Join(src, "SKILL.md"), []byte("hi\n"), 0o644)
	if err := workspace.CopyTree(src, dst); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(filepath.Join(dst, "SKILL.md"))
	if err != nil || string(raw) != "hi\n" {
		t.Fatalf("copy failed: %v %q", err, raw)
	}
}
