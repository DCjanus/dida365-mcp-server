package utils

import (
	"buf.build/go/protoyaml"
	"fmt"
	"os"

	"github.com/bufbuild/protovalidate-go"
	"github.com/dcjanus/dida365-mcp-server/gen/proto/configuration"
	"github.com/pkg/errors"
)

func LoadConfig(path string) (*configuration.Config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read config file")
	}
	text := string(content)
	expandedText := os.ExpandEnv(text)
	config := &configuration.Config{
		Server: &configuration.Server{
			Listen: "localhost:8080",
		},
		Logging: &configuration.Logging{
			Level: "info",
		},
	}
	if err := (protoyaml.UnmarshalOptions{DiscardUnknown: true}).Unmarshal([]byte(expandedText), config); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal config")
	}

	masked := ProtoClone(config)
	masked.GetOauth().ClientSecret = "[REDACTED]"
	dumped, err := (protoyaml.MarshalOptions{}).Marshal(masked)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal config")
	}
	fmt.Printf("using config: \n%s\n", dumped)

	if err := protovalidate.Validate(config); err != nil {
		return nil, errors.Wrap(err, "failed to validate config")
	}
	return config, nil
}
