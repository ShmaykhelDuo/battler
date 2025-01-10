package gamestate

import (
	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/match"
)

type Character struct {
	Number              int            `json:"number"`
	SkillAvailabilities []bool         `json:"skill_availabilities"`
	HP                  int            `json:"hp"`
	MaxHP               int            `json:"max_hp"`
	Effects             []Effect       `json:"effects"`
	Defences            map[Colour]int `json:"defences"`
}

func NewChacacter(c *game.Character, opp *game.Character, turnState game.TurnState) Character {
	effs := c.Effects()
	defs := c.Defences()

	res := Character{
		Number:              c.Desc().Number,
		SkillAvailabilities: make([]bool, 4),
		HP:                  c.HP(),
		MaxHP:               c.MaxHP(),
		Effects:             make([]Effect, len(effs)),
		Defences:            make(map[Colour]int, len(defs)),
	}

	for i, s := range c.Skills() {
		res.SkillAvailabilities[i] = s.IsAvailable(c, opp, turnState)
	}

	for i, e := range effs {
		res.Effects[i] = NewEffect(e, turnState)
	}

	for c, v := range defs {
		res.Defences[NewColour(c)] = v
	}

	return res
}

type TurnState struct {
	Number  int  `json:"number"`
	IsFirst bool `json:"first"`
	End     bool `json:"end"`
}

func NewTurnState(s game.TurnState) TurnState {
	return TurnState{
		Number:  s.TurnNum,
		IsFirst: s.IsGoingFirst,
		End:     s.IsTurnEnd,
	}
}

type GameState struct {
	Character    Character `json:"character"`
	Opponent     Character `json:"opponent"`
	Turn         TurnState `json:"turn"`
	SkillsLeft   int       `json:"skills_left"`
	SkillLog     SkillLog  `json:"skill_log"`
	IsPlayerTurn bool      `json:"player_turn"`
	AsOpp        bool      `json:"as_opp"`
}

func NewGameState(state match.GameState) GameState {
	return GameState{
		Character:    NewChacacter(state.Character, state.Opponent, state.TurnState),
		Opponent:     NewChacacter(state.Opponent, state.Character, state.TurnState),
		Turn:         NewTurnState(state.TurnState),
		SkillsLeft:   state.SkillsLeft,
		SkillLog:     NewSkillLog(state.SkillLog, state.TurnState),
		IsPlayerTurn: state.PlayerTurn,
		AsOpp:        state.AsOpp,
	}
}
