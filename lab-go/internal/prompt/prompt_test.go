package prompt_test

import (
	"path/filepath"
	"testing"

	"skills-laboratory/lab-go/internal/prompt"
)

func TestStratifiedSplit(t *testing.T) {
	var items []prompt.Item
	for i := 1; i <= 10; i++ {
		items = append(items, prompt.Item{ID: i, Query: "pos", ShouldTrigger: true})
	}
	for i := 11; i <= 20; i++ {
		items = append(items, prompt.Item{ID: i, Query: "neg", ShouldTrigger: false})
	}

	train, val := prompt.StratifiedSplit(items, 0.6, 42)
	if len(train) != 12 || len(val) != 8 {
		t.Fatalf("sizes: train=%d val=%d want 12/8", len(train), len(val))
	}
	posT, negT := prompt.CountByTrigger(train)
	posV, negV := prompt.CountByTrigger(val)
	if posT != 6 || negT != 6 || posV != 4 || negV != 4 {
		t.Fatalf("stratification train=%d/%d val=%d/%d", posT, negT, posV, negV)
	}

	seen := map[int]bool{}
	for _, it := range append(append([]prompt.Item{}, train...), val...) {
		if seen[it.ID] {
			t.Fatalf("duplicate id %d", it.ID)
		}
		seen[it.ID] = true
	}
	if len(seen) != 20 {
		t.Fatalf("expected 20 unique ids, got %d", len(seen))
	}
}

func TestLoadAndWrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "prompts.json")
	items := []prompt.Item{
		{ID: 1, Query: "a", ShouldTrigger: true},
		{ID: 2, Query: "b", ShouldTrigger: false},
	}
	if err := prompt.WriteJSON(path, items); err != nil {
		t.Fatal(err)
	}
	got, err := prompt.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0].ID != 1 || got[1].ShouldTrigger {
		t.Fatalf("unexpected load: %+v", got)
	}
}
