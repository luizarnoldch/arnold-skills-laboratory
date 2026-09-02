package tools

import (
	"context"

	"github.com/luizarnoldch/skills-lab-mcp/pkg/client/labapi"
	"github.com/luizarnoldch/skills-lab-mcp/pkg/client/orch"
	"github.com/luizarnoldch/skills-lab-mcp/pkg/config"
	"github.com/luizarnoldch/skills-lab-mcp/pkg/resolve"
)

type EvalCommonArgs struct {
	resolve.SkillRef
	DescriptionID     *int64   `json:"description_id,omitempty" jsonschema:"Description a evaluar (default: test)"`
	ContentID         *int64   `json:"content_id,omitempty" jsonschema:"Content (default: current)"`
	RunsDefault       *int64   `json:"runs_default,omitempty"`
	MajorityThreshold *float64 `json:"majority_threshold,omitempty"`
	Workspace         string   `json:"workspace,omitempty"`
	Model             string   `json:"model,omitempty"`
	Provider          string   `json:"provider,omitempty" jsonschema:"deepseek o cursor"`
	MaxTurns          int      `json:"max_turns,omitempty"`
	TimeoutMs         int      `json:"timeout_ms,omitempty"`
	MaxParallel       *int     `json:"max_parallel,omitempty"`
	StopOnSkillCall   *bool    `json:"stop_on_skill_call,omitempty"`
}

type BaselineEvalStartArgs struct {
	EvalCommonArgs
}

type TriggerEvalStartArgs struct {
	EvalCommonArgs
	Split string `json:"split" jsonschema:"train, validation o all"`
}

type JobStartResult struct {
	JobID         string `json:"job_id"`
	JobType       string `json:"job_type"`
	Status        string `json:"status"`
	SkillID       int64  `json:"skill_id"`
	DescriptionID int64  `json:"description_id"`
	ContentID     int64  `json:"content_id"`
	Split         string `json:"split,omitempty"`
}

type Evals struct {
	Cfg  config.Config
	Lab  *labapi.Client
	Orch *orch.Client
}

func (e *Evals) BaselineStart(ctx context.Context, args BaselineEvalStartArgs) (JobStartResult, error) {
	if err := args.Validate(); err != nil {
		return JobStartResult{}, err
	}

	skill, err := resolve.ResolveSkill(ctx, e.Lab, args.SkillRef)
	if err != nil {
		return JobStartResult{}, err
	}
	descID, err := resolve.DescriptionID(skill, args.DescriptionID)
	if err != nil {
		return JobStartResult{}, err
	}
	contentID, err := resolve.ContentID(skill, args.ContentID)
	if err != nil {
		return JobStartResult{}, err
	}

	workspace := args.Workspace
	if workspace == "" {
		workspace = e.Cfg.WorkspaceDefault
	}

	resp, err := e.Orch.CreateBaselineEvalJob(ctx, orch.CreateBaselineEvalJobRequest{
		SkillID:           skill.ID,
		DescriptionID:     descID,
		ContentID:         contentID,
		RunsDefault:       args.RunsDefault,
		MajorityThreshold: args.MajorityThreshold,
		Workspace:         workspace,
		Model:             args.Model,
		Provider:          args.Provider,
		MaxTurns:          args.MaxTurns,
		TimeoutMs:         args.TimeoutMs,
		MaxParallel:       args.MaxParallel,
		StopOnSkillCall:   args.StopOnSkillCall,
	})
	if err != nil {
		return JobStartResult{}, err
	}

	return JobStartResult{
		JobID:         resp.ID,
		JobType:       "baseline_eval",
		Status:        resp.Status,
		SkillID:       skill.ID,
		DescriptionID: descID,
		ContentID:     contentID,
	}, nil
}

func (e *Evals) TriggerStart(ctx context.Context, args TriggerEvalStartArgs) (JobStartResult, error) {
	if err := args.Validate(); err != nil {
		return JobStartResult{}, err
	}
	if args.Split == "" {
		return JobStartResult{}, errSplitRequired
	}

	skill, err := resolve.ResolveSkill(ctx, e.Lab, args.SkillRef)
	if err != nil {
		return JobStartResult{}, err
	}
	descID, err := resolve.DescriptionID(skill, args.DescriptionID)
	if err != nil {
		return JobStartResult{}, err
	}
	contentID, err := resolve.ContentID(skill, args.ContentID)
	if err != nil {
		return JobStartResult{}, err
	}

	workspace := args.Workspace
	if workspace == "" {
		workspace = e.Cfg.WorkspaceDefault
	}

	resp, err := e.Orch.CreateTriggerEvalBatch(ctx, orch.CreateTriggerEvalBatchRequest{
		SkillID:           skill.ID,
		DescriptionID:     descID,
		ContentID:         contentID,
		Split:             args.Split,
		RunsDefault:       args.RunsDefault,
		MajorityThreshold: args.MajorityThreshold,
		Workspace:         workspace,
		Model:             args.Model,
		Provider:          args.Provider,
		MaxTurns:          args.MaxTurns,
		TimeoutMs:         args.TimeoutMs,
		MaxParallel:       args.MaxParallel,
		StopOnSkillCall:   args.StopOnSkillCall,
	})
	if err != nil {
		return JobStartResult{}, err
	}

	return JobStartResult{
		JobID:         resp.ID,
		JobType:       "trigger_eval",
		Status:        resp.Status,
		SkillID:       skill.ID,
		DescriptionID: descID,
		ContentID:     contentID,
		Split:         args.Split,
	}, nil
}

var errSplitRequired = &evalError{msg: "split es requerido (train, validation o all)"}

type evalError struct{ msg string }

func (e *evalError) Error() string { return e.msg }
