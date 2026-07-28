package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"skills-laboratory/lab-go-quality/internal/evalset"
	"skills-laboratory/lab-go-quality/internal/workspace"
)

func main() {
	skillName := flag.String("skill-name", "", "Skill name")
	skillDir := flag.String("skill-dir", "", "Canonical skill directory (contains or will contain SKILL.md)")
	workspaceRoot := flag.String("workspace-root", "workspace/quality", "Root for quality workspaces")
	force := flag.Bool("force", false, "Overwrite evals.json if it exists")
	flag.Parse()

	if *skillName == "" || *skillDir == "" {
		fmt.Fprintln(os.Stderr, "error: -skill-name and -skill-dir are required")
		flag.Usage()
		os.Exit(2)
	}

	evalsDir := filepath.Join(*skillDir, "evals")
	filesDir := filepath.Join(evalsDir, "files")
	must(workspace.EnsureDir(filesDir))

	evalsPath := filepath.Join(evalsDir, "evals.json")
	if _, err := os.Stat(evalsPath); err == nil && !*force {
		fmt.Printf("exists: %s (use -force to overwrite)\n", evalsPath)
	} else {
		set := evalset.Set{
			SkillName: *skillName,
			Evals: []evalset.Case{
				{
					ID:             1,
					Name:           "happy-path",
					Prompt:         "TODO: realistic user prompt for a successful skill run",
					ExpectedOutput: "TODO: human-readable description of success",
					Files:          []string{},
					Assertions: []string{
						"TODO: verifiable assertion about the output",
					},
				},
				{
					ID:             2,
					Name:           "edge-case",
					Prompt:         "TODO: edge-case or ambiguous request",
					ExpectedOutput: "TODO: expected behavior on the edge case",
					Assertions: []string{
						"TODO: assertion for the edge case",
					},
				},
			},
		}
		must(evalset.Write(evalsPath, set))
		fmt.Printf("wrote %s\n", evalsPath)
	}

	ws := filepath.Join(*workspaceRoot, *skillName)
	must(workspace.EnsureDir(ws))
	gitignore := filepath.Join(ws, ".gitignore")
	if _, err := os.Stat(gitignore); os.IsNotExist(err) {
		_ = os.WriteFile(gitignore, []byte("iteration-*/\n"), 0o644)
	}
	fmt.Printf("workspace: %s\n", ws)
	fmt.Println("edit evals.json, then: go -C lab-go-quality run ./cmd/runevals ...")
}

func must(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
