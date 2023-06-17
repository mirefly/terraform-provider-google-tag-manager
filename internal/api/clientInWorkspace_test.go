package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/tagmanager/v2"
)

var testClientInWorkspaceOptions = &ClientInWorkspaceOptions{
	WorkspaceName: "test-client-in-workspace",
	ClientOptions: testClientOptions,
}

func TestNewClientInWorkspace(t *testing.T) {
	client, err := NewClientInWorkspace(testClientInWorkspaceOptions)
	if err == nil {
		defer func() {
			err := client.DeleteWorkspace(client.Options.WorkspaceId)
			if err != nil {
				t.Error(err)
			}
		}()
	}

	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotZero(t, client.Options.WorkspaceId)

	_, err = client.CreateTrigger(&tagmanager.Trigger{
		Name:  "test-trigger-2",
		Type:  "click",
		Notes: "updated by unit test",
		Parameter: []*tagmanager.Parameter{
			{Key: "clickText", Value: "Button", Type: "template"},
		},
	})
	assert.NoError(t, err)
}
