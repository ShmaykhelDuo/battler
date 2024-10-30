package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"

	"github.com/ShmaykhelDuo/battler/internal/bot/minimax"
	"github.com/ShmaykhelDuo/battler/internal/bot/ml"
	"github.com/ShmaykhelDuo/battler/internal/bot/ml/formats"
	"github.com/ShmaykhelDuo/battler/internal/game"
	"github.com/ShmaykhelDuo/battler/internal/game/match"
	"golang.org/x/sync/errgroup"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	fullOutputFile := flag.String("outfull", "", "full state output file name")
	movesOutputFile := flag.String("outmoves", "", "moves output file name")
	charNumber := flag.Int("c", 0, "character number")
	oppNumber := flag.Int("opp", 0, "opponent number")
	depth := flag.Int("d", 10, "depth")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to `file`")
	memprofile := flag.String("memprofile", "", "write memory profile to `file`")
	flag.Parse()

	if *fullOutputFile == "" {
		return errors.New("full output file name required")
	}
	if *movesOutputFile == "" {
		return errors.New("moves output file name required")
	}
	if *charNumber == 0 || *oppNumber == 0 {
		return errors.New("character and opponent numbers required")
	}

	data1, ok := chars[*charNumber]
	if !ok {
		return errors.New("invalid character number")
	}
	data2, ok := chars[*oppNumber]
	if !ok {
		return errors.New("invalid opponent number")
	}

	if *depth <= 0 {
		return errors.New("invalid depth")
	}

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			return fmt.Errorf("could not create CPU profile: %w", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			return fmt.Errorf("could not start CPU profile: %w", err)
		}
		defer pprof.StopCPUProfile()
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	var res1, res2 []minimax.Entry

	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		var err error
		res1, err = runMiniMax(egCtx, data1, data2, *depth, true)
		if err != nil {
			return fmt.Errorf("minimax first: %w", err)
		}
		return nil
	})
	eg.Go(func() error {
		var err error
		res2, err = runMiniMax(egCtx, data1, data2, *depth, false)
		if err != nil {
			return fmt.Errorf("minimax second: %w", err)
		}
		return nil
	})

	err := eg.Wait()
	if err != nil {
		return err
	}

	res := append(res1, res2...)

	err = output(*fullOutputFile, res, formats.FullStateFormat{})
	if err != nil {
		return err
	}

	err = output(*movesOutputFile, res, formats.PrevMovesFormat{})
	if err != nil {
		return err
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			return fmt.Errorf("could not create memory profile: %w", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.Lookup("allocs").WriteTo(f, 0); err != nil {
			return fmt.Errorf("could not write memory profile: %w", err)
		}
	}

	return err
}

func runMiniMax(ctx context.Context, data1, data2 game.CharacterData, depth int, goingFirst bool) ([]minimax.Entry, error) {
	c1 := game.NewCharacter(data1)
	c2 := game.NewCharacter(data2)
	turnState := game.TurnState{
		TurnNum:      1,
		IsGoingFirst: true,
	}

	state := match.GameState{
		Character:  c1,
		Opponent:   c2,
		TurnState:  turnState,
		SkillsLeft: 1,
		SkillLog:   make(match.SkillLog),
		PlayerTurn: true,
		AsOpp:      !goingFirst,
	}

	res, err := minimax.MemOptConcurrentRunner.MiniMax(ctx, state, depth)
	if err != nil {
		return nil, err
	}

	return res.Entries, nil
}

func output(filename string, res []minimax.Entry, format ml.Format) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}
	defer f.Close()

	err = ml.ExportDataset(f, res, format)
	if err != nil {
		return fmt.Errorf("export: %w", err)
	}

	return nil
}
