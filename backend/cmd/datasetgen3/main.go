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
	w, err := csvWriter("dataset-full-Ruby-vs-Milana-2.csv")
	if err != nil {
		log.Fatal(err)
	}

	fields := []string{
		"first",
	}
	for i := range 10 {
		fields = append(fields, fmt.Sprintf("skill%d", i+1))
	}
	for i := range 10 {
		fields = append(fields, fmt.Sprintf("oppskill%d", i+1))
	}
	fields = append(fields, "res")
	w.Write(fields)

	c := make(chan Out, 100)
	done := make(chan bool)

	go func() {
		n := 0
		cache := make(map[string]struct{})
		for {
			rec, ok := <-c
			if !ok {
				w.Flush()
				break
			}

			str := fmt.Sprintf("%#v", rec)
			if _, ok := cache[str]; ok {
				panic(str)
			} else {
				cache[str] = struct{}{}
			}

			n++
			if n%1000000 == 0 {
				fmt.Printf("%v\n", n)
			}
			intRec := make([]int, 22)
			if rec.First {
				intRec[0] = 1
			}
			ourAdd := 0
			oppAdd := 1
			if !rec.First {
				oppAdd = 0
				ourAdd = 1
			}
			for i := range 10 {
				if i*2+ourAdd < len(rec.PrevMoves) {
					intRec[1+i] = rec.PrevMoves[i*2+ourAdd]
				} else {
					intRec[1+i] = -1
				}

				if i*2+oppAdd < len(rec.PrevMoves) {
					intRec[11+i] = rec.PrevMoves[i*2+oppAdd]
				} else {
					intRec[11+i] = -1
				}
			}
			intRec[21] = rec.Strategy[0]
			strRec := make([]string, len(intRec))
			for i, n := range intRec {
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

	MiniMax(c1, c2, gameCtx, 1, 9, false, nil, c)

	c1 = game.NewCharacter(ruby.CharacterRuby)
	c2 = game.NewCharacter(milana.CharacterMilana)
	gameCtx = game.Context{
		TurnNum:      1,
		IsGoingFirst: true,
	}

	MiniMax(c1, c2, gameCtx, 1, 9, true, nil, c)

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
