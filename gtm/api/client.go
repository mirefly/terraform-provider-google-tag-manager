package api

import (
	"context"
	"errors"
	"time"

	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/tagmanager/v2"
)

type ClientOptions struct {
	CredentialFile             string
	AccountId                  string
	ContainerId                string
	WaitingTimeBeforeEachQuery time.Duration
}

type Client struct {
	*tagmanager.Service

	Options *ClientOptions
}

func NewClient(opts *ClientOptions) (*Client, error) {
	var ctx = context.Background()

	srv, err := tagmanager.NewService(ctx, option.WithCredentialsFile(opts.CredentialFile))
	if err != nil {
		return nil, err
	}

	return &Client{Service: srv, Options: opts}, nil
}

func (c *Client) containerPath() string {
	opts := c.Options
	return "accounts/" + opts.AccountId + "/containers/" + opts.ContainerId
}

func (c *Client) beforeEachQuery() {
	time.Sleep(c.Options.WaitingTimeBeforeEachQuery)
}

var ErrNotExist = errors.New("not exist")

func (c *Client) CreateWorkspace(ws *tagmanager.Workspace) (*tagmanager.Workspace, error) {
	c.beforeEachQuery()
	return c.Accounts.Containers.Workspaces.Create(c.containerPath(), ws).Do()
}

func (c *Client) ListWorkspaces() ([]*tagmanager.Workspace, error) {
	c.beforeEachQuery()
	resp, err := c.Accounts.Containers.Workspaces.List(c.containerPath()).Do()
	if err != nil {
		return nil, err
	} else {
		return resp.Workspace, nil
	}
}

func (c *Client) Workspace(id string) (*tagmanager.Workspace, error) {
	c.beforeEachQuery()
	ws, err := c.Accounts.Containers.Workspaces.Get(c.containerPath() + "/workspaces/" + id).Do()

	if errTyped, ok := err.(*googleapi.Error); ok && errTyped.Code == 404 {
		return nil, ErrNotExist
	} else {
		return ws, err
	}
}

func (c *Client) UpdateWorkspaces(id string, ws *tagmanager.Workspace) (*tagmanager.Workspace, error) {
	c.beforeEachQuery()
	return c.Accounts.Containers.Workspaces.Update(c.containerPath()+"/workspaces/"+id, ws).Do()
}

func (c *Client) DeleteWorkspace(id string) error {
	c.beforeEachQuery()
	return c.Accounts.Containers.Workspaces.Delete(c.containerPath() + "/workspaces/" + id).Do()
}

func (c *Client) workspacePath(id string) string {
	return c.containerPath() + "/workspaces/" + id
}

func (c *Client) CreateTag(workspaceId string, tag *tagmanager.Tag) (*tagmanager.Tag, error) {
	c.beforeEachQuery()
	return c.Accounts.Containers.Workspaces.Tags.Create(c.workspacePath(workspaceId), tag).Do()
}

func (c *Client) ListTags(workspaceId string) ([]*tagmanager.Tag, error) {
	c.beforeEachQuery()
	resp, err := c.Accounts.Containers.Workspaces.Tags.List(c.workspacePath(workspaceId)).Do()
	if err != nil {
		return nil, err
	} else {
		return resp.Tag, nil
	}
}

func (c *Client) Tag(workspaceId string, tagId string) (*tagmanager.Tag, error) {
	c.beforeEachQuery()
	tag, err := c.Accounts.Containers.Workspaces.Tags.Get(c.workspacePath(workspaceId) + "/tags/" + tagId).Do()

	if errTyped, ok := err.(*googleapi.Error); ok && errTyped.Code == 404 {
		return nil, ErrNotExist
	} else {
		return tag, err
	}
}

func (c *Client) UpdateTag(workspaceId string, tagId string, tag *tagmanager.Tag) (*tagmanager.Tag, error) {
	c.beforeEachQuery()
	return c.Accounts.Containers.Workspaces.Tags.Update(c.workspacePath(workspaceId)+"/tags/"+tagId, tag).Do()
}

func (c *Client) DeleteTag(workspaceId string, tagId string) error {
	c.beforeEachQuery()
	return c.Accounts.Containers.Workspaces.Tags.Delete(c.workspacePath(workspaceId) + "/tags/" + tagId).Do()
}

func (c *Client) CreateVariable(workspaceId string, variable *tagmanager.Variable) (*tagmanager.Variable, error) {
	c.beforeEachQuery()
	return c.Accounts.Containers.Workspaces.Variables.Create(c.workspacePath(workspaceId), variable).Do()
}

func (c *Client) ListVariables(workspaceId string) ([]*tagmanager.Variable, error) {
	c.beforeEachQuery()
	resp, err := c.Accounts.Containers.Workspaces.Variables.List(c.workspacePath(workspaceId)).Do()
	if err != nil {
		return nil, err
	} else {
		return resp.Variable, nil
	}
}

func (c *Client) Variable(workspaceId string, variableId string) (*tagmanager.Variable, error) {
	c.beforeEachQuery()
	variable, err := c.Accounts.Containers.Workspaces.Variables.Get(c.workspacePath(workspaceId) + "/variables/" + variableId).Do()

	if errTyped, ok := err.(*googleapi.Error); ok && errTyped.Code == 404 {
		return nil, ErrNotExist
	} else {
		return variable, err
	}
}

func (c *Client) UpdateVariable(workspaceId string, variableId string, variable *tagmanager.Variable) (*tagmanager.Variable, error) {
	c.beforeEachQuery()
	return c.Accounts.Containers.Workspaces.Variables.Update(c.workspacePath(workspaceId)+"/variables/"+variableId, variable).Do()
}

func (c *Client) DeleteVariable(workspaceId string, variableId string) error {
	c.beforeEachQuery()
	return c.Accounts.Containers.Workspaces.Variables.Delete(c.workspacePath(workspaceId) + "/variables/" + variableId).Do()
}

func (c *Client) CreateTrigger(workspaceId string, trigger *tagmanager.Trigger) (*tagmanager.Trigger, error) {
	c.beforeEachQuery()
	return c.Accounts.Containers.Workspaces.Triggers.Create(c.workspacePath(workspaceId), trigger).Do()
}

func (c *Client) ListTriggers(workspaceId string) ([]*tagmanager.Trigger, error) {
	c.beforeEachQuery()
	resp, err := c.Accounts.Containers.Workspaces.Triggers.List(c.workspacePath(workspaceId)).Do()
	if err != nil {
		return nil, err
	} else {
		return resp.Trigger, nil
	}
}

func (c *Client) Trigger(workspaceId string, triggerId string) (*tagmanager.Trigger, error) {
	c.beforeEachQuery()
	trigger, err := c.Accounts.Containers.Workspaces.Triggers.Get(c.workspacePath(workspaceId) + "/triggers/" + triggerId).Do()

	if errTyped, ok := err.(*googleapi.Error); ok && errTyped.Code == 404 {
		return nil, ErrNotExist
	} else {
		return trigger, err
	}
}

func (c *Client) UpdateTrigger(workspaceId string, triggerId string, trigger *tagmanager.Trigger) (*tagmanager.Trigger, error) {
	c.beforeEachQuery()
	return c.Accounts.Containers.Workspaces.Triggers.Update(c.workspacePath(workspaceId)+"/triggers/"+triggerId, trigger).Do()
}

func (c *Client) DeleteTrigger(workspaceId string, triggerId string) error {
	c.beforeEachQuery()
	return c.Accounts.Containers.Workspaces.Triggers.Delete(c.workspacePath(workspaceId) + "/triggers/" + triggerId).Do()
}
