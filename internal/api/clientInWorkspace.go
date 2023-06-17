package api

import (
	"google.golang.org/api/tagmanager/v2"
)

type ClientInWorkspaceOptions struct {
	*ClientOptions
	WorkspaceName string
	WorkspaceId   string
}

type ClientInWorkspace struct {
	*Client

	Options *ClientInWorkspaceOptions
}

func NewClientInWorkspace(options *ClientInWorkspaceOptions) (*ClientInWorkspace, error) {
	client, err := NewClient(options.ClientOptions)
	if err != nil {
		return nil, err
	}

	workspaces, err := client.ListWorkspaces()
	if err != nil {
		return nil, err
	}

	for _, workspace := range workspaces {
		if workspace.Name == options.WorkspaceName {
			options.WorkspaceId = workspace.WorkspaceId

			return &ClientInWorkspace{
				Client:  client,
				Options: options,
			}, nil
		}
	}

	workspace, err := client.CreateWorkspace(&tagmanager.Workspace{Name: options.WorkspaceName})
	if err != nil {
		return nil, err
	} else {
		options.WorkspaceId = workspace.WorkspaceId
		return &ClientInWorkspace{
			Client:  client,
			Options: options,
		}, nil
	}
}

// Tag CRUD

func (c *ClientInWorkspace) CreateTag(tag *tagmanager.Tag) (*tagmanager.Tag, error) {
	return c.Client.CreateTag(c.Options.WorkspaceId, tag)
}

func (c *ClientInWorkspace) ListTags() ([]*tagmanager.Tag, error) {
	return c.Client.ListTags(c.Options.WorkspaceId)
}

func (c *ClientInWorkspace) Tag(tagId string) (*tagmanager.Tag, error) {
	return c.Client.Tag(c.Options.WorkspaceId, tagId)
}

func (c *ClientInWorkspace) UpdateTag(tagId string, tag *tagmanager.Tag) (*tagmanager.Tag, error) {
	return c.Client.UpdateTag(c.Options.WorkspaceId, tagId, tag)
}

func (c *ClientInWorkspace) DeleteTag(tagId string) error {
	return c.Client.DeleteTag(c.Options.WorkspaceId, tagId)
}

// Variable CRUD

func (c *ClientInWorkspace) CreateVariable(variable *tagmanager.Variable) (*tagmanager.Variable, error) {
	return c.Client.CreateVariable(c.Options.WorkspaceId, variable)
}

func (c *ClientInWorkspace) ListVariables() ([]*tagmanager.Variable, error) {
	return c.Client.ListVariables(c.Options.WorkspaceId)
}

func (c *ClientInWorkspace) Variable(variableId string) (*tagmanager.Variable, error) {
	return c.Client.Variable(c.Options.WorkspaceId, variableId)
}

func (c *ClientInWorkspace) UpdateVariable(variableId string, variable *tagmanager.Variable) (*tagmanager.Variable, error) {
	return c.Client.UpdateVariable(c.Options.WorkspaceId, variableId, variable)
}

func (c *ClientInWorkspace) DeleteVariable(variableId string) error {
	return c.Client.DeleteVariable(c.Options.WorkspaceId, variableId)
}

// Trigger CRUD

func (c *ClientInWorkspace) CreateTrigger(trigger *tagmanager.Trigger) (*tagmanager.Trigger, error) {
	return c.Client.CreateTrigger(c.Options.WorkspaceId, trigger)
}

func (c *ClientInWorkspace) ListTriggers() ([]*tagmanager.Trigger, error) {
	return c.Client.ListTriggers(c.Options.WorkspaceId)
}

func (c *ClientInWorkspace) Trigger(triggerId string) (*tagmanager.Trigger, error) {
	return c.Client.Trigger(c.Options.WorkspaceId, triggerId)
}

func (c *ClientInWorkspace) UpdateTrigger(triggerId string, trigger *tagmanager.Trigger) (*tagmanager.Trigger, error) {
	return c.Client.UpdateTrigger(c.Options.WorkspaceId, triggerId, trigger)
}

func (c *ClientInWorkspace) DeleteTrigger(triggerId string) error {
	return c.Client.DeleteTrigger(c.Options.WorkspaceId, triggerId)
}
