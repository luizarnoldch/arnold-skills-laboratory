package tools

import (
	"context"

	"github.com/luizarnoldch/skills-lab-mcp/pkg/client/labapi"
	"github.com/luizarnoldch/skills-lab-mcp/pkg/resolve"
)

type SkillsListArgs struct{}

type SkillGetArgs struct {
	resolve.SkillRef
}

type SkillGetResult struct {
	Skill                labapi.Skill            `json:"skill"`
	Current              SkillCurrent            `json:"current"`
	TestDescriptionID    *int64                  `json:"test_description_id"`
	DescriptionVersions  []labapi.DescriptionView `json:"description_versions"`
	ContentVersions      []labapi.ContentView     `json:"content_versions"`
}

type SkillCurrent struct {
	Description *labapi.DescriptionView `json:"description"`
	Content     *labapi.ContentView     `json:"content"`
}

type Skills struct {
	Lab *labapi.Client
}

func (s *Skills) List(ctx context.Context, _ SkillsListArgs) ([]labapi.Skill, error) {
	return s.Lab.ListSkills(ctx)
}

func (s *Skills) Get(ctx context.Context, args SkillGetArgs) (SkillGetResult, error) {
	if err := args.Validate(); err != nil {
		return SkillGetResult{}, err
	}

	detail, err := resolve.ResolveSkillDetail(ctx, s.Lab, args.SkillRef)
	if err != nil {
		return SkillGetResult{}, err
	}

	descs, err := s.Lab.ListDescriptions(ctx, detail.ID)
	if err != nil {
		return SkillGetResult{}, err
	}
	contents, err := s.Lab.ListContents(ctx, detail.ID)
	if err != nil {
		return SkillGetResult{}, err
	}

	return SkillGetResult{
		Skill:               detail.Skill,
		Current:             SkillCurrent{Description: detail.Description, Content: detail.Content},
		TestDescriptionID:   resolve.EffectiveTestDescriptionID(detail.Skill),
		DescriptionVersions: descs,
		ContentVersions:     contents,
	}, nil
}
