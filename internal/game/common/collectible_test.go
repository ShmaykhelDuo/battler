package common_test

import (
	"fmt"
	"testing"

	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/common"
	"github.com/stretchr/testify/assert"
)

func TestNewCollectible(t *testing.T) {
	t.Parallel()

	tests := []struct {
		amount int
	}{
		{
			amount: 10,
		},
		{
			amount: 50,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprint(tt.amount), func(t *testing.T) {
			t.Parallel()

			eff := common.NewCollectible(tt.amount)

			assert.Equal(t, tt.amount, eff.Amount())
		})
	}
}

func TestCollectible_Increase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		initAmount int
		incrAmount int
		amount     int
	}{
		{
			initAmount: 5,
			incrAmount: 5,
			amount:     10,
		},
		{
			initAmount: 4,
			incrAmount: 3,
			amount:     7,
		},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("%d+%d=%d", tt.initAmount, tt.incrAmount, tt.amount)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			eff := common.NewCollectible(tt.initAmount)
			eff.Increase(tt.incrAmount)
			assert.Equal(t, tt.amount, eff.Amount())
		})
	}
}

func TestCollectible_Decrease(t *testing.T) {
	t.Parallel()

	tests := []struct {
		initAmount int
		decrAmount int
		amount     int
	}{
		{
			initAmount: 10,
			decrAmount: 5,
			amount:     5,
		},
		{
			initAmount: 7,
			decrAmount: 3,
			amount:     4,
		},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("%d-%d=%d", tt.initAmount, tt.decrAmount, tt.amount)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			eff := common.NewCollectible(tt.initAmount)
			eff.Decrease(tt.decrAmount)
			assert.Equal(t, tt.amount, eff.Amount())
		})
	}
}

func TestCollectible_HasExpired(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		amount    int
		isExpired bool
	}{
		{
			name:      "Positive",
			amount:    5,
			isExpired: false,
		},
		{
			name:      "Zero",
			amount:    0,
			isExpired: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			eff := common.NewCollectible(tt.amount)

			assert.Equal(t, tt.isExpired, eff.HasExpired(game.TurnState{}))
		})
	}
}
