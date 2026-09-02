package resolve

import (
	"fmt"

	"github.com/luizarnoldch/skills-lab-mcp/pkg/client/labapi"
)

func DescriptionID(skill labapi.Skill, explicit *int64) (int64, error) {
	if explicit != nil {
		return *explicit, nil
	}
	id := EffectiveTestDescriptionID(skill)
	if id == nil {
		return 0, fmt.Errorf("skill %d (%s) no tiene description de prueba ni current", skill.ID, skill.Name)
	}
	return *id, nil
}

func ContentID(skill labapi.Skill, explicit *int64) (int64, error) {
	if explicit != nil {
		return *explicit, nil
	}
	if skill.IDContent == nil {
		return 0, fmt.Errorf("skill %d (%s) no tiene content asignado", skill.ID, skill.Name)
	}
	return *skill.IDContent, nil
}

func StartingDescriptionID(skill labapi.Skill, explicit *int64) (int64, error) {
	return DescriptionID(skill, explicit)
}
