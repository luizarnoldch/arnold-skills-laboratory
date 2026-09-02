package labapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

func (c *Client) ListSkills(ctx context.Context) ([]Skill, error) {
	var out []Skill
	if err := c.doJSON(ctx, http.MethodGet, "/api/v1/skills", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) GetSkill(ctx context.Context, id int64) (SkillDetail, error) {
	var out SkillDetail
	path := fmt.Sprintf("/api/v1/skills/%d", id)
	if err := c.doJSON(ctx, http.MethodGet, path, nil, &out); err != nil {
		return SkillDetail{}, err
	}
	return out, nil
}

func (c *Client) ListDescriptions(ctx context.Context, skillID int64) ([]DescriptionView, error) {
	var out []DescriptionView
	path := fmt.Sprintf("/api/v1/skills/%d/descriptions", skillID)
	if err := c.doJSON(ctx, http.MethodGet, path, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) ListContents(ctx context.Context, skillID int64) ([]ContentView, error) {
	var out []ContentView
	path := fmt.Sprintf("/api/v1/skills/%d/contents", skillID)
	if err := c.doJSON(ctx, http.MethodGet, path, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) ListTriggerQueries(ctx context.Context, skillID int64, split string) ([]TriggerQuery, error) {
	path := fmt.Sprintf("/api/v1/skills/%d/trigger-queries", skillID)
	if split != "" {
		path += "?split=" + url.QueryEscape(split)
	}
	var out []TriggerQuery
	if err := c.doJSON(ctx, http.MethodGet, path, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) BulkReplaceTriggerQueries(ctx context.Context, skillID int64, items []BulkItem) ([]TriggerQuery, error) {
	var out []TriggerQuery
	path := fmt.Sprintf("/api/v1/skills/%d/trigger-queries", skillID)
	if err := c.doJSON(ctx, http.MethodPut, path, items, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) SplitTriggerQueries(ctx context.Context, skillID int64, req SplitRequest) (SplitResponse, error) {
	var out SplitResponse
	path := fmt.Sprintf("/api/v1/skills/%d/trigger-queries/split", skillID)
	if err := c.doJSON(ctx, http.MethodPost, path, req, &out); err != nil {
		return SplitResponse{}, err
	}
	return out, nil
}

func (c *Client) PatchTriggerQuery(ctx context.Context, id int64, req PatchTriggerQueryRequest) (TriggerQuery, error) {
	var out TriggerQuery
	path := fmt.Sprintf("/api/v1/trigger-queries/%d", id)
	if err := c.doJSON(ctx, http.MethodPatch, path, req, &out); err != nil {
		return TriggerQuery{}, err
	}
	return out, nil
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
		return fmt.Errorf("no se pudo conectar con laboratory-api (%s): %w", c.base, err)
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
		return fmt.Errorf("laboratory-api %s %s: %s", method, path, msg)
	}

	if dst == nil || res.StatusCode == http.StatusNoContent || len(raw) == 0 {
		return nil
	}
	if err := json.Unmarshal(raw, dst); err != nil {
		return fmt.Errorf("decode laboratory-api %s: %w", path, err)
	}
	return nil
}
