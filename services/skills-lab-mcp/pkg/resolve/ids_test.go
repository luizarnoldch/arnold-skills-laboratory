package resolve_test

import (
	"testing"

	"github.com/luizarnoldch/skills-lab-mcp/pkg/client/labapi"
	"github.com/luizarnoldch/skills-lab-mcp/pkg/resolve"
)

func TestEffectiveTestDescriptionID(t *testing.T) {
	testID := int64(5)
	currentID := int64(3)

	skill := labapi.Skill{IDTestDescription: &testID, IDDescription: &currentID}
	if got := resolve.EffectiveTestDescriptionID(skill); got == nil || *got != 5 {
		t.Fatalf("expected test id 5, got %v", got)
	}

	skill = labapi.Skill{IDDescription: &currentID}
	if got := resolve.EffectiveTestDescriptionID(skill); got == nil || *got != 3 {
		t.Fatalf("expected fallback to current id 3, got %v", got)
	}
}

func TestDescriptionID(t *testing.T) {
	desc := int64(7)
	skill := labapi.Skill{ID: 1, Name: "x", IDDescription: &desc}

	got, err := resolve.DescriptionID(skill, nil)
	if err != nil || got != 7 {
		t.Fatalf("got %d err %v", got, err)
	}

	explicit := int64(9)
	got, err = resolve.DescriptionID(skill, &explicit)
	if err != nil || got != 9 {
		t.Fatalf("got %d err %v", got, err)
	}
}

func TestContentID(t *testing.T) {
	content := int64(4)
	skill := labapi.Skill{ID: 1, Name: "x", IDContent: &content}

	got, err := resolve.ContentID(skill, nil)
	if err != nil || got != 4 {
		t.Fatalf("got %d err %v", got, err)
	}

	_, err = resolve.ContentID(labapi.Skill{ID: 1, Name: "x"}, nil)
	if err == nil {
		t.Fatal("expected error when no content")
	}
}

func TestSkillRefValidate(t *testing.T) {
	name := "foo"
	ref := resolve.SkillRef{SkillName: &name}
	if err := ref.Validate(); err != nil {
		t.Fatal(err)
	}

	if err := (resolve.SkillRef{}).Validate(); err == nil {
		t.Fatal("expected validation error")
	}
}
