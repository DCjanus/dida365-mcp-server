package main

import (
	"context"
	"encoding/json"

	"github.com/bufbuild/protovalidate-go"
	"github.com/cockroachdb/errors"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/dcjanus/dida365-mcp-server/gen/api"
	"github.com/dcjanus/dida365-mcp-server/internal/dida"
)

type DidaWrapper struct {
	log *zap.Logger
	cli *dida.Client
	ctx context.Context
}

func NewDidaWrapper(ctx context.Context, log *zap.Logger, token string) (*DidaWrapper, error) {
	cli := dida.NewClient(log, token)
	if _, err := cli.ListProjects(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to check dida token")
	}

	return &DidaWrapper{
		log: log.With(zap.String("component", "mcp.DidaWrapper")),
		cli: cli,
		ctx: ctx,
	}, nil
}

func (t *DidaWrapper) handleJSONResponse(data any) (*mcp.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal response")
	}
	return mcp.NewToolResultText(string(jsonBytes)), nil
}

func (t *DidaWrapper) parseJSONRequest(request mcp.CallToolRequest, target proto.Message) error {
	if len(request.Params.Arguments) == 0 {
		return errors.New("invalid request: missing parameters")
	}
	jsonBytes, err := json.Marshal(request.Params.Arguments)
	if err != nil {
		return errors.Wrap(err, "failed to marshal request arguments")
	}
	if err := protojson.Unmarshal(jsonBytes, target); err != nil {
		return errors.Wrap(err, "failed to parse request parameters")
	}
	if err := protovalidate.Validate(target); err != nil {
		return errors.Wrap(err, "failed to validate request parameters")
	}
	return nil
}

func (t *DidaWrapper) Tools() []server.ServerTool {
	return []server.ServerTool{
		t.ListProjects(t.ctx),
		t.GetProject(t.ctx),
		t.GetProjectData(t.ctx),
		t.CreateProject(t.ctx),
		t.UpdateProject(t.ctx),
		t.DeleteProject(t.ctx),
		t.GetTask(t.ctx),
		t.CreateTask(t.ctx),
		t.UpdateTask(t.ctx),
		t.CompleteTask(t.ctx),
		t.DeleteTask(t.ctx),
	}
}

func (t *DidaWrapper) ListProjects(ctx context.Context) server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("list_projects",
			mcp.WithDescription("List all projects, projects are the top level container for tasks"),
			mcp.WithString("random_string",
				mcp.Description("Dummy parameter for no-parameter tools"),
				mcp.Required(),
			),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			projects, err := t.cli.ListProjects(ctx)
			if err != nil {
				return nil, errors.Wrap(err, "failed to list projects")
			}
			return t.handleJSONResponse(projects)
		},
	}
}

func (t *DidaWrapper) GetProject(ctx context.Context) server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("get_project",
			mcp.WithDescription("Get a project by ID"),
			mcp.WithString("projectId",
				mcp.Description("ID of the project to get"),
				mcp.Required(),
			),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var req api.GetProjectRequest
			if err := t.parseJSONRequest(request, &req); err != nil {
				return nil, err
			}
			project, err := t.cli.GetProject(ctx, req.ProjectId)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get project")
			}
			return t.handleJSONResponse(project)
		},
	}
}

func (t *DidaWrapper) GetProjectData(ctx context.Context) server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("get_project_data",
			mcp.WithDescription("Get project data"),
			mcp.WithString("projectId",
				mcp.Description("ID of the project to get data for"),
				mcp.Required(),
			),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var req api.GetProjectDataRequest
			if err := t.parseJSONRequest(request, &req); err != nil {
				return nil, err
			}
			data, err := t.cli.GetProjectData(ctx, req.ProjectId)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get project data")
			}
			return t.handleJSONResponse(data)
		},
	}
}

func (t *DidaWrapper) CreateProject(ctx context.Context) server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("create_project",
			mcp.WithDescription("Create a new project"),
			mcp.WithString("name",
				mcp.Description("Name of the project"),
				mcp.Required(),
			),
			mcp.WithString("color",
				mcp.Description("Color of the project in hex format (e.g. \"#F18181\")"),
			),
			mcp.WithNumber("sortOrder",
				mcp.Description("Sort order of the project"),
			),
			mcp.WithString("viewMode",
				mcp.Description("View mode of the project, must be one of: \"list\", \"kanban\", \"timeline\""),
			),
			mcp.WithString("kind",
				mcp.Description("Kind of the project, must be one of: \"TASK\", \"NOTE\""),
			),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var req api.CreateProjectRequest
			if err := t.parseJSONRequest(request, &req); err != nil {
				return nil, err
			}
			createdProject, err := t.cli.CreateProject(ctx, &req)
			if err != nil {
				return nil, errors.Wrap(err, "failed to create project")
			}
			return t.handleJSONResponse(createdProject)
		},
	}
}

func (t *DidaWrapper) UpdateProject(ctx context.Context) server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("update_project",
			mcp.WithDescription("Update an existing project"),
			mcp.WithString("projectId",
				mcp.Description("ID of the project to update"),
				mcp.Required(),
			),
			mcp.WithString("name",
				mcp.Description("New name of the project"),
			),
			mcp.WithString("color",
				mcp.Description("New color of the project"),
			),
			mcp.WithNumber("sortOrder",
				mcp.Description("New sort order of the project"),
			),
			mcp.WithString("viewMode",
				mcp.Description("New view mode of the project (list, kanban, timeline)"),
			),
			mcp.WithString("kind",
				mcp.Description("New kind of the project (TASK, NOTE)"),
			),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var req api.UpdateProjectRequest
			if err := t.parseJSONRequest(request, &req); err != nil {
				return nil, err
			}
			updatedProject, err := t.cli.UpdateProject(ctx, &req)
			if err != nil {
				return nil, errors.Wrap(err, "failed to update project")
			}
			return t.handleJSONResponse(updatedProject)
		},
	}
}

func (t *DidaWrapper) DeleteProject(ctx context.Context) server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("delete_project",
			mcp.WithDescription("Delete a project"),
			mcp.WithString("projectId",
				mcp.Description("ID of the project to delete"),
				mcp.Required(),
			),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var req api.DeleteProjectRequest
			if err := t.parseJSONRequest(request, &req); err != nil {
				return nil, err
			}
			err := t.cli.DeleteProject(ctx, req.ProjectId)
			if err != nil {
				return nil, errors.Wrap(err, "failed to delete project")
			}
			return mcp.NewToolResultText("Project deleted successfully"), nil
		},
	}
}

func (t *DidaWrapper) GetTask(ctx context.Context) server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("get_task",
			mcp.WithDescription("Get a task by ID"),
			mcp.WithString("projectId",
				mcp.Description("ID of the project containing the task"),
				mcp.Required(),
			),
			mcp.WithString("taskId",
				mcp.Description("ID of the task to get"),
				mcp.Required(),
			),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var req api.GetTaskRequest
			if err := t.parseJSONRequest(request, &req); err != nil {
				return nil, err
			}
			task, err := t.cli.GetTask(ctx, req.ProjectId, req.TaskId)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get task")
			}
			return t.handleJSONResponse(task)
		},
	}
}

func (t *DidaWrapper) CreateTask(ctx context.Context) server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("create_task",
			mcp.WithDescription("Create a new task"),
			mcp.WithString("projectId",
				mcp.Description("ID of the project to create the task in"),
				mcp.Required(),
			),
			mcp.WithString("title",
				mcp.Description("Title of the task"),
				mcp.Required(),
			),
			mcp.WithString("content",
				mcp.Description("Content of the task"),
			),
			mcp.WithString("desc",
				mcp.Description("Description of the task"),
			),
			mcp.WithBoolean("isAllDay",
				mcp.Description("Whether the task is an all-day task"),
			),
			mcp.WithString("startDate",
				mcp.Description("Start date of the task in format \"yyyy-MM-dd'T'HH:mm:ssZ\" (e.g. \"2019-11-13T03:00:00+0000\")"),
			),
			mcp.WithString("dueDate",
				mcp.Description("Due date of the task in format \"yyyy-MM-dd'T'HH:mm:ssZ\" (e.g. \"2019-11-13T03:00:00+0000\")"),
			),
			mcp.WithString("timeZone",
				mcp.Description("Time zone of the task (e.g. \"America/Los_Angeles\")"),
			),
			mcp.WithArray("reminders",
				mcp.Description("Reminder times for the task"),
				mcp.Items(map[string]interface{}{
					"type":        "string",
					"description": "Reminder time in RRULE format (e.g. \"TRIGGER:P0DT9H0M0S\")",
				}),
			),
			mcp.WithString("repeatFlag",
				mcp.Description("Repeat flag for the task in RRULE format (e.g. \"RRULE:FREQ=DAILY;INTERVAL=1\")"),
			),
			mcp.WithNumber("priority",
				mcp.Description("Priority of the task (0: none, 1: low, 3: medium, 5: high)"),
			),
			mcp.WithNumber("sortOrder",
				mcp.Description("Sort order of the task"),
			),
			mcp.WithArray("items",
				mcp.Description("Checklist items of the task"),
				mcp.Items(map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"id": map[string]interface{}{
							"type":        "string",
							"description": "Unique identifier of the checklist item",
						},
						"status": map[string]interface{}{
							"type":        "number",
							"description": "Status of the checklist item (0: normal, 1: completed)",
							"enum":        []float64{0, 1},
						},
						"title": map[string]interface{}{
							"type":        "string",
							"description": "Title of the checklist item",
							"required":    true,
						},
						"sortOrder": map[string]interface{}{
							"type":        "number",
							"description": "Sort order of the checklist item",
						},
						"startDate": map[string]interface{}{
							"type":        "string",
							"description": "Start date of the checklist item in format \"yyyy-MM-dd'T'HH:mm:ssZ\"",
						},
						"isAllDay": map[string]interface{}{
							"type":        "boolean",
							"description": "Whether the checklist item is an all-day item",
						},
						"timeZone": map[string]interface{}{
							"type":        "string",
							"description": "Time zone of the checklist item (e.g. \"America/Los_Angeles\")",
						},
						"completedTime": map[string]interface{}{
							"type":        "string",
							"description": "Completion time of the checklist item in format \"yyyy-MM-dd'T'HH:mm:ssZ\"",
						},
					},
				}),
			),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var req api.CreateTaskRequest
			if err := t.parseJSONRequest(request, &req); err != nil {
				return nil, err
			}
			createdTask, err := t.cli.CreateTask(ctx, &req)
			if err != nil {
				return nil, errors.Wrap(err, "failed to create task")
			}
			return t.handleJSONResponse(createdTask)
		},
	}
}

func (t *DidaWrapper) UpdateTask(ctx context.Context) server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("update_task",
			mcp.WithDescription("Update an existing task"),
			mcp.WithString("taskId",
				mcp.Description("ID of the task to update"),
				mcp.Required(),
			),
			mcp.WithString("projectId",
				mcp.Description("ID of the project containing the task"),
				mcp.Required(),
			),
			mcp.WithString("title",
				mcp.Description("New title of the task"),
			),
			mcp.WithString("content",
				mcp.Description("New content of the task"),
			),
			mcp.WithString("desc",
				mcp.Description("New description of the task"),
			),
			mcp.WithBoolean("isAllDay",
				mcp.Description("Whether the task is an all-day task"),
			),
			mcp.WithString("startDate",
				mcp.Description("New start date of the task in format \"yyyy-MM-dd'T'HH:mm:ssZ\" (e.g. \"2019-11-13T03:00:00+0000\")"),
			),
			mcp.WithString("dueDate",
				mcp.Description("New due date of the task in format \"yyyy-MM-dd'T'HH:mm:ssZ\" (e.g. \"2019-11-13T03:00:00+0000\")"),
			),
			mcp.WithString("timeZone",
				mcp.Description("New time zone of the task (e.g. \"America/Los_Angeles\")"),
			),
			mcp.WithArray("reminders",
				mcp.Description("New reminder times for the task"),
				mcp.Items(map[string]interface{}{
					"type":        "string",
					"description": "Reminder time in RRULE format (e.g. \"TRIGGER:P0DT9H0M0S\")",
				}),
			),
			mcp.WithString("repeatFlag",
				mcp.Description("New repeat flag for the task in RRULE format (e.g. \"RRULE:FREQ=DAILY;INTERVAL=1\")"),
			),
			mcp.WithNumber("priority",
				mcp.Description("New priority of the task (0: none, 1: low, 3: medium, 5: high)"),
			),
			mcp.WithNumber("sortOrder",
				mcp.Description("New sort order of the task"),
			),
			mcp.WithArray("items",
				mcp.Description("New checklist items of the task"),
				mcp.Items(map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"id": map[string]interface{}{
							"type":        "string",
							"description": "Unique identifier of the checklist item",
						},
						"status": map[string]interface{}{
							"type":        "number",
							"description": "Status of the checklist item (0: normal, 1: completed)",
							"enum":        []float64{0, 1},
						},
						"title": map[string]interface{}{
							"type":        "string",
							"description": "Title of the checklist item",
							"required":    true,
						},
						"sortOrder": map[string]interface{}{
							"type":        "number",
							"description": "Sort order of the checklist item",
						},
						"startDate": map[string]interface{}{
							"type":        "string",
							"description": "Start date of the checklist item in format \"yyyy-MM-dd'T'HH:mm:ssZ\"",
						},
						"isAllDay": map[string]interface{}{
							"type":        "boolean",
							"description": "Whether the checklist item is an all-day item",
						},
						"timeZone": map[string]interface{}{
							"type":        "string",
							"description": "Time zone of the checklist item (e.g. \"America/Los_Angeles\")",
						},
						"completedTime": map[string]interface{}{
							"type":        "string",
							"description": "Completion time of the checklist item in format \"yyyy-MM-dd'T'HH:mm:ssZ\"",
						},
					},
				}),
			),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var req api.UpdateTaskRequest
			if err := t.parseJSONRequest(request, &req); err != nil {
				return nil, err
			}
			updatedTask, err := t.cli.UpdateTask(ctx, &req)
			if err != nil {
				return nil, errors.Wrap(err, "failed to update task")
			}
			return t.handleJSONResponse(updatedTask)
		},
	}
}

func (t *DidaWrapper) CompleteTask(ctx context.Context) server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("complete_task",
			mcp.WithDescription("Mark a task as completed"),
			mcp.WithString("projectId",
				mcp.Description("ID of the project containing the task"),
				mcp.Required(),
			),
			mcp.WithString("taskId",
				mcp.Description("ID of the task to complete"),
				mcp.Required(),
			),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var req api.CompleteTaskRequest
			if err := t.parseJSONRequest(request, &req); err != nil {
				return nil, err
			}
			err := t.cli.CompleteTask(ctx, req.ProjectId, req.TaskId)
			if err != nil {
				return nil, errors.Wrap(err, "failed to complete task")
			}
			return mcp.NewToolResultText("Task marked as completed"), nil
		},
	}
}

func (t *DidaWrapper) DeleteTask(ctx context.Context) server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("delete_task",
			mcp.WithDescription("Delete a task"),
			mcp.WithString("projectId",
				mcp.Description("ID of the project containing the task"),
				mcp.Required(),
			),
			mcp.WithString("taskId",
				mcp.Description("ID of the task to delete"),
				mcp.Required(),
			),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var req api.DeleteTaskRequest
			if err := t.parseJSONRequest(request, &req); err != nil {
				return nil, err
			}
			err := t.cli.DeleteTask(ctx, req.ProjectId, req.TaskId)
			if err != nil {
				return nil, errors.Wrap(err, "failed to delete task")
			}
			return mcp.NewToolResultText("Task deleted successfully"), nil
		},
	}
}
