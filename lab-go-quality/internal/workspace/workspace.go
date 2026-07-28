package workspace

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

var iterationRe = regexp.MustCompile(`^iteration-(\d+)$`)
var runRe = regexp.MustCompile(`^run_(\d+)$`)

// NextIterationDir creates iteration-N under base (iteration-1, iteration-2, ...).
func NextIterationDir(base string) (string, error) {
	if err := os.MkdirAll(base, 0o755); err != nil {
		return "", err
	}
	entries, err := os.ReadDir(base)
	if err != nil {
		return "", err
	}
	maxN := 0
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		m := iterationRe.FindStringSubmatch(e.Name())
		if len(m) != 2 {
			continue
		}
		n, err := strconv.Atoi(m[1])
		if err == nil && n > maxN {
			maxN = n
		}
	}
	next := filepath.Join(base, fmt.Sprintf("iteration-%d", maxN+1))
	if err := os.MkdirAll(next, 0o755); err != nil {
		return "", err
	}
	return next, nil
}

// NextRunDir creates the next run_NNN directory under configDir.
func NextRunDir(configDir string) (string, error) {
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return "", err
	}
	entries, err := os.ReadDir(configDir)
	if err != nil {
		return "", err
	}
	maxRun := 0
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		m := runRe.FindStringSubmatch(e.Name())
		if len(m) != 2 {
			continue
		}
		n, err := strconv.Atoi(m[1])
		if err == nil && n > maxRun {
			maxRun = n
		}
	}
	next := filepath.Join(configDir, fmt.Sprintf("run_%03d", maxRun+1))
	if err := os.MkdirAll(filepath.Join(next, "outputs"), 0o755); err != nil {
		return "", err
	}
	return next, nil
}

// RelPath returns path relative to iterationDir; falls back to path on error.
func RelPath(iterationDir, path string) string {
	rel, err := filepath.Rel(iterationDir, path)
	if err != nil {
		return path
	}
	return rel
}

// CopyTree copies src directory to dst (recursive).
func CopyTree(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return copyFile(src, dst)
	}
	if err := os.MkdirAll(dst, info.Mode()); err != nil {
		return err
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, e := range entries {
		s := filepath.Join(src, e.Name())
		d := filepath.Join(dst, e.Name())
		if e.Type()&os.ModeSymlink != 0 {
			target, err := os.Readlink(s)
			if err != nil {
				return err
			}
			if err := os.Symlink(target, d); err != nil {
				return err
			}
			continue
		}
		if e.IsDir() {
			if err := CopyTree(s, d); err != nil {
				return err
			}
			continue
		}
		if err := copyFile(s, d); err != nil {
			return err
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	info, err := in.Stat()
	if err != nil {
		return err
	}
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

// EnsureDir creates dir if missing.
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0o755)
}
