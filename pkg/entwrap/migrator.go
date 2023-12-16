package entwrap

import (
	"bytes"
	"context"

	"github.com/PopescuStefanRadu/ent-demo/pkg/ent"
	"github.com/rs/zerolog"
)

type Migrator struct {
	Ent    *ent.Client
	Logger zerolog.Logger
}

func (m Migrator) Migrate(ctx context.Context) error {
	buffer := &bytes.Buffer{}
	if err := m.Ent.Schema.WriteTo(ctx, buffer); err != nil {
		return err
	}

	m.Logger.Info().Msgf("Migrated schema with following changes: %s", buffer.String())

	return m.Ent.Schema.Create(ctx)
}
