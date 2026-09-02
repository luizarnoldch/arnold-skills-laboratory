package orch

type CreateBaselineEvalJobRequest struct {
	SkillID           int64    `json:"skill_id"`
	DescriptionID     int64    `json:"description_id"`
	ContentID         int64    `json:"content_id"`
	RunsDefault       *int64   `json:"runs_default,omitempty"`
	MajorityThreshold *float64 `json:"majority_threshold,omitempty"`
	Workspace         string   `json:"workspace,omitempty"`
	Model             string   `json:"model,omitempty"`
	Provider          string   `json:"provider,omitempty"`
	MaxTurns          int      `json:"max_turns,omitempty"`
	TimeoutMs         int      `json:"timeout_ms,omitempty"`
	MaxParallel       *int     `json:"max_parallel,omitempty"`
	StopOnSkillCall   *bool    `json:"stop_on_skill_call,omitempty"`
}

type CreateTriggerEvalBatchRequest struct {
	SkillID           int64    `json:"skill_id"`
	DescriptionID     int64    `json:"description_id"`
	ContentID         int64    `json:"content_id"`
	Split             string   `json:"split"`
	RunsDefault       *int64   `json:"runs_default,omitempty"`
	MajorityThreshold *float64 `json:"majority_threshold,omitempty"`
	Workspace         string   `json:"workspace,omitempty"`
	Model             string   `json:"model,omitempty"`
	Provider          string   `json:"provider,omitempty"`
	MaxTurns          int      `json:"max_turns,omitempty"`
	TimeoutMs         int      `json:"timeout_ms,omitempty"`
	MaxParallel       *int     `json:"max_parallel,omitempty"`
	StopOnSkillCall   *bool    `json:"stop_on_skill_call,omitempty"`
}

type CreateOptimizeJobRequest struct {
	SkillID               int64    `json:"skill_id"`
	StartingDescriptionID int64    `json:"starting_description_id"`
	ContentID             int64    `json:"content_id"`
	MaxIters              *int64   `json:"max_iters,omitempty"`
	Threshold             *float64 `json:"threshold,omitempty"`
	RunsPerQuery          *int64   `json:"runs_per_query,omitempty"`
	MajorityThreshold     *float64 `json:"majority_threshold,omitempty"`
	Workspace             string   `json:"workspace,omitempty"`
	Model                 string   `json:"model,omitempty"`
	Provider              string   `json:"provider,omitempty"`
	MaxTurns              int      `json:"max_turns,omitempty"`
	TimeoutMs             int      `json:"timeout_ms,omitempty"`
	MaxParallel           *int     `json:"max_parallel,omitempty"`
}

type JobCreateResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type BaselineEvalJob struct {
	ID                    string   `json:"id"`
	SkillID               int64    `json:"skill_id"`
	DescriptionID         int64    `json:"description_id"`
	ContentID             int64    `json:"content_id"`
	SkillName             string   `json:"skill_name"`
	RunsDefault           int64    `json:"runs_default"`
	MajorityThreshold     float64  `json:"majority_threshold"`
	Workspace             string   `json:"workspace"`
	Model                 string   `json:"model,omitempty"`
	Provider              string   `json:"provider,omitempty"`
	MaxParallel           int      `json:"max_parallel"`
	Status                string   `json:"status"`
	ProgressDone          int      `json:"progress_done"`
	ProgressTotal         int      `json:"progress_total"`
	TrainEvalRunID        *int64   `json:"train_eval_run_id,omitempty"`
	ValidationEvalRunID   *int64   `json:"validation_eval_run_id,omitempty"`
	TrainAccuracy         *float64 `json:"train_accuracy,omitempty"`
	ValidationAccuracy    *float64 `json:"validation_accuracy,omitempty"`
	Error                 string   `json:"error,omitempty"`
	CreatedAt             string   `json:"created_at"`
	UpdatedAt             string   `json:"updated_at"`
	CompletedAt           string   `json:"completed_at,omitempty"`
}

type TriggerEvalBatch struct {
	ID                string   `json:"id"`
	SkillID           int64    `json:"skill_id"`
	DescriptionID     int64    `json:"description_id"`
	ContentID         int64    `json:"content_id"`
	SkillName         string   `json:"skill_name"`
	Split             string   `json:"split"`
	RunsDefault       int64    `json:"runs_default"`
	MajorityThreshold float64  `json:"majority_threshold"`
	Workspace         string   `json:"workspace"`
	Model             string   `json:"model,omitempty"`
	Provider          string   `json:"provider,omitempty"`
	Status            string   `json:"status"`
	ProgressDone      int      `json:"progress_done"`
	ProgressTotal     int      `json:"progress_total"`
	LabEvalRunID      *int64   `json:"lab_eval_run_id,omitempty"`
	Accuracy          *float64 `json:"accuracy,omitempty"`
	Correct           *int64   `json:"correct,omitempty"`
	Total             *int64   `json:"total,omitempty"`
	Error             string   `json:"error,omitempty"`
	CreatedAt         string   `json:"created_at"`
	UpdatedAt         string   `json:"updated_at"`
	CompletedAt       string   `json:"completed_at,omitempty"`
}

type OptimizeJob struct {
	ID                    string  `json:"id"`
	SkillID               int64   `json:"skill_id"`
	ContentID             int64   `json:"content_id"`
	StartingDescriptionID int64   `json:"starting_description_id"`
	LabOptimizeRunID      *int64  `json:"lab_optimize_run_id,omitempty"`
	SkillName             string  `json:"skill_name"`
	Workspace             string  `json:"workspace"`
	Model                 string  `json:"model,omitempty"`
	Provider              string  `json:"provider,omitempty"`
	MaxIters              int64   `json:"max_iters"`
	Threshold             float64 `json:"threshold"`
	RunsPerQuery          int64   `json:"runs_per_query"`
	MajorityThreshold     float64 `json:"majority_threshold"`
	MaxParallel           *int    `json:"max_parallel,omitempty"`
	Status                string  `json:"status"`
	CurrentIteration      int     `json:"current_iteration"`
	BestDescriptionID     *int64  `json:"best_description_id,omitempty"`
	BestIteration         *int    `json:"best_iteration,omitempty"`
	Error                 string  `json:"error,omitempty"`
	StopReason            string  `json:"stop_reason,omitempty"`
	CreatedAt             string  `json:"created_at"`
	UpdatedAt             string  `json:"updated_at"`
	CompletedAt           string  `json:"completed_at,omitempty"`
}
