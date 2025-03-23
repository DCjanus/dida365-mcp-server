// Package dida provides a client for interacting with the Dida365 API.
// It implements all the basic operations for managing projects and tasks.
package dida

import (
	"context"

	"github.com/bufbuild/protovalidate-go"
	"github.com/cockroachdb/errors"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"resty.dev/v3"

	"github.com/dcjanus/dida365-mcp-server/gen/api"
)

// Client represents a Dida365 API client.
type Client struct {
	log *zap.Logger
	cli *resty.Client
}

// NewClient creates a new Dida365 API client with the given logger and authentication token.
func NewClient(log *zap.Logger, token string) *Client {
	return &Client{
		log: log.With(zap.String("component", "dida.Client")),
		cli: resty.
			New().
			SetBaseURL("https://api.dida365.com").
			SetAuthToken(token),
	}
}

// -----------------------------------------------------------------------------
// Common request handling
// -----------------------------------------------------------------------------

// request represents a generic API request.
type request struct {
	method string
	path   string
	params map[string]string
	body   proto.Message
	result proto.Message
}

// doRequest executes an API request and handles common error cases.
func (c *Client) doRequest(ctx context.Context, req request) error {
	r := c.cli.R().SetContext(ctx)

	if req.params != nil {
		r.SetPathParams(req.params)
	}
	if req.body != nil {
		if err := protovalidate.Validate(req.body); err != nil {
			return errors.Wrap(err, "invalid request body")
		}

		body, err := protojson.MarshalOptions{UseEnumNumbers: true}.Marshal(req.body)
		if err != nil {
			return errors.Wrap(err, "failed to marshal request body")
		}
		r.SetContentType("application/json")
		r.SetBody(string(body))
	}
	if req.result != nil {
		r.SetExpectResponseContentType("application/json")
	}

	res, err := r.Execute(req.method, req.path)
	if err != nil {
		return errors.Wrap(err, "failed to execute request")
	}
	if res.IsError() {
		return errors.Errorf("request failed with status %d: %s", res.StatusCode(), res.String())
	}

	if req.result != nil {
		if err := (protojson.UnmarshalOptions{DiscardUnknown: true}).Unmarshal(res.Bytes(), req.result); err != nil {
			return errors.Wrap(err, "failed to unmarshal response")
		}
	}

	return nil
}

// -----------------------------------------------------------------------------
// Task operations
// -----------------------------------------------------------------------------

// GetTask retrieves a specific task by project ID and task ID.
func (c *Client) GetTask(ctx context.Context, projectID, taskID string) (*api.Task, error) {
	c.log.Debug("GetTask", zap.String("projectID", projectID), zap.String("taskID", taskID))

	var result api.Task
	err := c.doRequest(ctx, request{
		method: "GET",
		path:   "/open/v1/project/{project_id}/task/{task_id}",
		params: map[string]string{
			"project_id": projectID,
			"task_id":    taskID,
		},
		result: &result,
	})
	if err != nil {
		return nil, errors.Wrap(err, "get task")
	}

	return &result, nil
}

// CreateTask creates a new task with the given parameters.
func (c *Client) CreateTask(ctx context.Context, req *api.CreateTaskRequest) (*api.Task, error) {
	c.log.Debug("CreateTask", zap.String("projectID", req.ProjectId))

	var result api.Task
	err := c.doRequest(ctx, request{
		method: "POST",
		path:   "/open/v1/task",
		body:   req,
		result: &result,
	})
	if err != nil {
		return nil, errors.Wrap(err, "create task")
	}

	return &result, nil
}

// UpdateTask updates an existing task with the given parameters.
func (c *Client) UpdateTask(ctx context.Context, req *api.UpdateTaskRequest) (*api.Task, error) {
	c.log.Debug("UpdateTask", zap.String("taskID", req.TaskId))

	var result api.Task
	err := c.doRequest(ctx, request{
		method: "POST",
		path:   "/open/v1/task/{task_id}",
		params: map[string]string{
			"task_id": req.TaskId,
		},
		body:   req,
		result: &result,
	})
	if err != nil {
		return nil, errors.Wrap(err, "update task")
	}

	return &result, nil
}

// CompleteTask marks a task as completed.
func (c *Client) CompleteTask(ctx context.Context, projectID, taskID string) error {
	c.log.Debug("CompleteTask", zap.String("projectID", projectID), zap.String("taskID", taskID))

	err := c.doRequest(ctx, request{
		method: "POST",
		path:   "/open/v1/project/{project_id}/task/{task_id}/complete",
		params: map[string]string{
			"project_id": projectID,
			"task_id":    taskID,
		},
	})
	if err != nil {
		return errors.Wrap(err, "complete task")
	}

	return nil
}

// DeleteTask deletes a task by its ID.
func (c *Client) DeleteTask(ctx context.Context, projectID, taskID string) error {
	c.log.Debug("DeleteTask", zap.String("projectID", projectID), zap.String("taskID", taskID))

	err := c.doRequest(ctx, request{
		method: "DELETE",
		path:   "/open/v1/project/{project_id}/task/{task_id}",
		params: map[string]string{
			"project_id": projectID,
			"task_id":    taskID,
		},
	})
	if err != nil {
		return errors.Wrap(err, "delete task")
	}

	return nil
}

// -----------------------------------------------------------------------------
// Project operations
// -----------------------------------------------------------------------------

// ListProjects returns a list of all projects.
func (c *Client) ListProjects(ctx context.Context) ([]*api.Project, error) {
	c.log.Debug("ListProjects")

	var reply structpb.ListValue
	err := c.doRequest(ctx, request{
		method: "GET",
		path:   "/open/v1/project",
		result: &reply,
	})
	if err != nil {
		return nil, errors.Wrap(err, "list projects")
	}

	projects := make([]*api.Project, 0, len(reply.Values))
	for _, v := range reply.Values {
		buf, err := v.MarshalJSON()
		if err != nil {
			return nil, errors.Wrap(err, "marshal project")
		}
		var p api.Project
		if err := protojson.Unmarshal(buf, &p); err != nil {
			return nil, errors.Wrap(err, "unmarshal project")
		}
		projects = append(projects, &p)
	}

	return projects, nil
}

// GetProject retrieves a specific project by its ID.
func (c *Client) GetProject(ctx context.Context, projectID string) (*api.Project, error) {
	c.log.Debug("GetProject", zap.String("projectID", projectID))

	var result api.Project
	err := c.doRequest(ctx, request{
		method: "GET",
		path:   "/open/v1/project/{project_id}",
		params: map[string]string{
			"project_id": projectID,
		},
		result: &result,
	})
	if err != nil {
		return nil, errors.Wrap(err, "get project")
	}

	return &result, nil
}

// GetProjectData retrieves detailed data for a specific project.
func (c *Client) GetProjectData(ctx context.Context, projectID string) (*api.ProjectData, error) {
	c.log.Debug("GetProjectData", zap.String("projectID", projectID))

	var result api.ProjectData
	err := c.doRequest(ctx, request{
		method: "GET",
		path:   "/open/v1/project/{project_id}/data",
		params: map[string]string{
			"project_id": projectID,
		},
		result: &result,
	})
	if err != nil {
		return nil, errors.Wrap(err, "get project data")
	}

	return &result, nil
}

// CreateProject creates a new project with the given parameters.
func (c *Client) CreateProject(ctx context.Context, req *api.CreateProjectRequest) (*api.Project, error) {
	c.log.Debug("CreateProject", zap.String("name", req.Name))

	var result api.Project
	err := c.doRequest(ctx, request{
		method: "POST",
		path:   "/open/v1/project",
		body:   req,
		result: &result,
	})
	if err != nil {
		return nil, errors.Wrap(err, "create project")
	}

	return &result, nil
}

// UpdateProject updates an existing project with the given parameters.
func (c *Client) UpdateProject(ctx context.Context, req *api.UpdateProjectRequest) (*api.Project, error) {
	c.log.Debug("UpdateProject", zap.String("projectID", req.ProjectId))

	var result api.Project
	err := c.doRequest(ctx, request{
		method: "POST",
		path:   "/open/v1/project/{project_id}",
		params: map[string]string{
			"project_id": req.ProjectId,
		},
		body:   req,
		result: &result,
	})
	if err != nil {
		return nil, errors.Wrap(err, "update project")
	}

	return &result, nil
}

// DeleteProject deletes a project by its ID.
func (c *Client) DeleteProject(ctx context.Context, projectID string) error {
	c.log.Debug("DeleteProject", zap.String("projectID", projectID))

	err := c.doRequest(ctx, request{
		method: "DELETE",
		path:   "/open/v1/project/{project_id}",
		params: map[string]string{
			"project_id": projectID,
		},
	})
	if err != nil {
		return errors.Wrap(err, "delete project")
	}

	return nil
}
