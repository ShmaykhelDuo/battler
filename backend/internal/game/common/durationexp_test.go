package common_test

import (
	"testing"

	"github.com/ShmaykhelDuo/battler/backend/internal/game"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/common"
	"github.com/stretchr/testify/assert"
)

func TestDurationExpirable_TurnsLeft(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		expCtx    game.Context
		gameCtx   game.Context
		turnsLeft int
	}{
		{
			name: "ExpiresThisTurn",
			expCtx: game.Context{
				TurnNum:      5,
				IsGoingFirst: true,
			},
			gameCtx: game.Context{
				TurnNum:      5,
				IsGoingFirst: true,
			},
			turnsLeft: 1,
		},
		{
			name: "ExpiresOpponentsNextTurnGoingFirst",
			expCtx: game.Context{
				TurnNum:      5,
				IsGoingFirst: false,
			},
			gameCtx: game.Context{
				TurnNum:      5,
				IsGoingFirst: true,
			},
			turnsLeft: 1,
		},
		{
			name: "ExpiresOpponentsNextTurnGoingSecond",
			expCtx: game.Context{
				TurnNum:      5,
				IsGoingFirst: true,
			},
			gameCtx: game.Context{
				TurnNum:      4,
				IsGoingFirst: false,
			},
			turnsLeft: 1,
		},
		{
			name: "ExpiresYourNextTurn",
			expCtx: game.Context{
				TurnNum:      6,
				IsGoingFirst: true,
			},
			gameCtx: game.Context{
				TurnNum:      5,
				IsGoingFirst: true,
			},
			turnsLeft: 2,
		},
		{
			name: "ExpiresOpponentsSecondNextTurnGoingFirst",
			expCtx: game.Context{
				TurnNum:      6,
				IsGoingFirst: false,
			},
			gameCtx: game.Context{
				TurnNum:      5,
				IsGoingFirst: true,
			},
			turnsLeft: 2,
		},
		{
			name: "ExpiresOpponentsSecondNextTurnGoingSecond",
			expCtx: game.Context{
				TurnNum:      7,
				IsGoingFirst: true,
			},
			gameCtx: game.Context{
				TurnNum:      5,
				IsGoingFirst: false,
			},
			turnsLeft: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			exp := common.NewDurationExpirable(tt.expCtx)

			turnsLeft := exp.TurnsLeft(tt.gameCtx)
			assert.Equal(t, tt.turnsLeft, turnsLeft)
		})
	}
}
