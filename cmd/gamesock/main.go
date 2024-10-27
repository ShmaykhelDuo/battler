package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/bot/minimax"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/milana"
	"github.com/ShmaykhelDuo/battler/internal/game/characters/ruby"
	"github.com/ShmaykhelDuo/battler/internal/game/match"
)

func main() {
	sock, err := net.Listen("unix", "/tmp/test.sock")
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	// Cleanup the sockfile.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Remove("/tmp/test.sock")
		os.Exit(1)
	}()

	fmt.Println("Listening on socket..")

	for {
		fmt.Println("Waiting for connection..")
		conn, err := sock.Accept()
		if err != nil {
			log.Fatalf("accept: %v", err)
		}

		fmt.Println("Got connection")

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	learner := NewDQLLearnerBot(conn)
	bot := minimax.NewBot(4)
	// bot := &bot.RandomBot{}

	// c1, c2 := getRandomPair()
	c1 := game.NewCharacter(milana.CharacterMilana)
	c2 := game.NewCharacter(ruby.CharacterRuby)

	p1 := match.CharacterPlayer{
		Character: c1,
		Player:    learner,
	}
	p2 := match.CharacterPlayer{
		Character: c2,
		Player:    bot,
	}

	invertedOrder := rand.IntN(2) == 1
	m := match.New(p1, p2, invertedOrder)

	err := m.Run()
	if err != nil {
		log.Printf("match: %v\n", err)
		return
	}

	res, err := m.Result()
	if err != nil {
		log.Printf("match result: %v\n", err)
		return
	}

	switch res {
	case match.ResultWonSecond:
		fmt.Println("Lost")
	case match.ResultWonFirst:
		fmt.Println("Won")
	default:
		fmt.Println("Draw")
	}
}
