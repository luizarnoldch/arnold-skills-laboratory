package prompt

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"sort"
)

// Item is one indexed prompt/query for skill trigger evaluation.
type Item struct {
	ID            int    `json:"id"`
	Query         string `json:"query"`
	ShouldTrigger bool   `json:"should_trigger"`
	Runs          *int   `json:"runs,omitempty"`
}

// Load reads a JSON array of prompt items from path.
func Load(path string) ([]Item, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var items []Item
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("%s: expected a non-empty JSON array", path)
	}
	for i, it := range items {
		if it.ID < 1 {
			return nil, fmt.Errorf("%s: item[%d] missing valid id", path, i)
		}
		if it.Query == "" {
			return nil, fmt.Errorf("%s: item[%d] missing query", path, i)
		}
	}
	return items, nil
}

// WriteJSON writes items as indented JSON with trailing newline.
func WriteJSON(path string, items []Item) error {
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}

// StratifiedSplit splits prompts into train/validation by should_trigger.
// trainRatio default 0.6 yields 6/4 for 10 items per class.
func StratifiedSplit(prompts []Item, trainRatio float64, seed int64) (train, validation []Item) {
	if trainRatio <= 0 || trainRatio >= 1 {
		trainRatio = 0.6
	}
	rng := rand.New(rand.NewSource(seed))

	var positives, negatives []Item
	for _, p := range prompts {
		if p.ShouldTrigger {
			positives = append(positives, p)
		} else {
			negatives = append(negatives, p)
		}
	}
	rng.Shuffle(len(positives), func(i, j int) { positives[i], positives[j] = positives[j], positives[i] })
	rng.Shuffle(len(negatives), func(i, j int) { negatives[i], negatives[j] = negatives[j], negatives[i] })

	take := func(items []Item) (tr, val []Item) {
		nTrain := int(float64(len(items))*trainRatio + 0.5) // round
		tr = append([]Item(nil), items[:nTrain]...)
		val = append([]Item(nil), items[nTrain:]...)
		return tr, val
	}

	posTrain, posVal := take(positives)
	negTrain, negVal := take(negatives)

	train = append(posTrain, negTrain...)
	validation = append(posVal, negVal...)
	sort.Slice(train, func(i, j int) bool { return train[i].ID < train[j].ID })
	sort.Slice(validation, func(i, j int) bool { return validation[i].ID < validation[j].ID })
	return train, validation
}

// CountByTrigger returns positive and negative counts.
func CountByTrigger(items []Item) (pos, neg int) {
	for _, it := range items {
		if it.ShouldTrigger {
			pos++
		} else {
			neg++
		}
	}
	return pos, neg
}
