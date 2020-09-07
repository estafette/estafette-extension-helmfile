package helmfile

import (
	"context"
	"fmt"

	foundation "github.com/estafette/estafette-foundation"
)

type Client interface {
	Lint(ctx context.Context) (err error)
	Diff(ctx context.Context) (err error)
	Apply(ctx context.Context) (err error)
}

// NewClient returns a new helmfile.Client
func NewClient(ctx context.Context, file, logLevel string) (Client, error) {
	if file == "" {
		return nil, fmt.Errorf("file is empty, this is now allowed")
	}
	if logLevel == "" {
		logLevel = "info"
	}

	return &client{
		file: file,
	}, nil
}

type client struct {
	file     string
	logLevel string
}

func (c *client) Lint(ctx context.Context) (err error) {
	return foundation.RunCommandExtended(ctx, "helmfile --file %v --log-level %v lint", c.file, c.logLevel)
}

func (c *client) Diff(ctx context.Context) (err error) {
	return foundation.RunCommandExtended(ctx, "helmfile --file %v --log-level %v diff", c.file, c.logLevel)
}

func (c *client) Apply(ctx context.Context) (err error) {
	return foundation.RunCommandExtended(ctx, "helmfile --file %v --log-level %v apply", c.file, c.logLevel)
}
