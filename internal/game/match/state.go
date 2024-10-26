package match

import (
	"slices"

	"github.com/ShmaykhelDuo/battler/internal/game"
)

type SkillLogItem struct {
	Character  *game.Character
	SkillIndex int
}

type SkillLog struct {
	items []SkillLogItem
}

func NewSkillLog() *SkillLog {
	return &SkillLog{}
}

func (l *SkillLog) Add(c *game.Character, i int) {
	l.items = append(l.items, SkillLogItem{Character: c, SkillIndex: i})
}

func (l *SkillLog) Items() []SkillLogItem {
	return slices.Clone(l.items)
}

type GameState struct {
	Character  *game.Character
	Opponent   *game.Character
	TurnState  game.TurnState
	SkillLog   []SkillLogItem
	PlayerTurn bool
	AsOpp      bool
}
