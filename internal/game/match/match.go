package match

import (
	"context"
	"errors"
	"maps"

	"github.com/ShmaykhelDuo/battler/internal/game"
)

var ErrMatchNotEnded = errors.New("match has not ended")

type CharacterPlayer struct {
	Character *game.Character
	Player    Player
}

type Result struct {
	Player1 ResultPlayer
	Player2 ResultPlayer
}

type ResultPlayer struct {
	Status     ResultStatus
	HasGivenUp bool
}

type ResultStatus int

const (
	ResultStatusUnknown ResultStatus = 0
	ResultStatusWon     ResultStatus = 1
	ResultStatusLost    ResultStatus = 2
	ResultStatusDraw    ResultStatus = 3
)

type ResultError struct {
	Res Result
	Err error
}

// Match is a game between two characters.
type Match struct {
	p1, p2        CharacterPlayer
	invertedOrder bool
	skillLog      SkillLog
	result        chan ResultError
}

// New creates a new match.
func New(p1, p2 CharacterPlayer, invertedOrder bool) *Match {
	if invertedOrder {
		p1, p2 = p2, p1
	}

	return &Match{
		p1:            p1,
		p2:            p2,
		invertedOrder: invertedOrder,
		skillLog:      make(SkillLog),
		result:        make(chan ResultError),
	}
}

// Run runs the match.
func (m *Match) Run(ctx context.Context) {
	done := make(chan ResultError)

	turnsCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		res, err := m.runTurns(turnsCtx)
		done <- ResultError{
			Res: res,
			Err: err,
		}
		close(done)
	}()

	select {
	case res := <-done:
		m.result <- res
	case <-m.p1.Player.GivenUp():
		m.result <- ResultError{
			Res: m.rightOrderResult(
				ResultPlayer{
					Status:     ResultStatusLost,
					HasGivenUp: true,
				},
				ResultPlayer{
					Status: ResultStatusWon,
				},
			),
		}
	case <-m.p2.Player.GivenUp():
		m.result <- ResultError{
			Res: m.rightOrderResult(
				ResultPlayer{
					Status: ResultStatusWon,
				},
				ResultPlayer{
					Status:     ResultStatusLost,
					HasGivenUp: true,
				},
			),
		}
	}
}

func (m *Match) Result() <-chan ResultError {
	return m.result
}

func (m *Match) rightOrderResult(result1 ResultPlayer, result2 ResultPlayer) Result {
	if m.invertedOrder {
		return Result{
			Player1: result2,
			Player2: result1,
		}
	}

	return Result{
		Player1: result1,
		Player2: result2,
	}
}

func (m *Match) runTurns(ctx context.Context) (Result, error) {
	var turnState game.TurnState
	for turnState = game.StartTurnState(); turnState.TurnNum <= game.MaxTurnNumber; turnState = turnState.Next() {
		end, err := m.runTurn(ctx, turnState)
		if err != nil {
			return Result{}, err
		}
		if end {
			break
		}
	}

	err := m.sendState(ctx, m.p1, m.p2, turnState.WithTurnEnd(), 0, false, false)
	if err != nil {
		return Result{}, err
	}
	err = m.sendState(ctx, m.p2, m.p1, turnState.WithTurnEnd(), 0, false, false)
	if err != nil {
		return Result{}, err
	}
	err = m.p1.Player.SendEnd(ctx)
	if err != nil {
		return Result{}, err
	}
	err = m.p2.Player.SendEnd(ctx)
	if err != nil {
		return Result{}, err
	}

	var res1, res2 ResultStatus
	if m.p1.Character.HP() > m.p2.Character.HP() {
		res1 = ResultStatusWon
		res2 = ResultStatusLost
	} else if m.p1.Character.HP() < m.p2.Character.HP() {
		res1 = ResultStatusLost
		res2 = ResultStatusWon
	} else {
		res1 = ResultStatusDraw
		res2 = ResultStatusDraw
	}

	return m.rightOrderResult(
		ResultPlayer{
			Status: res1,
		},
		ResultPlayer{
			Status: res2,
		},
	), nil
}

func (m *Match) runTurn(ctx context.Context, turnState game.TurnState) (end bool, err error) {
	var c, opp CharacterPlayer
	if turnState.IsGoingFirst {
		c = m.p1
		opp = m.p2
	} else {
		c = m.p2
		opp = m.p1
	}

	asOpp := c.Character.IsControlledByOpp()

	var control, observer CharacterPlayer
	if asOpp {
		control = opp
		observer = c
	} else {
		control = c
		observer = opp
	}

	for skillsLeft := c.Character.SkillsPerTurn(); skillsLeft > 0; skillsLeft-- {
		err = m.sendState(ctx, control, observer, turnState, skillsLeft, true, asOpp)
		if err != nil {
			return true, err
		}
		err = m.sendState(ctx, observer, control, turnState, skillsLeft, false, asOpp)
		if err != nil {
			return true, err
		}

		for {
			i, err := control.Player.RequestSkill(ctx)
			if err != nil {
				return true, err
			}

			// log.Printf("Player %s has selected skill %d\n", c.Desc().Name, i)

			err = c.Character.Skills()[i].Use(c.Character, opp.Character, turnState)
			if err == nil {
				m.skillLog[turnState] = append(m.skillLog[turnState], i)
				break
			}

			err = control.Player.SendError(ctx, err)
			if err != nil {
				return true, err
			}
		}

		if m.isEnd() {
			return true, nil
		}
	}

	endState := turnState.WithTurnEnd()
	c.Character.OnTurnEnd(opp.Character, endState)
	opp.Character.OnTurnEnd(c.Character, endState)
	return m.isEnd(), nil
}

func (m *Match) sendState(ctx context.Context, c, opp CharacterPlayer, turnState game.TurnState, skillsLeft int, playerTurn bool, asOpp bool) error {
	state := GameState{
		Character:  c.Character.Clone(),
		Opponent:   opp.Character.Clone(),
		TurnState:  turnState,
		SkillsLeft: skillsLeft,
		SkillLog:   maps.Clone(m.skillLog),
		PlayerTurn: playerTurn,
		AsOpp:      asOpp,
	}
	return c.Player.SendState(ctx, state)
}

func (m *Match) isEnd() bool {
	return m.p1.Character.HP() <= 0 || m.p2.Character.HP() <= 0
}
