package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
}
