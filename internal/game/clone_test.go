package game_test

import (
	"testing"

	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/euphoria"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/storyteller"
)

func BenchmarkCloneCharacter(b *testing.B) {
	c := game.NewCharacter(storyteller.CharacterStoryteller)
	c.AddEffect(storyteller.NewEffectCannotUse(game.TurnState{}, game.ColourNone))
	c.AddEffect(storyteller.NewEffectControlled(game.TurnState{}))
	c.AddEffect(euphoria.NewEffectEuphoricSource(13))

	for i := 0; i < b.N; i++ {
		game.CloneCharacter(c)
	}
}
