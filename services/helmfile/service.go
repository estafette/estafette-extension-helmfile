package helmfile

import (
	"context"
	"fmt"

	"github.com/estafette/estafette-extension-helmfile/clients/credentials"
	"github.com/estafette/estafette-extension-helmfile/services/kind"
	foundation "github.com/estafette/estafette-foundation"
)

type Service interface {
	Lint(ctx context.Context) (err error)
	Diff(ctx context.Context) (err error)
	Apply(ctx context.Context) (err error)
}

// NewService returns a new orchestrator.Service
func NewService(ctx context.Context, credentialsClient credentials.Client, kindService kind.Service, file string) (Service, error) {
	if kindService == nil {
		return nil, fmt.Errorf("kindService is nil, this is now allowed")
	}
	if file == "" {
		return nil, fmt.Errorf("file is empty, this is now allowed")
	}

	return &service{
		credentialsClient: credentialsClient,
		kindService:       kindService,
		file:              file,
	}, nil
}

type service struct {
	credentialsClient credentials.Client
	kindService       kind.Service
	file              string
}

func (s *service) Lint(ctx context.Context) (err error) {
	err = s.init(ctx)
	if err != nil {
		return
	}

	return foundation.RunCommandExtended(ctx, "helmfile --file %v lint", s.file)
}

func (s *service) Diff(ctx context.Context) (err error) {
	err = s.init(ctx)
	if err != nil {
		return
	}

	return foundation.RunCommandExtended(ctx, "helmfile --file %v diff", s.file)
}

func (s *service) Apply(ctx context.Context) (err error) {
	err = s.init(ctx)
	if err != nil {
		return
	}

	return foundation.RunCommandExtended(ctx, "helmfile --file %v apply", s.file)
}

func (s *service) init(ctx context.Context) (err error) {
	err = s.credentialsClient.Init(ctx)
	if err != nil {
		return
	}

	err = s.kindService.WaitForReadiness(ctx)
	if err != nil {
		return
	}
	err = s.kindService.PrepareKubeConfig(ctx)
	if err != nil {
		return
	}

	return nil
}
