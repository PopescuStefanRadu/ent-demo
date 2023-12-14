package entwrap

import (
	"context"
	"github.com/PopescuStefanRadu/ent-demo/pkg/ent"
)

type Migrator struct {
	Ent *ent.Client
}

func (m Migrator) Migrate(ctx context.Context) error {
	return m.Ent.Schema.Create(ctx)
}
