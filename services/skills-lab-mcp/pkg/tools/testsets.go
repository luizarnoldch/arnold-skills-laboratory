package tools

import (
	"context"
	"fmt"

	"github.com/luizarnoldch/skills-lab-mcp/pkg/client/labapi"
	"github.com/luizarnoldch/skills-lab-mcp/pkg/resolve"
)

type PromptItem struct {
	PromptIndex   int64  `json:"prompt_index" jsonschema:"Índice del prompt (1-based)"`
	Query         string `json:"query" jsonschema:"Texto del prompt de prueba"`
	ShouldTrigger bool   `json:"should_trigger" jsonschema:"Si el skill debería activarse"`
	Runs          *int64 `json:"runs,omitempty" jsonschema:"Runs por query (override)"`
}

type TestSetUploadArgs struct {
	resolve.SkillRef
	Prompts    []PromptItem `json:"prompts" jsonschema:"Lista de prompts a cargar"`
	TrainRatio float64      `json:"train_ratio,omitempty" jsonschema:"Proporción train (default 0.6)"`
	Seed       int64          `json:"seed,omitempty" jsonschema:"Seed para split estratificado (default 42)"`
}

type TestSetUploadResult struct {
	SkillID         int64                `json:"skill_id"`
	SkillName       string               `json:"skill_name"`
	TrainCount      int64                `json:"train_count"`
	ValidationCount int64                `json:"validation_count"`
	Seed            int64                `json:"seed"`
	TrainRatio      float64              `json:"train_ratio"`
	Prompts         []labapi.TriggerQuery `json:"prompts"`
}

type TestSetListArgs struct {
	resolve.SkillRef
	Split string `json:"split,omitempty" jsonschema:"Filtrar: unassigned, train, validation, all"`
}

type TestSetListResult struct {
	SkillID     int64                `json:"skill_id"`
	SkillName   string               `json:"skill_name"`
	SplitFilter string               `json:"split_filter"`
	Prompts     []labapi.TriggerQuery `json:"prompts"`
}

type TestSetUpdateArgs struct {
	TriggerQueryID *int64  `json:"trigger_query_id,omitempty" jsonschema:"ID del trigger query"`
	SkillID        *int64  `json:"skill_id,omitempty"`
	SkillName      *string `json:"skill_name,omitempty"`
	PromptIndex    *int64  `json:"prompt_index,omitempty" jsonschema:"Índice si no se usa trigger_query_id"`
	Query          *string `json:"query,omitempty"`
	ShouldTrigger  *bool   `json:"should_trigger,omitempty"`
	Split          *string `json:"split,omitempty" jsonschema:"train, validation o unassigned"`
	Runs           *int64  `json:"runs,omitempty"`
}

type TestSets struct {
	Lab *labapi.Client
}

func (t *TestSets) Upload(ctx context.Context, args TestSetUploadArgs) (TestSetUploadResult, error) {
	if err := args.Validate(); err != nil {
		return TestSetUploadResult{}, err
	}
	if len(args.Prompts) == 0 {
		return TestSetUploadResult{}, fmt.Errorf("se requiere al menos un prompt")
	}

	trainRatio := args.TrainRatio
	if trainRatio == 0 {
		trainRatio = 0.6
	}
	seed := args.Seed
	if seed == 0 {
		seed = 42
	}

	skill, err := resolve.ResolveSkill(ctx, t.Lab, args.SkillRef)
	if err != nil {
		return TestSetUploadResult{}, err
	}

	items := make([]labapi.BulkItem, len(args.Prompts))
	for i, p := range args.Prompts {
		items[i] = labapi.BulkItem{
			ID:            p.PromptIndex,
			Query:         p.Query,
			ShouldTrigger: p.ShouldTrigger,
			Runs:          p.Runs,
		}
	}

	if _, err := t.Lab.BulkReplaceTriggerQueries(ctx, skill.ID, items); err != nil {
		return TestSetUploadResult{}, err
	}

	splitResult, err := t.Lab.SplitTriggerQueries(ctx, skill.ID, labapi.SplitRequest{
		Seed:       seed,
		TrainRatio: &trainRatio,
	})
	if err != nil {
		return TestSetUploadResult{}, err
	}

	return TestSetUploadResult{
		SkillID:         skill.ID,
		SkillName:       skill.Name,
		TrainCount:      splitResult.TrainCount,
		ValidationCount: splitResult.ValidationCount,
		Seed:            splitResult.Seed,
		TrainRatio:      splitResult.TrainRatio,
		Prompts:         splitResult.Queries,
	}, nil
}

func (t *TestSets) List(ctx context.Context, args TestSetListArgs) (TestSetListResult, error) {
	if err := args.Validate(); err != nil {
		return TestSetListResult{}, err
	}

	skill, err := resolve.ResolveSkill(ctx, t.Lab, args.SkillRef)
	if err != nil {
		return TestSetListResult{}, err
	}

	splitFilter := args.Split
	if splitFilter == "" {
		splitFilter = "all"
	}
	split := splitFilter
	if split == "all" {
		split = ""
	}

	prompts, err := t.Lab.ListTriggerQueries(ctx, skill.ID, split)
	if err != nil {
		return TestSetListResult{}, err
	}

	return TestSetListResult{
		SkillID:     skill.ID,
		SkillName:   skill.Name,
		SplitFilter: splitFilter,
		Prompts:     prompts,
	}, nil
}

func (t *TestSets) Update(ctx context.Context, args TestSetUpdateArgs) (labapi.TriggerQuery, error) {
	queryID := args.TriggerQueryID

	if queryID == nil {
		if args.PromptIndex == nil {
			return labapi.TriggerQuery{}, fmt.Errorf("se requiere trigger_query_id o prompt_index")
		}
		ref := resolve.SkillRef{SkillID: args.SkillID, SkillName: args.SkillName}
		if err := ref.Validate(); err != nil {
			return labapi.TriggerQuery{}, err
		}
		skill, err := resolve.ResolveSkill(ctx, t.Lab, ref)
		if err != nil {
			return labapi.TriggerQuery{}, err
		}
		prompts, err := t.Lab.ListTriggerQueries(ctx, skill.ID, "")
		if err != nil {
			return labapi.TriggerQuery{}, err
		}
		var found *labapi.TriggerQuery
		for i := range prompts {
			if prompts[i].PromptIndex == *args.PromptIndex {
				found = &prompts[i]
				break
			}
		}
		if found == nil {
			return labapi.TriggerQuery{}, fmt.Errorf("no se encontró prompt con prompt_index=%d", *args.PromptIndex)
		}
		queryID = &found.ID
	}

	if args.Query == nil && args.ShouldTrigger == nil && args.Split == nil && args.Runs == nil {
		return labapi.TriggerQuery{}, fmt.Errorf("se requiere al menos un campo a actualizar: query, should_trigger, split, runs")
	}

	patch := labapi.PatchTriggerQueryRequest{
		Query:         args.Query,
		ShouldTrigger: args.ShouldTrigger,
		Split:         args.Split,
		Runs:          args.Runs,
	}
	return t.Lab.PatchTriggerQuery(ctx, *queryID, patch)
}
