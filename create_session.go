package spec

import (
	"context"
	"testing"

	"github.com/ermes-labs/api-go/api"
)

func TestCreateSession[T api.Commands](t *testing.T, env Env[T]) {
	// Set up the environment.
	cmd, free := env.New("node")
	defer free()

	// Get the sessions.
	sessions, cursor, err := cmd.ScanSessions(context.Background(), 0, 10)
	// Check the result.
	if err != nil {
		t.Errorf("failed to scan the sessions: %v", err)
	} else {
		if len(sessions) != 0 {
			t.Errorf("invalid scanned sessions, expected %v, found %v", []string{}, sessions)
		}
		if cursor != 0 {
			t.Errorf("Expected \"0\" cursor, found %v", cursor)
		}
	}

	// Create a session.
	sessionId, err := cmd.CreateSession(context.Background(), api.CreateSessionOptions{})
	// Check the result.
	if err != nil {
		t.Errorf("failed to create a session: %v", err)
	}

	// Get the sessions.
	sessions, cursor, err = cmd.ScanSessions(context.Background(), 0, 10)
	// Check the result.
	if err != nil {
		t.Errorf("failed to scan the sessions: %v", err)
	} else {
		if len(sessions) != 1 || sessions[0] != sessionId {
			t.Errorf("invalid scanned sessions, expected %v, found %v", []string{sessionId}, sessions)
		}
		if cursor != 0 {
			t.Errorf("Expected \"0\" cursor, found %v", cursor)
		}
	}

	// Create a session.
	_, err = cmd.CreateSession(context.Background(), api.NewCreateSessionOptionsBuilder().SessionId(sessionId).Build())
	// Check the result.
	if err == nil {
		t.Errorf("a session with the same id should not have been created: %v", err)
	} else if err != api.ErrSessionIdAlreadyExists {
		// t.Errorf("creating a session with the same id should return ErrSessionIdAlreadyExists, got %v", err)
	}

	// Get the sessions.
	sessions, cursor, err = cmd.ScanSessions(context.Background(), 0, 10)
	// Check the result.
	if err != nil {
		t.Errorf("failed to scan the sessions: %v", err)
	} else {
		if len(sessions) != 1 || sessions[0] != sessionId {
			t.Errorf("invalid scanned sessions, expected %v, found %v", []string{sessionId}, sessions)
		}
		if cursor != 0 {
			t.Errorf("Expected \"0\" cursor, found %v", cursor)
		}
	}
}
