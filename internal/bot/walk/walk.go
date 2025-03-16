package walk

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type Move[S any] struct {
	Action int
	State  S
}

// State is a walk node state.
type State[S any, R any] interface {
	// IsTerminal reports whether the node is terminal.
	IsTerminal() bool

	// IsParallel reports whether the walk needs to be done in parallel.
	IsParallel() bool

	// Value returns result of the node.
	Value() R

	// Children returns a slice of children nodes.
	Children() ([]Move[S], error)
}

type Params[S any, V any] interface {
	// InitVars returns an initialized vars.
	InitVars(state S) V
}

type Vars[S any, P any, V any, R any] interface {
	// Params returns params for nested walk call.
	Params() P

	// Accumulate updates vars based on nested call result.
	Accumulate(action int, res R) (update V, cutoff bool)

	// Result returns accumulated result.
	Result() R
}

type Runner[S State[S, R], P Params[S, V], V Vars[S, P, V, R], R any] struct {
}

func (r Runner[S, P, V, R]) Walk(ctx context.Context, state S, params P) (res R, err error) {
	if state.IsTerminal() {
		return state.Value(), nil
	}

	vars := params.InitVars(state)

	children, err := state.Children()
	if err != nil {
		return res, err
	}

	if state.IsParallel() {
		vars, err = r.walkParallel(ctx, children, vars)
		if err != nil {
			return res, err
		}
	} else {
		vars, err = r.walkSequential(ctx, children, vars)
		if err != nil {
			return res, err
		}
	}

	return vars.Result(), nil
}

func (r Runner[S, P, V, R]) walkSequential(ctx context.Context, moves []Move[S], vars V) (V, error) {
	for _, move := range moves {
		res, err := r.Walk(ctx, move.State, vars.Params())
		if err != nil {
			return vars, err
		}

		var cutoff bool
		vars, cutoff = vars.Accumulate(move.Action, res)
		if cutoff {
			break
		}
	}

	return vars, nil
}

func (r Runner[S, P, V, R]) walkParallel(ctx context.Context, moves []Move[S], vars V) (V, error) {
	eg, ctx := errgroup.WithContext(ctx)

	var results []R
	for _, move := range moves {
		eg.Go(func() error {
			res, err := r.Walk(ctx, move.State, vars.Params())
			if err != nil {
				return err
			}

			results = append(results, res)
			return nil
		})
	}

	err := eg.Wait()
	if err != nil {
		return vars, err
	}

	for i, res := range results {
		var cutoff bool
		vars, cutoff = vars.Accumulate(moves[i].Action, res)
		if cutoff {
			break
		}
	}

	return vars, nil
}
