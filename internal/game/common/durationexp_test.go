package common_test

import (
	"testing"

	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/common"
	"github.com/stretchr/testify/assert"
)

func TestDurationExpirable_TurnsLeft(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		expCtx    game.TurnState
		turnState game.TurnState
		turnsLeft int
	}{
		{
			name: "ExpiresThisTurn",
			expCtx: game.TurnState{
				TurnNum:      5,
				IsGoingFirst: true,
			},
			turnState: game.TurnState{
				TurnNum:      5,
				IsGoingFirst: true,
			},
			turnsLeft: 1,
		},
		{
			name: "ExpiresOpponentsNextTurnGoingFirst",
			expCtx: game.TurnState{
				TurnNum:      5,
				IsGoingFirst: false,
			},
			turnState: game.TurnState{
				TurnNum:      5,
				IsGoingFirst: true,
			},
			turnsLeft: 1,
		},
		{
			name: "ExpiresOpponentsNextTurnGoingSecond",
			expCtx: game.TurnState{
				TurnNum:      5,
				IsGoingFirst: true,
			},
			turnState: game.TurnState{
				TurnNum:      4,
				IsGoingFirst: false,
			},
			turnsLeft: 1,
		},
		{
			name: "ExpiresYourNextTurn",
			expCtx: game.TurnState{
				TurnNum:      6,
				IsGoingFirst: true,
			},
			turnState: game.TurnState{
				TurnNum:      5,
				IsGoingFirst: true,
			},
			turnsLeft: 2,
		},
		{
			name: "ExpiresOpponentsSecondNextTurnGoingFirst",
			expCtx: game.TurnState{
				TurnNum:      6,
				IsGoingFirst: false,
			},
			turnState: game.TurnState{
				TurnNum:      5,
				IsGoingFirst: true,
			},
			turnsLeft: 2,
		},
		{
			name: "ExpiresOpponentsSecondNextTurnGoingSecond",
			expCtx: game.TurnState{
				TurnNum:      7,
				IsGoingFirst: true,
			},
			turnState: game.TurnState{
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

			turnsLeft := exp.TurnsLeft(tt.turnState)
			assert.Equal(t, tt.turnsLeft, turnsLeft)
		})
	}
}
