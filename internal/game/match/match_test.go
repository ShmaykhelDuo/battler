package match_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/milana"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/z89"
	"github.com/ShmaykhelDuo/battler/internal/game/match"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type FakePlayer struct {
	moves    []int
	hasState bool
	hasEnd   bool
}

func (p *FakePlayer) SendState(ctx context.Context, state match.GameState) error {
	p.hasState = true
	return nil
}

func (p *FakePlayer) SendError(ctx context.Context, err error) error {
	return fmt.Errorf("got error in fake: %w", err)
}

func (p *FakePlayer) SendEnd(ctx context.Context) error {
	p.hasEnd = true
	return nil
}

func (p *FakePlayer) RequestSkill(ctx context.Context) (int, error) {
	if !p.hasState {
		return 0, errors.New("has no state")
	}

	if len(p.moves) == 0 {
		return 0, errors.New("no moves left")
	}

	m := p.moves[0]
	p.moves = p.moves[1:]

	return m, nil
}

func (p *FakePlayer) GivenUp() <-chan any {
	return nil
}

func TestMatch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		char1, char2   *game.CharacterData
		moves1, moves2 []int
		invertedOrder  bool
		result         match.Result
		hp1, hp2       int
	}{
		{
			name:  "FastEnd",
			char1: z89.CharacterZ89,
			char2: milana.CharacterMilana,
			moves1: []int{
				z89.SkillScarcityIndex,
				z89.SkillIndifferenceIndex,
				z89.SkillScarcityIndex,
				z89.SkillGreenSphereIndex,
				z89.SkillIndifferenceIndex,
				z89.SkillScarcityIndex,
				z89.SkillGreenSphereIndex,
				z89.SkillScarcityIndex,
				z89.SkillDespondencyIndex,
			},
			moves2: []int{
				milana.SkillMintMistIndex,
				milana.SkillRoyalMoveIndex,
				milana.SkillRoyalMoveIndex,
				milana.SkillComposureIndex,
				milana.SkillMintMistIndex,
				milana.SkillRoyalMoveIndex,
				milana.SkillRoyalMoveIndex,
				milana.SkillRoyalMoveIndex,
			},
			invertedOrder: false,
			result: match.Result{
				Player1: match.ResultPlayer{
					Status: match.ResultStatusWon,
				},
				Player2: match.ResultPlayer{
					Status: match.ResultStatusLost,
				},
			},
			hp1: 29,
			hp2: -22,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p1 := &FakePlayer{moves: tt.moves1}
			cp1 := match.CharacterPlayer{
				Character: game.NewCharacter(tt.char1),
				Player:    p1,
			}
			p2 := &FakePlayer{moves: tt.moves2}
			cp2 := match.CharacterPlayer{
				Character: game.NewCharacter(tt.char2),
				Player:    p2,
			}

			m := match.New(cp1, cp2, tt.invertedOrder)

			go m.Run(context.Background())

			reserr := <-m.Result()
			require.NoError(t, reserr.Err, "error")
			assert.Equal(t, tt.result, reserr.Res)

			assert.Equal(t, tt.hp1, cp1.Character.HP(), "hp1")
			assert.Equal(t, tt.hp2, cp2.Character.HP(), "hp2")

			assert.True(t, p1.hasEnd, "p1 end")
			assert.True(t, p2.hasEnd, "p2 end")
		})
	}
}
