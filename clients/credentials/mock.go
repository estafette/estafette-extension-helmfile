package credentials

import (
	"context"
)

type MockClient struct {
	InitFunc func(ctx context.Context) (err error)
}

func (c MockClient) Init(ctx context.Context) (err error) {
	if c.InitFunc == nil {
		return
	}

	return c.InitFunc(ctx)
}
