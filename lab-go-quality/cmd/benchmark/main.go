package main

import (
	"flag"
	"fmt"
	"os"

	"skills-laboratory/lab-go-quality/internal/benchmark"
)

func main() {
	iteration := flag.String("iteration", "", "Path to iteration-N directory")
	flag.Parse()
	if *iteration == "" {
		fmt.Fprintln(os.Stderr, "error: -iteration is required")
		flag.Usage()
		os.Exit(2)
	}
	rep, err := benchmark.Compute(*iteration)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	if err := benchmark.Write(*iteration, rep); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	rs := rep.RunSummary
	fmt.Printf("with_skill  pass_rate=%.3f±%.3f time_s=%.1f tokens=%.0f\n",
		rs.WithSkill.PassRate.Mean, rs.WithSkill.PassRate.Stddev,
		rs.WithSkill.TimeSeconds.Mean, rs.WithSkill.Tokens.Mean)
	fmt.Printf("without_skill pass_rate=%.3f±%.3f time_s=%.1f tokens=%.0f\n",
		rs.WithoutSkill.PassRate.Mean, rs.WithoutSkill.PassRate.Stddev,
		rs.WithoutSkill.TimeSeconds.Mean, rs.WithoutSkill.Tokens.Mean)
	fmt.Printf("delta       pass_rate=%+.3f time_s=%+.1f tokens=%+.0f\n",
		rs.Delta.PassRate, rs.Delta.TimeSeconds, rs.Delta.Tokens)
	fmt.Printf("wrote %s/benchmark.json\n", *iteration)
}
