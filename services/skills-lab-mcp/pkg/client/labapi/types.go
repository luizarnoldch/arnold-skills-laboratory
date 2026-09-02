package labapi

type Skill struct {
	ID                int64  `json:"id"`
	Name              string `json:"name"`
	IDDescription     *int64 `json:"id_description"`
	IDTestDescription *int64 `json:"id_test_description"`
	IDContent         *int64 `json:"id_content"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

type DescriptionView struct {
	ID               int64  `json:"id"`
	SkillID          int64  `json:"skill_id"`
	SkillDescription string `json:"skill_description"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
	IsCurrent        bool   `json:"is_current,omitempty"`
	IsForTests       bool   `json:"is_for_tests,omitempty"`
}

type ContentView struct {
	ID           int64  `json:"id"`
	SkillID      int64  `json:"skill_id"`
	SkillContent string `json:"skill_content"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type SkillDetail struct {
	Skill
	Description *DescriptionView `json:"description,omitempty"`
	Content     *ContentView     `json:"content,omitempty"`
}

type TriggerQuery struct {
	ID            int64  `json:"id"`
	SkillID       int64  `json:"skill_id"`
	PromptIndex   int64  `json:"prompt_index"`
	Query         string `json:"query"`
	ShouldTrigger bool   `json:"should_trigger"`
	Split         string `json:"split"`
	Runs          *int64 `json:"runs,omitempty"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type BulkItem struct {
	ID            int64  `json:"id"`
	Query         string `json:"query"`
	ShouldTrigger bool   `json:"should_trigger"`
	Runs          *int64 `json:"runs,omitempty"`
}

type SplitRequest struct {
	Seed        int64    `json:"seed"`
	TrainRatio  *float64 `json:"train_ratio,omitempty"`
}

type SplitResponse struct {
	TrainCount      int64          `json:"train_count"`
	ValidationCount int64          `json:"validation_count"`
	Seed            int64          `json:"seed"`
	TrainRatio      float64        `json:"train_ratio"`
	Queries         []TriggerQuery `json:"queries"`
}

type PatchTriggerQueryRequest struct {
	Query         *string  `json:"query,omitempty"`
	ShouldTrigger *bool    `json:"should_trigger,omitempty"`
	Split         *string  `json:"split,omitempty"`
	Runs          *int64   `json:"runs,omitempty"`
}
