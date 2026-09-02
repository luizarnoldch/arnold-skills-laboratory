package tools

import (
	"context"

	"github.com/luizarnoldch/skills-lab-mcp/pkg/client/labapi"
	"github.com/luizarnoldch/skills-lab-mcp/pkg/client/orch"
	"github.com/luizarnoldch/skills-lab-mcp/pkg/config"
	"github.com/luizarnoldch/skills-lab-mcp/pkg/resolve"
)

type OptimizeStartArgs struct {
	resolve.SkillRef
	StartingDescriptionID *int64   `json:"starting_description_id,omitempty"`
	ContentID             *int64   `json:"content_id,omitempty"`
	MaxIters              *int64   `json:"max_iters,omitempty"`
	Threshold             *float64 `json:"threshold,omitempty"`
	RunsPerQuery          *int64   `json:"runs_per_query,omitempty"`
	MajorityThreshold     *float64 `json:"majority_threshold,omitempty"`
	Workspace             string   `json:"workspace,omitempty"`
	Model                 string   `json:"model,omitempty"`
	Provider              string   `json:"provider,omitempty" jsonschema:"deepseek o cursor"`
	MaxTurns              int      `json:"max_turns,omitempty"`
	TimeoutMs             int      `json:"timeout_ms,omitempty"`
	MaxParallel           *int     `json:"max_parallel,omitempty"`
}

type OptimizeStartResult struct {
	JobID                 string `json:"job_id"`
	JobType               string `json:"job_type"`
	Status                string `json:"status"`
	SkillID               int64  `json:"skill_id"`
	StartingDescriptionID int64  `json:"starting_description_id"`
	ContentID             int64  `json:"content_id"`
}

type Optimize struct {
	Cfg  config.Config
	Lab  *labapi.Client
	Orch *orch.Client
}

func (o *Optimize) Start(ctx context.Context, args OptimizeStartArgs) (OptimizeStartResult, error) {
	if err := args.Validate(); err != nil {
		return OptimizeStartResult{}, err
	}

	skill, err := resolve.ResolveSkill(ctx, o.Lab, args.SkillRef)
	if err != nil {
		return OptimizeStartResult{}, err
	}
	startDescID, err := resolve.StartingDescriptionID(skill, args.StartingDescriptionID)
	if err != nil {
		return OptimizeStartResult{}, err
	}
	contentID, err := resolve.ContentID(skill, args.ContentID)
	if err != nil {
		return OptimizeStartResult{}, err
	}

	workspace := args.Workspace
	if workspace == "" {
		workspace = o.Cfg.WorkspaceDefault
	}

	resp, err := o.Orch.CreateOptimizeJob(ctx, orch.CreateOptimizeJobRequest{
		SkillID:               skill.ID,
		StartingDescriptionID: startDescID,
		ContentID:             contentID,
		MaxIters:              args.MaxIters,
		Threshold:             args.Threshold,
		RunsPerQuery:          args.RunsPerQuery,
		MajorityThreshold:     args.MajorityThreshold,
		Workspace:             workspace,
		Model:                 args.Model,
		Provider:              args.Provider,
		MaxTurns:              args.MaxTurns,
		TimeoutMs:             args.TimeoutMs,
		MaxParallel:           args.MaxParallel,
	})
	if err != nil {
		return OptimizeStartResult{}, err
	}

	return OptimizeStartResult{
		JobID:                 resp.ID,
		JobType:               "optimize",
		Status:                resp.Status,
		SkillID:               skill.ID,
		StartingDescriptionID: startDescID,
		ContentID:             contentID,
	}, nil
}
