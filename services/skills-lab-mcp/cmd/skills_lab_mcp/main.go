package main

import (
	"context"
	"log"
	"os"

	"github.com/luizarnoldch/skills-lab-mcp/pkg/client/labapi"
	"github.com/luizarnoldch/skills-lab-mcp/pkg/client/orch"
	"github.com/luizarnoldch/skills-lab-mcp/pkg/config"
	mcpserver "github.com/luizarnoldch/skills-lab-mcp/pkg/mcp"
	"github.com/luizarnoldch/skills-lab-mcp/pkg/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	log.SetOutput(os.Stderr)

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	lab := labapi.New(cfg.LabAPIURL, cfg.HTTPTimeout)
	orchClient := orch.New(cfg.OrchAPIURL, cfg.HTTPTimeout)

	ctx := context.Background()
	if !lab.CheckReady(ctx) {
		log.Printf("ADVERTENCIA: laboratory-api no responde en %s/ready", cfg.LabAPIURL)
	}
	if !orchClient.CheckReady(ctx) {
		log.Printf("ADVERTENCIA: orchestrator no responde en %s/ready", cfg.OrchAPIURL)
	}

	skills := &tools.Skills{Lab: lab}
	testSets := &tools.TestSets{Lab: lab}
	evals := &tools.Evals{Cfg: cfg, Lab: lab, Orch: orchClient}
	optimize := &tools.Optimize{Cfg: cfg, Lab: lab, Orch: orchClient}
	jobs := &tools.Jobs{Orch: orchClient}

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "skills-lab",
		Version: "0.1.0",
	}, nil)

	mcpserver.Register(server, mcpserver.Deps{
		Skills:   skills,
		TestSets: testSets,
		Evals:    evals,
		Optimize: optimize,
		Jobs:     jobs,
	})

	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		log.Fatalf("server: %v", err)
	}
}
