package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"skills-laboratory/lab-go/internal/prompt"
)

func main() {
	input := flag.String("i", "", "Full prompts.json")
	trainOut := flag.String("train-out", "", "Output train.json (default: sibling train.json)")
	valOut := flag.String("validation-out", "", "Output validation.json (default: sibling validation.json)")
	trainRatio := flag.Float64("train-ratio", 0.6, "Train fraction (default 0.6)")
	seed := flag.Int64("seed", 42, "RNG seed for reproducibility")
	flag.Parse()

	if *input == "" {
		fmt.Fprintln(os.Stderr, "error: -i is required")
		flag.Usage()
		os.Exit(2)
	}

	items, err := prompt.Load(*input)
	if err != nil {
		fatal(err)
	}

	train, validation := prompt.StratifiedSplit(items, *trainRatio, *seed)

	tOut := *trainOut
	if tOut == "" {
		tOut = filepath.Join(filepath.Dir(*input), "train.json")
	}
	vOut := *valOut
	if vOut == "" {
		vOut = filepath.Join(filepath.Dir(*input), "validation.json")
	}

	if err := os.MkdirAll(filepath.Dir(tOut), 0o755); err != nil {
		fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(vOut), 0o755); err != nil {
		fatal(err)
	}
	if err := prompt.WriteJSON(tOut, train); err != nil {
		fatal(err)
	}
	if err := prompt.WriteJSON(vOut, validation); err != nil {
		fatal(err)
	}

	posT, negT := prompt.CountByTrigger(train)
	posV, negV := prompt.CountByTrigger(validation)
	fmt.Printf("Wrote %s (%d = %d pos / %d neg)\n", tOut, len(train), posT, negT)
	fmt.Printf("Wrote %s (%d = %d pos / %d neg)\n", vOut, len(validation), posV, negV)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}
