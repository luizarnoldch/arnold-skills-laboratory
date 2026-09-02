package tools_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/luizarnoldch/skills-lab-mcp/pkg/client/labapi"
	"github.com/luizarnoldch/skills-lab-mcp/pkg/tools"
)

func TestSkillsList(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/skills" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode([]labapi.Skill{{ID: 1, Name: "demo"}})
	}))
	defer srv.Close()

	lab := labapi.New(srv.URL, 0)
	skills := &tools.Skills{Lab: lab}

	out, err := skills.List(context.Background(), tools.SkillsListArgs{})
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 1 || out[0].Name != "demo" {
		t.Fatalf("unexpected result: %+v", out)
	}
}
