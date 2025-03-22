package utils

import (
	"os"

	"github.com/cockroachdb/errors"
	prettyconsole "github.com/thessem/zap-prettyconsole"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/dcjanus/dida365-mcp-server/gen/conf"
)

func NewLogger(cfg *conf.Logging) (*zap.Logger, error) {
	config := prettyconsole.NewEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder

	level, err := zapcore.ParseLevel(cfg.GetLevel())
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse log level")
	}

	logger := zap.New(zapcore.NewCore(
		prettyconsole.NewEncoder(config),
		os.Stdout,
		level,
	))
	return logger, nil
}
