package orch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	base   string
	client *http.Client
}

func New(baseURL string, timeout time.Duration) *Client {
	return &Client{
		base: strings.TrimRight(baseURL, "/"),
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) CheckReady(ctx context.Context) bool {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.base+"/ready", nil)
	if err != nil {
		return false
	}
	res, err := c.client.Do(req)
	if err != nil {
		return false
	}
	defer res.Body.Close()
	return res.StatusCode == http.StatusOK
}

func (c *Client) CreateBaselineEvalJob(ctx context.Context, req CreateBaselineEvalJobRequest) (JobCreateResponse, error) {
	var out JobCreateResponse
	if err := c.doJSON(ctx, http.MethodPost, "/api/v1/baseline-eval-jobs", req, &out); err != nil {
		return JobCreateResponse{}, err
	}
	return out, nil
}

func (c *Client) GetBaselineEvalJob(ctx context.Context, id string) (BaselineEvalJob, error) {
	var out BaselineEvalJob
	path := "/api/v1/baseline-eval-jobs/" + urlPathEscape(id)
	if err := c.doJSON(ctx, http.MethodGet, path, nil, &out); err != nil {
		return BaselineEvalJob{}, err
	}
	return out, nil
}

func (c *Client) CreateTriggerEvalBatch(ctx context.Context, req CreateTriggerEvalBatchRequest) (JobCreateResponse, error) {
	var out JobCreateResponse
	if err := c.doJSON(ctx, http.MethodPost, "/api/v1/trigger-eval-batches", req, &out); err != nil {
		return JobCreateResponse{}, err
	}
	return out, nil
}

func (c *Client) GetTriggerEvalBatch(ctx context.Context, id string) (TriggerEvalBatch, error) {
	var out TriggerEvalBatch
	path := "/api/v1/trigger-eval-batches/" + urlPathEscape(id)
	if err := c.doJSON(ctx, http.MethodGet, path, nil, &out); err != nil {
		return TriggerEvalBatch{}, err
	}
	return out, nil
}

func (c *Client) CreateOptimizeJob(ctx context.Context, req CreateOptimizeJobRequest) (JobCreateResponse, error) {
	var out JobCreateResponse
	if err := c.doJSON(ctx, http.MethodPost, "/api/v1/optimize-jobs", req, &out); err != nil {
		return JobCreateResponse{}, err
	}
	return out, nil
}

func (c *Client) GetOptimizeJob(ctx context.Context, id string) (OptimizeJob, error) {
	var out OptimizeJob
	path := "/api/v1/optimize-jobs/" + urlPathEscape(id)
	if err := c.doJSON(ctx, http.MethodGet, path, nil, &out); err != nil {
		return OptimizeJob{}, err
	}
	return out, nil
}

func urlPathEscape(id string) string {
	return strings.ReplaceAll(id, "/", "%2F")
}

func (c *Client) doJSON(ctx context.Context, method, path string, body any, dst any) error {
	var rdr io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}
		rdr = bytes.NewReader(data)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.base+path, rdr)
	if err != nil {
		return err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	res, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("no se pudo conectar con orchestrator (%s): %w", c.base, err)
	}
	defer res.Body.Close()

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		msg := strings.TrimSpace(string(raw))
		var errBody struct {
			Error string `json:"error"`
		}
		if json.Unmarshal(raw, &errBody) == nil && errBody.Error != "" {
			msg = errBody.Error
		}
		if msg == "" {
			msg = res.Status
		}
		return fmt.Errorf("orchestrator %s %s: %s", method, path, msg)
	}

	if dst == nil || res.StatusCode == http.StatusNoContent || len(raw) == 0 {
		return nil
	}
	if err := json.Unmarshal(raw, dst); err != nil {
		return fmt.Errorf("decode orchestrator %s: %w", path, err)
	}
	return nil
}
