package utils

import (
	"os"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/jsternberg/zap-logfmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/dcjanus/dida365-mcp-server/gen/proto/configuration"
)

func NewLogger(cfg *configuration.Logging) (*zap.Logger, error) {
	config := zap.NewProductionEncoderConfig()

	level, err := zapcore.ParseLevel(cfg.GetLevel())
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse log level")
	}
	config.EncodeTime = func(ts time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(ts.Local().Format(time.RFC3339))
	}
	logger := zap.New(zapcore.NewCore(
		zaplogfmt.NewEncoder(config),
		os.Stdout,
		level,
	))
	return logger, nil
}
