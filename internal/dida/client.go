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

	var reply []*api.Project
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

func (c *Client) CreateTask(ctx context.Context, req *api.CreateTaskRequest) (*api.Task, error) {
	c.log.Debug("CreateTask", zap.String("projectID", req.ProjectId))
	path := "/open/v1/task"

	var reply api.Task
	res, err := c.cli.
		R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&reply).
		Post(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create task")
	}
	if res.IsError() {
		return nil, errors.Errorf("failed to create task, status: %d, body: %s", res.StatusCode(), res.String())
	}
	return &reply, nil
}

func (c *Client) UpdateTask(ctx context.Context, req *api.UpdateTaskRequest) (*api.Task, error) {
	c.log.Debug("UpdateTask", zap.String("taskID", req.TaskId))
	path := "/open/v1/task/{task_id}"

	var reply api.Task
	res, err := c.cli.
		R().
		SetContext(ctx).
		SetPathParam("task_id", req.TaskId).
		SetBody(req).
		SetResult(&reply).
		Post(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update task")
	}
	if res.IsError() {
		return nil, errors.Errorf("failed to update task, status: %d, body: %s", res.StatusCode(), res.String())
	}
	return &reply, nil
}

func (c *Client) CompleteTask(ctx context.Context, projectID, taskID string) error {
	c.log.Debug("CompleteTask", zap.String("projectID", projectID), zap.String("taskID", taskID))
	path := "/open/v1/project/{project_id}/task/{task_id}/complete"

	res, err := c.cli.
		R().
		SetContext(ctx).
		SetPathParams(map[string]string{
			"project_id": projectID,
			"task_id":    taskID,
		}).
		Post(path)
	if err != nil {
		return errors.Wrap(err, "failed to complete task")
	}
	if res.IsError() {
		return errors.Errorf("failed to complete task, status: %d, body: %s", res.StatusCode(), res.String())
	}
	return nil
}

func (c *Client) DeleteTask(ctx context.Context, projectID, taskID string) error {
	c.log.Debug("DeleteTask", zap.String("projectID", projectID), zap.String("taskID", taskID))
	path := "/open/v1/project/{project_id}/task/{task_id}"

	res, err := c.cli.
		R().
		SetContext(ctx).
		SetPathParams(map[string]string{
			"project_id": projectID,
			"task_id":    taskID,
		}).
		Delete(path)
	if err != nil {
		return errors.Wrap(err, "failed to delete task")
	}
	if res.IsError() {
		return errors.Errorf("failed to delete task, status: %d, body: %s", res.StatusCode(), res.String())
	}
	return nil
}

func (c *Client) GetProject(ctx context.Context, projectID string) (*api.Project, error) {
	c.log.Debug("GetProject", zap.String("projectID", projectID))
	path := "/open/v1/project/{project_id}"

	var reply api.Project
	res, err := c.cli.
		R().
		SetContext(ctx).
		SetPathParam("project_id", projectID).
		SetResult(&reply).
		Get(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get project")
	}
	if res.IsError() {
		return nil, errors.Errorf("failed to get project, status: %d, body: %s", res.StatusCode(), res.String())
	}
	return &reply, nil
}

func (c *Client) GetProjectData(ctx context.Context, projectID string) (*api.ProjectData, error) {
	c.log.Debug("GetProjectData", zap.String("projectID", projectID))
	path := "/open/v1/project/{project_id}/data"

	var reply api.ProjectData
	res, err := c.cli.
		R().
		SetContext(ctx).
		SetPathParam("project_id", projectID).
		SetResult(&reply).
		Get(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get project data")
	}
	if res.IsError() {
		return nil, errors.Errorf("failed to get project data, status: %d, body: %s", res.StatusCode(), res.String())
	}
	return &reply, nil
}

func (c *Client) CreateProject(ctx context.Context, req *api.CreateProjectRequest) (*api.Project, error) {
	c.log.Debug("CreateProject", zap.String("name", req.Name))
	path := "/open/v1/project"

	var reply api.Project
	res, err := c.cli.
		R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&reply).
		Post(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create project")
	}
	if res.IsError() {
		return nil, errors.Errorf("failed to create project, status: %d, body: %s", res.StatusCode(), res.String())
	}
	return &reply, nil
}

func (c *Client) UpdateProject(ctx context.Context, req *api.UpdateProjectRequest) (*api.Project, error) {
	c.log.Debug("UpdateProject", zap.String("projectID", req.ProjectId))
	path := "/open/v1/project/{project_id}"

	var reply api.Project
	res, err := c.cli.
		R().
		SetContext(ctx).
		SetPathParam("project_id", req.ProjectId).
		SetBody(req).
		SetResult(&reply).
		Post(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update project")
	}
	if res.IsError() {
		return nil, errors.Errorf("failed to update project, status: %d, body: %s", res.StatusCode(), res.String())
	}
	return &reply, nil
}

func (c *Client) DeleteProject(ctx context.Context, projectID string) error {
	c.log.Debug("DeleteProject", zap.String("projectID", projectID))
	path := "/open/v1/project/{project_id}"

	res, err := c.cli.
		R().
		SetContext(ctx).
		SetPathParam("project_id", projectID).
		Delete(path)
	if err != nil {
		return errors.Wrap(err, "failed to delete project")
	}
	if res.IsError() {
		return errors.Errorf("failed to delete project, status: %d, body: %s", res.StatusCode(), res.String())
	}
	return nil
}
