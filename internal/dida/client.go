package dida

import (
	"context"

	"github.com/cockroachdb/errors"
	"go.uber.org/zap"
	"resty.dev/v3"

	"github.com/dcjanus/dida365-mcp-server/gen/api"
)

type Client struct {
	log *zap.Logger
	cli *resty.Client
}

func NewClient(log *zap.Logger, token string) *Client {
	log = log.With(zap.String("component", "dida.Client"))

	return &Client{
		log: log,
		cli: resty.New().SetBaseURL("https://api.dida365.com").SetAuthToken(token),
	}
}

func (c *Client) GetTask(ctx context.Context, projectID, taskID string) (*api.Task, error) {
	c.log.Debug("GetTask", zap.String("projectID", projectID), zap.String("taskID", taskID))
	path := "/open/v1/project/{project_id}/task/{task_id}"

	var reply api.Task

	res, err := c.cli.
		R().
		SetContext(ctx).
		SetPathParams(map[string]string{
			"project_id": projectID,
			"task_id":    taskID,
		}).
		SetResult(&reply).
		Get(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get task")
	}
	if res.IsError() {
		return nil, errors.Errorf("failed to get task, status: %d, body: %s", res.StatusCode(), res.String())
	}

	return &reply, nil
}

func (c *Client) ListProjects(ctx context.Context) ([]*api.Project, error) {
	c.log.Debug("ListProjects")
	path := "/open/v1/project"

	reply := []*api.Project{}
	res, err := c.cli.
		R().
		SetContext(ctx).
		SetResult(&reply).
		Get(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list projects")
	}
	if res.IsError() {
		return nil, errors.Errorf("failed to list projects, status: %d, body: %s", res.StatusCode(), res.String())
	}
	return reply, nil
}
