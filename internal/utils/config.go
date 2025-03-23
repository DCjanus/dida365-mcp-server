package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"buf.build/go/protoyaml"
	"github.com/bufbuild/protovalidate-go"
	"github.com/cockroachdb/errors"
	"github.com/joho/godotenv"

	"github.com/dcjanus/dida365-mcp-server/gen/conf"
)

// LoadDotEnvs loads the *.env file in the root directory.
func LoadDotEnvs() {
	files, err := filepath.Glob("*.env")
	if err != nil {
		panic(errors.Wrap(err, "failed to scan current directory"))
	}
	if len(files) > 0 {
		if err := godotenv.Load(files...); err != nil {
			panic(errors.Wrap(err, "failed to load .env file"))
		}
	}
}

func LoadConfig(path string) (*conf.Config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read config file")
	}
	text := string(content)
	expandedText := os.ExpandEnv(text)
	config := &conf.Config{
		Server: &conf.Server{
			Listen: "localhost:8080",
		},
		Logging: &conf.Logging{
			Level: "info",
		},
	}
	if err := (protoyaml.UnmarshalOptions{DiscardUnknown: true}).Unmarshal([]byte(expandedText), config); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal config")
	}

	masked := ProtoClone(config)
	if masked.Oauth != nil && masked.Oauth.ClientSecret != "" {
		masked.Oauth.ClientSecret = "********"
	}
	dumped, err := (protoyaml.MarshalOptions{UseProtoNames: true, EmitUnpopulated: true}).Marshal(masked)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal config")
	}
	fmt.Printf("using config: \n%s\n", dumped)

	if err := protovalidate.Validate(config); err != nil {
		return nil, errors.Wrap(err, "failed to validate config")
	}
	return config, nil
}
