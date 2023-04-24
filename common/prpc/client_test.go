package prpc

import (
	"testing"

	"github.com/Gezi-lzq/plato/common/config"
	ptrace "github.com/Gezi-lzq/plato/common/prpc/trace"
	"github.com/stretchr/testify/assert"
)

func TestNewPClient(t *testing.T) {
	config.Init("../../plato.yaml")
	ptrace.StartAgent()
	defer ptrace.StopAgent()

	_, err := NewPClient("plato_server")
	assert.NoError(t, err)
}
