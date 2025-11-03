package migrator

import (
	"context"
)

type Migrator interface {
	Up(ctx context.Context) error
	Down(ctx context.Context) error
}
