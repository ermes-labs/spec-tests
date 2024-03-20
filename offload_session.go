package spec

import (
	"context"
	"testing"

	"github.com/ermes-labs/api-go/api"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestOffloadSession[T api.Commands](t *testing.T, env Env[T]) {
	// Set up the environment.
	cmd, free := env.New("node")
	defer free()

	// Create a session.
	sessionId, err := cmd.CreateSession(context.Background(), api.CreateSessionOptions{})
	// Check the result.
	if err != nil {
		t.Errorf("failed to create a session: %v", err)
	}

	offloadedTo, err := cmd.AcquireSession(context.Background(), sessionId, api.NewAcquireSessionOptionsBuilder().AllowOffloading().Build())
	// Check the result.
	if err != nil {
		t.Errorf("failed to acquire a session: %v", err)
	}
	// Check the offloadedTo.
	if offloadedTo != nil {
		t.Errorf("a session should not have been offloaded, got %v", *offloadedTo)
	}

	// TODO: wrap api.Commands and require some methods to enforce the loader signature.
	// TODO: wrap api.Commands and require some methods to check the session data.
	// Offload the session.
	offloadSessionOptions := api.NewOffloadSessionOptionsBuilder().Build()
	_, _, err = cmd.OffloadSession(context.Background(), sessionId, offloadSessionOptions)
	// Check the result.
	if err != nil {
		t.Errorf("failed to offload a session: %v", err)
	}

	_, _, err = cmd.OffloadSession(context.Background(), sessionId, offloadSessionOptions)
	// Check the result.
	if err == nil {
		t.Errorf("a session should not have been offloaded twice, got %v", err)
	} // else if err != api.ErrSessionAlreadyOffloaded {
	// 	t.Errorf("offloading a session twice should return ErrSessionAlreadyOffloaded, got %v", err)
	// }

	// Acquire the offloaded session.
	offloadedTo, err = cmd.AcquireSession(context.Background(), sessionId, api.NewAcquireSessionOptionsBuilder().Build())
	// Check the result.
	if err == nil {
		t.Errorf("acquiring an offloading session should return an error, got %v", err)
	} //else if err != api.ErrSessionIsOffloading {
	// t.Errorf("acquiring an offloading session should return ErrSessionIsOffloading, got %v", err)
	// }
	if offloadedTo != nil {
		t.Errorf("session should be offloading, got offloadedTo %v", *offloadedTo)
	}

	ids, cursor, err := cmd.ScanOffloadedSessions(context.Background(), 0, 10)
	// Check the result.
	if err != nil {
		t.Errorf("failed to scan the offloadable sessions: %v", err)
	} else {
		if len(ids) != 0 {
			t.Errorf("invalid scanned offloadable sessions, expected %v, found %v", []string{}, ids)
		}
		if cursor != 0 {
			t.Errorf("Expected \"0\" cursor, found %v", cursor)
		}
	}

	newLocation := api.SessionLocation{}
	cmd.ConfirmSessionOffload(context.Background(), sessionId, newLocation, offloadSessionOptions, func(ctx context.Context, oldLocation api.SessionLocation) (bool, error) {
		return false, nil
	})

	ids, cursor, err = cmd.ScanOffloadedSessions(context.Background(), 0, 10)
	// Check the result.
	if err != nil {
		t.Errorf("failed to scan the offloadable sessions: %v", err)
	} else {
		if len(ids) != 1 || ids[0] != sessionId {
			t.Errorf("invalid scanned offloadable sessions, expected %v, found %v", []string{sessionId}, ids)
		}
		if cursor != 0 {
			t.Errorf("Expected \"0\" cursor, found %v", cursor)
		}
	}
}
