package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/tagmanager/v2"
)

var testClientOptions = &ClientOptions{
	CredentialFile: "./testdata/credentials-77b14e38b4dd.json",
	AccountId:      "6105084028",
	ContainerId:    "119458552",
}

func newTestClient(t *testing.T) *Client {
	client, err := NewClient(testClientOptions)

	assert.Nil(t, err)
	return client
}

func TestNewClient(t *testing.T) {
	client := newTestClient(t)
	assert.NotNil(t, client)
}

func currentTimeString() string {
	return time.Now().Format("2006-01-02-15-04-05")
}

func TestClientWorkSpaceCRUD(t *testing.T) {
	client := newTestClient(t)

	// Create workspace
	ws, err := client.CreateWorkspace(&tagmanager.Workspace{
		Name:        "test-workspace-CRUD-" + currentTimeString(),
		Description: "created by unit test",
	})
	assert.NoError(t, err)
	assert.NotNil(t, ws)

	// Get workspace
	fetched, err := client.Workspace(ws.WorkspaceId)
	assert.NoError(t, err)
	assert.Equal(t, ws.Name, fetched.Name)

	// List workspaces
	list, err := client.ListWorkspaces()
	assert.NoError(t, err)
	assert.Greater(t, len(list), 0)

	// Update workspace
	updated, err := client.UpdateWorkspaces(ws.WorkspaceId, &tagmanager.Workspace{
		Name:        "updated-workspace-" + currentTimeString(),
		Description: "updated by unit test",
	})
	assert.NoError(t, err)
	assert.Contains(t, updated.Name, "updated-workspace")

	// Delete workspace
	err = client.DeleteWorkspace(ws.WorkspaceId)
	assert.NoError(t, err)

	// Get nonexisting workspace
	ws, err = client.Workspace(ws.WorkspaceId)
	assert.Equal(t, ErrNotExist, err)
	assert.Nil(t, ws)
}

func TestClientVariableCRUD(t *testing.T) {
	client := newTestClient(t)
	ws, err := client.CreateWorkspace(&tagmanager.Workspace{
		Name:        "test-variable-CRUD-" + currentTimeString(),
		Description: "created by unit test",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer client.DeleteWorkspace(ws.WorkspaceId)

	// Create variable
	variable, err := client.CreateVariable(ws.WorkspaceId, &tagmanager.Variable{
		Name: "test-variable-1",
		Type: "v",
		Parameter: []*tagmanager.Parameter{
			{Key: "name", Type: "template", Value: "test"},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, "test-variable-1", variable.Name)

	// Get variable
	variable, err = client.Variable(ws.WorkspaceId, variable.VariableId)
	assert.NoError(t, err)
	assert.Equal(t, "test-variable-1", variable.Name)

	// Update variable
	variable, err = client.UpdateVariable(ws.WorkspaceId, variable.VariableId, &tagmanager.Variable{
		Name: "test-variable-2",
		Type: "v",
		Parameter: []*tagmanager.Parameter{
			{Key: "name", Type: "template", Value: "test"},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, "test-variable-2", variable.Name)

	// Delete variable
	err = client.DeleteVariable(ws.WorkspaceId, variable.VariableId)
	assert.NoError(t, err)
}

func TestClientTagCRUD(t *testing.T) {
	client := newTestClient(t)
	ws, err := client.CreateWorkspace(&tagmanager.Workspace{
		Name:        "test-tags-CRUD-" + currentTimeString(),
		Description: "created by unit test",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer client.DeleteWorkspace(ws.WorkspaceId)

	// Create tag
	tag, err := client.CreateTag(ws.WorkspaceId, &tagmanager.Tag{
		Name:  "test-tag-1",
		Notes: "created by unit test",
		Type:  "gaawe",
		Parameter: []*tagmanager.Parameter{
			{Key: "eventName", Value: "test", Type: "template"},
			{Key: "measurementId", Value: "test", Type: "template"},
			{Key: "eventParameters", Type: "list", List: []*tagmanager.Parameter{
				{Type: "map", Map: []*tagmanager.Parameter{
					{Key: "name", Type: "template", Value: "name-v"},
					{Key: "value", Type: "template", Value: "value-v"},
				}},
			}},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, "test-tag-1", tag.Name)

	// Get tag
	tag, err = client.Tag(ws.WorkspaceId, tag.TagId)
	assert.NoError(t, err)
	assert.Equal(t, "test-tag-1", tag.Name)

	// Update tag
	tag, err = client.UpdateTag(ws.WorkspaceId, tag.TagId, &tagmanager.Tag{
		Name:  "test-tag-2",
		Notes: "updated by unit test",
		Type:  "gaawe",
		Parameter: []*tagmanager.Parameter{
			{Key: "eventName", Value: "test", Type: "template"},
			{Key: "measurementId", Value: "test", Type: "template"},
		},
	})
	assert.NoError(t, err)

	// Delete tag
	err = client.DeleteTag(ws.WorkspaceId, tag.TagId)
	assert.NoError(t, err)
}

func TestClientTriggerCRUD(t *testing.T) {
	client := newTestClient(t)
	ws, err := client.CreateWorkspace(&tagmanager.Workspace{
		Name:        "test-triggers-CRUD-" + currentTimeString(),
		Description: "created by unit test",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer client.DeleteWorkspace(ws.WorkspaceId)

	// Create trigger
	trigger, err := client.CreateTrigger(ws.WorkspaceId, &tagmanager.Trigger{
		Name:  "test-trigger-1",
		Type:  "customEvent",
		Notes: "My Custom Event",
		CustomEventFilter: []*tagmanager.Condition{
			{
				Type: "equals",
				Parameter: []*tagmanager.Parameter{
					{Key: "arg0", Value: "{{_event}}", Type: "template"},
					{Key: "arg1", Value: "myEvent", Type: "template"},
				},
			},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, "test-trigger-1", trigger.Name)

	// Get trigger
	trigger, err = client.Trigger(ws.WorkspaceId, trigger.TriggerId)
	assert.NoError(t, err)
	assert.Equal(t, "test-trigger-1", trigger.Name)
	assert.Equal(t, "myEvent", trigger.CustomEventFilter[0].Parameter[1].Value)

	// Update trigger
	trigger, err = client.UpdateTrigger(ws.WorkspaceId, trigger.TriggerId, &tagmanager.Trigger{
		Name:  "test-trigger-2",
		Type:  "click",
		Notes: "updated by unit test",
		Parameter: []*tagmanager.Parameter{
			{Key: "clickText", Value: "Button", Type: "template"},
		},
	})
	assert.NoError(t, err)

	// Delete trigger
	err = client.DeleteTrigger(ws.WorkspaceId, trigger.TriggerId)
	assert.NoError(t, err)
}
