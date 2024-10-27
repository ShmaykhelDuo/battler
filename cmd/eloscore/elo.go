package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand/v2"
	"slices"

	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/match"
	"golang.org/x/sync/errgroup"
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

		eg, _ := errgroup.WithContext(context.Background())

		res := make([]match.Result, len(order))
		for j, pair := range order {
			eg.Go(func() error {
				var err error
				res[j], err = s.match(pair[0], pair[1])
				if err != nil {
					return err
				}

				p1 := pair[0]
				p2 := pair[1]

				switch res[j] {
				case match.ResultWonFirst:
					log.Printf("%s as %s won over %s as %s.", p1, s.c1.Desc.Name, p2, s.c2.Desc.Name)
				case match.ResultWonSecond:
					log.Printf("%s as %s won over %s as %s.", p2, s.c2.Desc.Name, p1, s.c1.Desc.Name)
				default:
					log.Printf("%s as %s drawed with %s as %s.", p1, s.c1.Desc.Name, p2, s.c2.Desc.Name)
				}

				return nil
			})
		}

		err := eg.Wait()
		if err != nil {
			return err
		}

		log.Printf("Round %d/%d SUMMARY\n\n", i+1, rounds)

		for j, pair := range order {
			log.Printf("Round %d/%d: match %d/%d\n", i+1, rounds, j+1, len(order))

			p1 := pair[0]
			p2 := pair[1]

			s.updateRatings(p1, p2, res[j])

			switch res[j] {
			case match.ResultWonFirst:
				log.Printf("%s as %s won over %s as %s.", p1, s.c1.Desc.Name, p2, s.c2.Desc.Name)
			case match.ResultWonSecond:
				log.Printf("%s as %s won over %s as %s.", p2, s.c2.Desc.Name, p1, s.c1.Desc.Name)
			default:
				log.Printf("%s as %s drawed with %s as %s.", p1, s.c1.Desc.Name, p2, s.c2.Desc.Name)
			}

			log.Printf("New ratings: %s=%d, %s=%d\n", p1, s.ratings[p1], p2, s.ratings[p2])
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

func (s *eloScoring) match(p1, p2 string) (match.Result, error) {
	cp1 := match.CharacterPlayer{
		Character: game.NewCharacter(s.c1),
		Player:    s.players[p1].Player(s.c1.Desc),
	}
	cp2 := match.CharacterPlayer{
		Character: game.NewCharacter(s.c2),
		Player:    s.players[p2].Player(s.c2.Desc),
	}

	invertedOrder := rand.IntN(2) == 1
	m := match.New(cp1, cp2, invertedOrder)

	err := m.Run(context.Background())
	if err != nil {
		return 0, fmt.Errorf("match: %w", err)
	}

	res, err := m.Result()
	if err != nil {
		return 0, fmt.Errorf("match result: %w", err)
	}

	return res, nil
}

func (s *eloScoring) updateRatings(p1, p2 string, res match.Result) {
	r1 := s.ratings[p1]
	r2 := s.ratings[p2]

	q1 := math.Pow(10, float64(r1)/400)
	q2 := math.Pow(10, float64(r2)/400)

	e1 := q1 / (q1 + q2)
	e2 := q2 / (q1 + q2)

	var s1, s2 float64
	switch res {
	case match.ResultWonFirst:
		s1 = 1.0
		s2 = 0.0
	case match.ResultWonSecond:
		s1 = 0.0
		s2 = 1.0
	default:
		s1 = 0.5
		s2 = 0.5
	}

	k := 32.0
	s.ratings[p1] = r1 + int(math.Round(k*(s1-e1)))
	s.ratings[p2] = r2 + int(math.Round(k*(s2-e2)))
}
