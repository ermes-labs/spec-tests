package spec

import (
	"testing"

	"github.com/ermes-labs/api-go/api"
)

type Env[T api.Commands] interface {
	New(node string) (T, func())
}

func RunTests[T api.Commands](t *testing.T, env Env[T]) {
	TestAcquireSession(t, env)
	TestCreateSession(t, env)
	TestOffloadSession(t, env)
}
