package match

import "context"

type Player interface {
	SendState(ctx context.Context, state GameState) error
	SendError(ctx context.Context, err error) error
	SendEnd(ctx context.Context) error
	RequestSkill(ctx context.Context) (int, error)
	GivenUp() <-chan any
}
