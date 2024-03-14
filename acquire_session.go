package spec

import (
	"context"
	"testing"

	"github.com/ermes-labs/api-go/api"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestAcquireSession[T api.Commands](t *testing.T, env Env[T]) {
	// Set up the environment.
	cmd, free := env.New("node")
	defer free()

	// Scan the offloadable sessions.
	keys, cursor, err := cmd.ScanOffloadableSessions(context.Background(), 0, 10)
	// Check the result.
	if err != nil {
		t.Errorf("failed to scan the offloadable sessions: %v", err)
	} else {
		if len(keys) != 0 {
			t.Errorf("invalid scanned offloadable sessions, expected %v, found %v", []string{}, keys)
		}
		if cursor != 0 {
			t.Errorf("Expected \"0\" cursor, found %v", cursor)
		}
	}

	// Acquire an invalid session.
	offloadedTo, err := cmd.AcquireSession(context.Background(), "sessionId", api.AcquireSessionOptions{})
	// Check the result.
	if err == nil {
		t.Errorf("a non-existing session is illegally acquired")
	} else if err != api.ErrSessionNotFound {
		// TODO:
		// t.Errorf("acquiring a non-existing session should return ErrSessionNotFound, got %v", err)
	}
	// Check the offloadedTo.
	if offloadedTo != nil {
		t.Errorf("a non-existing session cannot have been offloaded, got %v", *offloadedTo)
	}

	// Scan the offloadable sessions.
	keys, cursor, err = cmd.ScanOffloadableSessions(context.Background(), 0, 10)
	// Check the result.
	if err != nil {
		t.Errorf("failed to scan the offloadable sessions: %v", err)
	} else {
		if len(keys) != 0 {
			t.Errorf("invalid scanned offloadable sessions, expected %v, found %v", []string{}, keys)
		}
		if cursor != 0 {
			t.Errorf("Expected \"0\" cursor, found %v", cursor)
		}
	}

	// Release non-existing session.
	offloadedTo, err = cmd.ReleaseSession(context.Background(), "sessionId", api.AcquireSessionOptions{})
	// Check the result.
	if err == nil {
		t.Errorf("a non-existing session is illegally released")
	} else if err != api.ErrSessionNotFound {
		// t.Errorf("releasing a non-existing session should return ErrSessionNotFound, got %v", err)
	}
	if offloadedTo != nil {
		t.Errorf("a non-existing session cannot have been offloaded, got %v", *offloadedTo)
	}

	// Create a session.
	sessionId, err := cmd.CreateSession(context.Background(), api.CreateSessionOptions{})
	// Check the result.
	if err != nil {
		t.Errorf("failed to create a session: %v", err)
	}

	offloadedTo, err = cmd.AcquireSession(context.Background(), sessionId, api.NewAcquireSessionOptionsBuilder().AllowOffloading().Build())
	// Check the result.
	if err != nil {
		t.Errorf("failed to acquire a session: %v", err)
	}
	// Check the offloadedTo.
	if offloadedTo != nil {
		t.Errorf("a session should not have been offloaded, got %v", *offloadedTo)
	}

	// Scan the offloadable sessions.
	keys, cursor, err = cmd.ScanOffloadableSessions(context.Background(), 0, 10)
	// Check the result.
	if err != nil {
		t.Errorf("failed to scan the offloadable sessions: %v", err)
	} else {
		if len(keys) != 1 || keys[0] != sessionId {
			t.Errorf("invalid scanned offloadable sessions, expected %v, found %v", []string{sessionId}, keys)
		}
		if cursor != 0 {
			t.Errorf("Expected \"0\" cursor, found %v", cursor)
		}
	}

	// Create a session.
	sessionId2, err := cmd.CreateSession(context.Background(), api.CreateSessionOptions{})
	// Check the result.
	if err != nil {
		t.Errorf("failed to create a session: %v", err)
	}

	offloadedTo, err = cmd.AcquireSession(context.Background(), sessionId2, api.NewAcquireSessionOptionsBuilder().Build())
	// Check the result.
	if err != nil {
		t.Errorf("failed to acquire a session: %v", err)
	}
	// Check the offloadedTo.
	if offloadedTo != nil {
		t.Errorf("a session should not have been offloaded, got %v", *offloadedTo)
	}

	offloadedTo, err = cmd.AcquireSession(context.Background(), sessionId, api.NewAcquireSessionOptionsBuilder().Build())
	// Check the result.
	if err != nil {
		t.Errorf("failed to acquire a session: %v", err)
	}
	// Check the offloadedTo.
	if offloadedTo != nil {
		t.Errorf("a session should not have been offloaded, got %v", *offloadedTo)
	}

	// Scan the offloadable sessions.
	keys, cursor, err = cmd.ScanOffloadableSessions(context.Background(), 0, 10)
	// Check the result.
	if err != nil {
		t.Errorf("failed to scan the offloadable sessions: %v", err)
	} else {
		if len(keys) != 0 {
			t.Errorf("invalid scanned offloadable sessions, expected %v, found %v", []string{}, keys)
		}
		if cursor != 0 {
			t.Errorf("Expected \"0\" cursor, found %v", cursor)
		}
	}
}
