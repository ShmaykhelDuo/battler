package auth_test

import (
	"context"
	"testing"

	model "github.com/ShmaykhelDuo/battler/internal/model/auth"
	"github.com/ShmaykhelDuo/battler/internal/pkg/auth"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSession(t *testing.T) {
	t.Parallel()

	t.Run("no session", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		_, err := auth.Session(ctx)
		assert.ErrorIs(t, err, model.ErrNoSession)
	})

	t.Run("session", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		s := model.Session{
			ID:     uuid.MustParse("a1583a16-48bb-45b2-9fbd-66370b567f38"),
			UserID: uuid.MustParse("48440578-dae4-4b6d-9fda-c7ce7a2a2b2a"),
		}
		ctx = auth.ContextWithSession(ctx, s)

		got, err := auth.Session(ctx)
		assert.NoError(t, err)
		assert.Equal(t, s, got)
	})
}
