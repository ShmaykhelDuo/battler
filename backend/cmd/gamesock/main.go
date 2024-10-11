package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/ShmaykhelDuo/battler/backend/internal/bot/ml"
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

type stateMsg struct {
	State  []int `json:"state"`
	End    bool  `json:"end"`
	Reward int   `json:"reward"`
}

type actionMsg struct {
	Action int `json:"action"`
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	encoder := json.NewEncoder(conn)
	decoder := json.NewDecoder(conn)

	out := make(chan PlayerState)
	in := make(chan int)

	p := Player{
		Out: out,
		In:  in,
	}
	bot := Bot()

	go Game(p, bot)

	totalReward := 0

	for {
		state := <-out

		reward := state.State.Character.HP() - state.State.Opponent.HP()
		if state.End {
			if state.Win {
				reward += 100
			} else {
				reward -= 100
			}
		}
		sMsg := stateMsg{
			State:  ml.NewState(state.State).ToSlice(),
			End:    state.End,
			Reward: reward - totalReward,
		}
		totalReward = reward
		fmt.Printf("Sending %#v\n", sMsg)
		err := encoder.Encode(sMsg)
		if err != nil {
			log.Printf("send: %v", err)
			return
		}

		if state.End {
			close(in)
			break
		}

		var msg actionMsg
		err = decoder.Decode(&msg)
		if err != nil {
			log.Printf("recv: %v", err)
			return
		}
		fmt.Printf("Got action %#v\n", msg)
		in <- msg.Action
	}
}
