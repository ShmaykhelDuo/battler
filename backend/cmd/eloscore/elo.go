package main

import (
	"fmt"
	"log"
	"math"
	"math/rand/v2"
	"slices"

	"github.com/ShmaykhelDuo/battler/backend/internal/game"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/match"
)

type eloScoring struct {
	c1, c2  game.CharacterData
	players map[string]playerGenerator
	ratings map[string]int
}

func newEloScoring(c1, c2 game.CharacterData, players map[string]playerGenerator) *eloScoring {
	s := &eloScoring{
		c1:      c1,
		c2:      c2,
		players: players,
		ratings: make(map[string]int, len(players)),
	}

	for p := range players {
		s.ratings[p] = 1800
	}

	return s
}

func (s *eloScoring) run(rounds int) error {
	pairs := s.pairs()

	for i := range rounds {
		order := slices.Clone(pairs)
		rand.Shuffle(len(order), func(i, j int) {
			order[i], order[j] = order[j], order[i]
		})
		log.Printf("%v\n", order)

		for j, pair := range order {
			log.Printf("Round %d/%d: match %d/%d\n", i+1, rounds, j+1, len(order))

			err := s.match(pair[0], pair[1])
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *eloScoring) pairs() [][2]string {
	pairs := make([][2]string, 0, len(s.players)*(len(s.players)-1))

	for p1 := range s.players {
		for p2 := range s.players {
			if p1 == p2 {
				continue
			}

			pairs = append(pairs, [2]string{p1, p2})
		}
	}

	return pairs
}

func (s *eloScoring) match(p1, p2 string) error {
	player1 := s.players[p1].Player(s.c1.Desc)
	player2 := s.players[p2].Player(s.c2.Desc)
	c1 := game.NewCharacter(s.c1)
	c2 := game.NewCharacter(s.c2)

	res, err := match.Match(c1, c2, player1, player2)
	if err != nil {
		return fmt.Errorf("match: %w", err)
	}

	s.updateRatings(p1, p2, res)

	if res == 1 {
		log.Printf("%s won over %s.", p2, p1)
	} else if res == -1 {
		log.Printf("%s won over %s.", p1, p2)
	} else {
		log.Printf("%s drawed with %s.", p1, p2)
	}

	log.Printf("New ratings: %s=%d, %s=%d\n", p1, s.ratings[p1], p2, s.ratings[p2])

	return nil
}

func (s *eloScoring) updateRatings(p1, p2 string, res int) {
	r1 := s.ratings[p1]
	r2 := s.ratings[p2]

	q1 := math.Pow(10, float64(r1)/400)
	q2 := math.Pow(10, float64(r2)/400)

	e1 := q1 / (q1 + q2)
	e2 := q2 / (q1 + q2)

	s1 := 0.5 * float64(1-res)
	s2 := 0.5 * float64(1+res)

	k := 32.0
	s.ratings[p1] = r1 + int(math.Round(k*(s1-e1)))
	s.ratings[p2] = r2 + int(math.Round(k*(s2-e2)))
}
