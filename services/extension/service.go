package extension

import (
	"context"
	"fmt"

	"github.com/estafette/estafette-extension-helmfile/clients/credentials"
	"github.com/estafette/estafette-extension-helmfile/clients/helmfile"
	"github.com/estafette/estafette-extension-helmfile/clients/kind"
)

type Service interface {
	ExecuteAction(ctx context.Context, action Action) (err error)
	Lint(ctx context.Context) (err error)
	Diff(ctx context.Context) (err error)
	Apply(ctx context.Context) (err error)
	Init(ctx context.Context) (err error)
}

// NewService returns a new extension.Service
func NewService(ctx context.Context, credentialsClient credentials.Client, kindClient kind.Client, helmfileClient helmfile.Client) (Service, error) {
	if credentialsClient == nil {
		return nil, fmt.Errorf("credentialsClient is nil, this is now allowed")
	}
	if kindClient == nil {
		return nil, fmt.Errorf("kindClient is nil, this is now allowed")
	}
	if helmfileClient == nil {
		return nil, fmt.Errorf("helmfileClient is nil, this is now allowed")
	}

	return &service{
		credentialsClient: credentialsClient,
		kindClient:        kindClient,
		helmfileClient:    helmfileClient,
	}, nil
}

type service struct {
	credentialsClient credentials.Client
	kindClient        kind.Client
	helmfileClient    helmfile.Client
}

func (s *service) ExecuteAction(ctx context.Context, action Action) (err error) {
	switch action {
	case ActionLint:
		return s.Lint(ctx)

	case ActionDiff:
		return s.Diff(ctx)

	case ActionApply:
		return s.Apply(ctx)

	default:
		return fmt.Errorf("action %v is not supported", action)
	}
}

func (s *service) Lint(ctx context.Context) (err error) {
	err = s.initCredentials(ctx)
	if err != nil {
		return
	}

	return s.helmfileClient.Lint(ctx)
}

func (s *service) Diff(ctx context.Context) (err error) {
	err = s.Init(ctx)
	if err != nil {
		return
	}

	return s.helmfileClient.Diff(ctx)
}

func (s *service) Apply(ctx context.Context) (err error) {
	err = s.Init(ctx)
	if err != nil {
		return
	}

	return s.helmfileClient.Apply(ctx)
}

func (s *service) initCredentials(ctx context.Context) (err error) {
	// extract credentials and write to location set in envvar GOOGLE_APPLICATION_CREDENTIALS
	err = s.credentialsClient.Init(ctx)
	if err != nil {
		return
	}

	return nil
}

func (s *service) initKindHost(ctx context.Context) (err error) {
	// wait for kind host to be ready
	err = s.kindClient.WaitForReadiness(ctx)
	if err != nil {
		return
	}

	// configure .kube/config for using kind host
	err = s.kindClient.PrepareKubeConfig(ctx)
	if err != nil {
		return
	}

	return nil
}

func (s *service) Init(ctx context.Context) (err error) {
	err = s.initCredentials(ctx)
	if err != nil {
		return
	}

	err = s.initKindHost(ctx)
	if err != nil {
		return
	}

	return nil
}
