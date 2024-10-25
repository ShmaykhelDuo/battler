package main

import (
	"context"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"strconv"

	"github.com/ShmaykhelDuo/battler/backend/internal/game"
	"golang.org/x/sync/errgroup"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	outputFile := flag.String("o", "", "output file name")
	charNumber := flag.Int("c", 0, "character number")
	oppNumber := flag.Int("opp", 0, "opponent number")
	bufSize := flag.Int("buf", 100, "channel buffer size")
	depth := flag.Int("d", 10, "depth")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to `file`")
	memprofile := flag.String("memprofile", "", "write memory profile to `file`")
	flag.Parse()

	if *outputFile == "" {
		return errors.New("output file name required")
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

	if *bufSize < 0 {
		return errors.New("invalid buffer size")
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

	f, err := os.Create(*outputFile)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}
	defer f.Close()

	data := make(chan []int, *bufSize)
	out := make(chan Out, *bufSize)

	eg, egCtx := errgroup.WithContext(ctx)

	d1 := make(chan []bool)
	d2 := make(chan []bool)

	eg.Go(func() error {
		err := writeData(egCtx, f, data)
		if err != nil {
			return fmt.Errorf("write: %w", err)
		}
		return nil
	})
	eg.Go(func() error {
		err := convertOut(egCtx, out, data)
		if err != nil {
			return fmt.Errorf("convert: %w", err)
		}
		return nil
	})
	eg.Go(func() error {
		err := runMiniMax(egCtx, data1, data2, *depth, true, out)
		close(d1)
		if err != nil {
			return fmt.Errorf("minimax first: %w", err)
		}
		return nil
	})
	eg.Go(func() error {
		err := runMiniMax(egCtx, data1, data2, *depth, false, out)
		close(d2)
		if err != nil {
			return fmt.Errorf("minimax second: %w", err)
		}
		return nil
	})
	eg.Go(func() error {
		<-d1
		<-d2
		close(out)
		return nil
	})

	err = eg.Wait()

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			return fmt.Errorf("could not create memory profile: %w", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			return fmt.Errorf("could not write memory profile: %w", err)
		}
	}

	return err
}

func writeData(ctx context.Context, out io.Writer, data <-chan []int) error {
	w := csv.NewWriter(out)
	defer func() {
		w.Flush()
		err := w.Error()
		if err != nil {
			log.Printf("flush: %v\n", err)
		}
	}()

	err := w.Write(headerRow())
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}

	for {
		select {
		case row, ok := <-data:
			if !ok {
				return nil
			}

			err := w.Write(toStringRow(row))
			if err != nil {
				return fmt.Errorf("write: %w", err)
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func headerRow() []string {
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
	return fields
}

func toStringRow(row []int) []string {
	res := make([]string, len(row))
	for i, v := range row {
		res[i] = strconv.Itoa(v)
	}
	return res
}

func runMiniMax(ctx context.Context, data1, data2 game.CharacterData, depth int, goingFirst bool, out chan<- Out) error {
	c1 := game.NewCharacter(data1)
	c2 := game.NewCharacter(data2)
	turnState := game.TurnState{
		TurnNum:      1,
		IsGoingFirst: true,
	}

	_, _, res, err := MiniMax(ctx, c1, c2, turnState, 1, depth, !goingFirst, nil)

	for _, o := range res {
		out <- o
	}
	return err
}

func convertOut(ctx context.Context, in <-chan Out, out chan<- []int) error {
	defer close(out)

	for {
		select {
		case rec, ok := <-in:
			if !ok {
				return nil
			}
			select {
			case out <- outToSlice(rec):
			case <-ctx.Done():
				return ctx.Err()
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func outToSlice(rec Out) []int {
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
	return intRec
}
