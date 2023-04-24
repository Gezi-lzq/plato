package logger

import (
	"context"
	"testing"
	"time"

	"github.com/Gezi-lzq/plato/common/config"
)

func TestLogger(t *testing.T) {
	config.Init("../../plato.yaml")
	NewLogger(WithLogDir("/Users/www/logs"))
	InfoCtx(context.Background(), "info test")
	DebugCtx(context.Background(), "debug test")
	WarnCtx(context.Background(), "warn test")
	ErrorCtx(context.Background(), "error test")
	time.Sleep(1 * time.Second)
}
