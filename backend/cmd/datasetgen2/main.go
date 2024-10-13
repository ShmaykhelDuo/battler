package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/ShmaykhelDuo/battler/backend/internal/game"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/characters/milana"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/characters/ruby"
)

func main() {
	w, err := csvWriter("dataset-full-Ruby-vs-Milana.csv")
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan []int, 100)
	done := make(chan bool)

	go func() {
		n := 0
		for {
			rec, ok := <-c
			if !ok {
				w.Flush()
				break
			}
			n++
			if n%1000000 == 0 {
				fmt.Printf("%v\n", n)
			}
			strRec := make([]string, len(rec))
			for i, n := range rec {
				strRec[i] = strconv.Itoa(n)
			}

			w.Write(strRec)
		}
		done <- true
	}()

	c1 := game.NewCharacter(ruby.CharacterRuby)
	c2 := game.NewCharacter(milana.CharacterMilana)
	gameCtx := game.Context{
		TurnNum:      1,
		IsGoingFirst: true,
	}

	MiniMax(c1, c2, gameCtx, 1, 8, false, c)

	c1 = game.NewCharacter(ruby.CharacterRuby)
	c2 = game.NewCharacter(milana.CharacterMilana)
	gameCtx = game.Context{
		TurnNum:      1,
		IsGoingFirst: true,
	}

	MiniMax(c1, c2, gameCtx, 1, 8, true, c)

	close(c)
	<-done
}

func csvWriter(filename string) (*csv.Writer, error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	return csv.NewWriter(f), nil
}
