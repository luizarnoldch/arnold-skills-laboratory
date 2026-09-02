package tools

import (
	"context"
	"fmt"

	"github.com/luizarnoldch/skills-lab-mcp/pkg/client/orch"
)

type JobGetArgs struct {
	JobType string `json:"job_type" jsonschema:"baseline_eval, trigger_eval u optimize"`
	JobID   string `json:"job_id" jsonschema:"UUID del job devuelto por *_start"`
}

type JobProgress struct {
	Done  int `json:"done"`
	Total int `json:"total"`
}

type JobGetResult struct {
	JobType string `json:"job_type"`
	JobID   string `json:"job_id"`
	Status  string `json:"status"`

	SkillID     int64  `json:"skill_id,omitempty"`
	SkillName   string `json:"skill_name,omitempty"`
	DescriptionID int64 `json:"description_id,omitempty"`
	ContentID   int64  `json:"content_id,omitempty"`
	Split       string `json:"split,omitempty"`

	Progress *JobProgress `json:"progress,omitempty"`

	TrainEvalRunID      *int64   `json:"train_eval_run_id,omitempty"`
	ValidationEvalRunID *int64   `json:"validation_eval_run_id,omitempty"`
	TrainAccuracy       *float64 `json:"train_accuracy,omitempty"`
	ValidationAccuracy  *float64 `json:"validation_accuracy,omitempty"`

	LabEvalRunID *int64   `json:"lab_eval_run_id,omitempty"`
	Accuracy     *float64 `json:"accuracy,omitempty"`
	Correct      *int64   `json:"correct,omitempty"`
	Total        *int64   `json:"total,omitempty"`

	StartingDescriptionID *int64  `json:"starting_description_id,omitempty"`
	LabOptimizeRunID      *int64  `json:"lab_optimize_run_id,omitempty"`
	CurrentIteration      int     `json:"current_iteration,omitempty"`
	MaxIters              int64   `json:"max_iters,omitempty"`
	Threshold             float64 `json:"threshold,omitempty"`
	RunsPerQuery          int64   `json:"runs_per_query,omitempty"`
	MajorityThreshold     float64 `json:"majority_threshold,omitempty"`
	BestDescriptionID     *int64  `json:"best_description_id,omitempty"`
	BestIteration         *int    `json:"best_iteration,omitempty"`
	StopReason            string  `json:"stop_reason,omitempty"`

	Error       string `json:"error,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
	CompletedAt string `json:"completed_at,omitempty"`
}

type Jobs struct {
	Orch *orch.Client
}

func (j *Jobs) Get(ctx context.Context, args JobGetArgs) (JobGetResult, error) {
	if args.JobType == "" || args.JobID == "" {
		return JobGetResult{}, fmt.Errorf("job_type y job_id son requeridos")
	}

	switch args.JobType {
	case "baseline_eval":
		job, err := j.Orch.GetBaselineEvalJob(ctx, args.JobID)
		if err != nil {
			return JobGetResult{}, err
		}
		return JobGetResult{
			JobType:               args.JobType,
			JobID:                 args.JobID,
			Status:                job.Status,
			SkillID:               job.SkillID,
			SkillName:             job.SkillName,
			DescriptionID:         job.DescriptionID,
			ContentID:             job.ContentID,
			Progress:              &JobProgress{Done: job.ProgressDone, Total: job.ProgressTotal},
			TrainEvalRunID:        job.TrainEvalRunID,
			ValidationEvalRunID:   job.ValidationEvalRunID,
			TrainAccuracy:         job.TrainAccuracy,
			ValidationAccuracy:    job.ValidationAccuracy,
			Error:                 job.Error,
			CreatedAt:             job.CreatedAt,
			UpdatedAt:             job.UpdatedAt,
			CompletedAt:           job.CompletedAt,
		}, nil

	case "trigger_eval":
		job, err := j.Orch.GetTriggerEvalBatch(ctx, args.JobID)
		if err != nil {
			return JobGetResult{}, err
		}
		return JobGetResult{
			JobType:       args.JobType,
			JobID:         args.JobID,
			Status:        job.Status,
			SkillID:       job.SkillID,
			SkillName:     job.SkillName,
			DescriptionID: job.DescriptionID,
			ContentID:     job.ContentID,
			Split:         job.Split,
			Progress:      &JobProgress{Done: job.ProgressDone, Total: job.ProgressTotal},
			LabEvalRunID:  job.LabEvalRunID,
			Accuracy:      job.Accuracy,
			Correct:       job.Correct,
			Total:         job.Total,
			Error:         job.Error,
			CreatedAt:     job.CreatedAt,
			UpdatedAt:     job.UpdatedAt,
			CompletedAt:   job.CompletedAt,
		}, nil

	case "optimize":
		job, err := j.Orch.GetOptimizeJob(ctx, args.JobID)
		if err != nil {
			return JobGetResult{}, err
		}
		return JobGetResult{
			JobType:               args.JobType,
			JobID:                 args.JobID,
			Status:                job.Status,
			SkillID:               job.SkillID,
			SkillName:             job.SkillName,
			StartingDescriptionID: &job.StartingDescriptionID,
			ContentID:             job.ContentID,
			LabOptimizeRunID:      job.LabOptimizeRunID,
			CurrentIteration:      job.CurrentIteration,
			MaxIters:              job.MaxIters,
			Threshold:             job.Threshold,
			RunsPerQuery:          job.RunsPerQuery,
			MajorityThreshold:     job.MajorityThreshold,
			BestDescriptionID:     job.BestDescriptionID,
			BestIteration:         job.BestIteration,
			StopReason:            job.StopReason,
			Error:                 job.Error,
			CreatedAt:             job.CreatedAt,
			UpdatedAt:             job.UpdatedAt,
			CompletedAt:           job.CompletedAt,
		}, nil

	default:
		return JobGetResult{}, fmt.Errorf("job_type inválido %q (use baseline_eval, trigger_eval u optimize)", args.JobType)
	}
}
