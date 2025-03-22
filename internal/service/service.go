package service

import (
	"context"
	"net/url"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"resty.dev/v3"

	"github.com/dcjanus/dida365-mcp-server/gen/api"
	"github.com/dcjanus/dida365-mcp-server/gen/conf"
	"github.com/dcjanus/dida365-mcp-server/gen/model"
)

type Dida365oAuthService struct {
	log *zap.Logger
	cfg *conf.Config
	cli *resty.Client
	api.UnimplementedDida365OAuthServiceServer
}

func NewDida365AuthService(log *zap.Logger, cfg *conf.Config) *Dida365oAuthService {
	cli := resty.New().SetHeader("User-Agent", "dida365-mcp-server")
	return &Dida365oAuthService{log: log, cfg: cfg, cli: cli}
}

func (d *Dida365oAuthService) Ping(ctx context.Context, _ *emptypb.Empty) (*wrapperspb.StringValue, error) {
	if err := grpc.SetHeader(ctx, map[string][]string{"key": {"value"}}); err != nil {
		return nil, errors.Wrap(err, "failed to set header")
	}
	return &wrapperspb.StringValue{Value: "Pong"}, nil
}

func (d *Dida365oAuthService) OAuthLogin(ctx context.Context, _ *emptypb.Empty) (*model.TemporaryRedirectResponse, error) {
	base, err := url.Parse("https://dida365.com/oauth/authorize")
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse base URL")
	}
	query := base.Query()
	query.Set("client_id", d.cfg.GetOauth().GetClientId())
	query.Set("scope", "tasks:write tasks:read")
	query.Set("state", uuid.NewString())
	query.Set("redirect_uri", d.cfg.GetOauth().GetRedirectUri())
	query.Set("response_type", "code")
	base.RawQuery = query.Encode()

	return &model.TemporaryRedirectResponse{Location: base.String()}, nil
}

func (d *Dida365oAuthService) OAuthCallback(ctx context.Context, req *api.OAuthCallbackRequest) (*api.OAuthCallbackResponse, error) {
	reply := struct {
		AccessToken string `json:"access_token"`
	}{}

	res, err := d.cli.
		R().
		SetBasicAuth(d.cfg.GetOauth().GetClientId(), d.cfg.GetOauth().GetClientSecret()).
		SetFormData(map[string]string{
			"code":         req.GetCode(),
			"grant_type":   "authorization_code",
			"scope":        "tasks:write tasks:read",
			"redirect_uri": d.cfg.GetOauth().GetRedirectUri(),
		}).
		SetResult(&reply).
		Post("https://dida365.com/oauth/token")
	if err != nil {
		return nil, errors.Wrap(err, "failed to request oauth token")
	}
	if res.IsError() {
		return nil, errors.Errorf("failed to request oauth token, status: %d, body: %s", res.StatusCode(), res.String())
	}

	return &api.OAuthCallbackResponse{AccessToken: reply.AccessToken, State: req.GetState()}, nil
}
