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
func NewClient(ctx context.Context, file string) (Client, error) {
	if file == "" {
		return nil, fmt.Errorf("file is empty, this is now allowed")
	}

	return &client{
		file: file,
	}, nil
}

type client struct {
	file string
}

func (c *client) Lint(ctx context.Context) (err error) {
	return foundation.RunCommandExtended(ctx, "helmfile --file %v lint", c.file)
}

func (c *client) Diff(ctx context.Context) (err error) {
	return foundation.RunCommandExtended(ctx, "helmfile --file %v diff", c.file)
}

func (c *client) Apply(ctx context.Context) (err error) {
	return foundation.RunCommandExtended(ctx, "helmfile --file %v apply", c.file)
}
