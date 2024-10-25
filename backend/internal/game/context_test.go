package game_test

import (
	"testing"

	"github.com/ShmaykhelDuo/battler/backend/internal/game"
	"github.com/stretchr/testify/assert"
)

func TestContext_IsAfter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		ctx     game.TurnState
		other   game.TurnState
		isAfter bool
	}{
		{
			name: "LessTurn",
			ctx: game.TurnState{
				TurnNum: 5,
			},
			other: game.TurnState{
				TurnNum: 6,
			},
			isAfter: false,
		},
		{
			name: "GreaterTurn",
			ctx: game.TurnState{
				TurnNum: 6,
			},
			other: game.TurnState{
				TurnNum: 5,
			},
			isAfter: true,
		},
		{
			name: "EqualTurnLessGoingFirst",
			ctx: game.TurnState{
				TurnNum:      5,
				IsGoingFirst: true,
			},
			other: game.TurnState{
				TurnNum:      5,
				IsGoingFirst: false,
			},
			isAfter: false,
		},
		{
			name: "EqualTurnGreaterGoingFirst",
			ctx: game.TurnState{
				TurnNum:      5,
				IsGoingFirst: false,
			},
			other: game.TurnState{
				TurnNum:      5,
				IsGoingFirst: true,
			},
			isAfter: true,
		},
		{
			name: "EqualTurnAndGoingFirstLessTurnEnd",
			ctx: game.TurnState{
				TurnNum:      5,
				IsGoingFirst: true,
				IsTurnEnd:    false,
			},
			other: game.TurnState{
				TurnNum:      5,
				IsGoingFirst: true,
				IsTurnEnd:    true,
			},
			isAfter: false,
		},
		{
			name: "EqualTurnAndGoingFirstGreaterTurnEnd",
			ctx: game.TurnState{
				TurnNum:      5,
				IsGoingFirst: true,
				IsTurnEnd:    true,
			},
			other: game.TurnState{
				TurnNum:      5,
				IsGoingFirst: true,
				IsTurnEnd:    false,
			},
			isAfter: true,
		},
		{
			name: "EqualAllTurnEnd",
			ctx: game.TurnState{
				TurnNum:      5,
				IsGoingFirst: true,
				IsTurnEnd:    true,
			},
			other: game.TurnState{
				TurnNum:      5,
				IsGoingFirst: true,
				IsTurnEnd:    true,
			},
			isAfter: false,
		},
		{
			name: "EqualAllNotTurnEnd",
			ctx: game.TurnState{
				TurnNum:      5,
				IsGoingFirst: true,
				IsTurnEnd:    false,
			},
			other: game.TurnState{
				TurnNum:      5,
				IsGoingFirst: true,
				IsTurnEnd:    false,
			},
			isAfter: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			isAfter := tt.ctx.IsAfter(tt.other)
			assert.Equal(t, tt.isAfter, isAfter)
		})
	}
}

func TestContext_AddTurns(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		ctx             game.TurnState
		turns           int
		isOpponentsTurn bool
		res             game.TurnState
	}{
		{
			name: "AddNone",
			ctx: game.TurnState{
				TurnNum:      5,
				IsGoingFirst: false,
			},
			turns:           0,
			isOpponentsTurn: false,
			res: game.TurnState{
				TurnNum:      5,
				IsGoingFirst: false,
			},
		},
		{
			name: "NextOpponentsTurnGoingFirst",
			ctx: game.TurnState{
				TurnNum:      5,
				IsGoingFirst: true,
			},
			turns:           0,
			isOpponentsTurn: true,
			res: game.TurnState{
				TurnNum:      5,
				IsGoingFirst: false,
			},
		},
		{
			name: "NextOpponentsTurnNotGoingFirst",
			ctx: game.TurnState{
				TurnNum:      5,
				IsGoingFirst: false,
			},
			turns:           0,
			isOpponentsTurn: true,
			res: game.TurnState{
				TurnNum:      6,
				IsGoingFirst: true,
			},
		},
		{
			name: "AddSomeTurns",
			ctx: game.TurnState{
				TurnNum:      5,
				IsGoingFirst: false,
			},
			turns:           3,
			isOpponentsTurn: false,
			res: game.TurnState{
				TurnNum:      8,
				IsGoingFirst: false,
			},
		},
		{
			name: "AddSomeTurnsWithOpponentsTurn",
			ctx: game.TurnState{
				TurnNum:      5,
				IsGoingFirst: false,
			},
			turns:           3,
			isOpponentsTurn: true,
			res: game.TurnState{
				TurnNum:      9,
				IsGoingFirst: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			res := tt.ctx.AddTurns(tt.turns, tt.isOpponentsTurn)
			assert.Equal(t, tt.res, res)
		})
	}
}
