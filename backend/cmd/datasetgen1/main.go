package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"strconv"

	"github.com/ShmaykhelDuo/battler/backend/internal/game"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/bot"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/bot/minimax"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/characters/storyteller"
	"github.com/ShmaykhelDuo/battler/backend/internal/game/match"
)

func main() {
	w, err := csvWriter("dataset-1.csv")
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan []int, 100)
	done := make(chan bool)

	go func() {
		for {
			rec, ok := <-c
			if !ok {
				w.Flush()
				break
			}
			strRec := make([]string, len(rec))
			for i, n := range rec {
				strRec[i] = strconv.Itoa(n)
			}

			w.Write(strRec)
		}
		done <- true
	}()

	for i := range 1000 {
		if i%100 == 0 {
			fmt.Printf("%d%%\n", i/100)
		}
		b1 := NewBot(4, c)

		botI := rand.IntN(4)
		var b2 match.Player
		switch botI {
		case 0:
			b2 = &bot.RandomBot{}
		default:
			p := rand.Float64()
			b2 = bot.NewRandomWrapperBot(minimax.NewBot(botI), p)
		}

		c1 := game.NewCharacter(storyteller.CharacterStoryteller)
		c2 := getRandomChar(storyteller.CharacterStoryteller)

		_, err := match.Match(c1, c2, b1, b2)
		if err != nil {
			log.Printf("match: %v\n", err)
			return
		}
		// switch res {
		// case 1:
		// 	fmt.Println("Lost")
		// case -1:
		// 	fmt.Println("Won")
		// default:
		// 	fmt.Println("Draw")
		// }
		// fmt.Printf("%d : %d\n", c1.HP(), c2.HP())
	}

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
