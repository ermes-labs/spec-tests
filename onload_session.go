package spec

import (
	"context"
	"testing"

	"github.com/ermes-labs/api-go/api"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestOnloadSession[T api.Commands](t *testing.T, env Env[T]) {
	// Set up the environment.
	cmd1, free1 := env.New("node1")
	cmd2, free2 := env.New("node2")
	defer free1()
	defer free2()

	// Create a session.
	sessionId, err := cmd1.CreateSession(context.Background(), api.CreateSessionOptions{})
	// Check the result.
	if err != nil {
		t.Errorf("failed to create a session: %v", err)
	}

	metadata, err := cmd1.GetSessionMetadata(context.Background(), sessionId)
	// Check the result.
	if err != nil {
		t.Errorf("failed to get session metadata: %v", err)
	}

	// TODO: wrap api.Commands and require some methods to enforce the loader signature.
	readCloser, _, err := cmd1.OffloadSession(context.Background(), sessionId, api.NewOffloadSessionOptionsBuilder().Build())
	// Check the result.
	if err != nil {
		t.Errorf("failed to offload session: %v", err)
	}

	cmd2.OnloadSession(context.Background(), metadata, readCloser, api.NewOnloadSessionOptionsBuilder().Build())
}
