package resolve

import (
	"context"
	"fmt"

	"github.com/luizarnoldch/skills-lab-mcp/pkg/client/labapi"
)

type SkillRef struct {
	SkillID   *int64  `json:"skill_id,omitempty" jsonschema:"ID numérico del skill"`
	SkillName *string `json:"skill_name,omitempty" jsonschema:"Nombre exacto del skill"`
}

func (r SkillRef) Validate() error {
	if r.SkillID == nil && (r.SkillName == nil || *r.SkillName == "") {
		return fmt.Errorf("se requiere skill_id o skill_name")
	}
	return nil
}

func EffectiveTestDescriptionID(skill labapi.Skill) *int64 {
	if skill.IDTestDescription != nil {
		return skill.IDTestDescription
	}
	return skill.IDDescription
}

func ResolveSkill(ctx context.Context, lab *labapi.Client, ref SkillRef) (labapi.Skill, error) {
	if err := ref.Validate(); err != nil {
		return labapi.Skill{}, err
	}

	if ref.SkillID != nil {
		detail, err := lab.GetSkill(ctx, *ref.SkillID)
		if err != nil {
			return labapi.Skill{}, err
		}
		return detail.Skill, nil
	}

	skills, err := lab.ListSkills(ctx)
	if err != nil {
		return labapi.Skill{}, err
	}

	name := *ref.SkillName
	var matches []labapi.Skill
	for _, s := range skills {
		if s.Name == name {
			matches = append(matches, s)
		}
	}
	switch len(matches) {
	case 0:
		return labapi.Skill{}, fmt.Errorf("skill no encontrado con nombre %q", name)
	case 1:
		return matches[0], nil
	default:
		ids := make([]int64, len(matches))
		for i, s := range matches {
			ids[i] = s.ID
		}
		return labapi.Skill{}, fmt.Errorf("nombre %q es ambiguo (%d skills); usa skill_id: %v", name, len(matches), ids)
	}
}

func ResolveSkillDetail(ctx context.Context, lab *labapi.Client, ref SkillRef) (labapi.SkillDetail, error) {
	skill, err := ResolveSkill(ctx, lab, ref)
	if err != nil {
		return labapi.SkillDetail{}, err
	}
	return lab.GetSkill(ctx, skill.ID)
}
