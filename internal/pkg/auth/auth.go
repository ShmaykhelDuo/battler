package auth

import (
	"context"

	model "github.com/ShmaykhelDuo/battler/internal/model/auth"
)

type key int

var sessionKey key

func ContextWithSession(ctx context.Context, session model.Session) context.Context {
	return context.WithValue(ctx, sessionKey, session)
}

func Session(ctx context.Context) (model.Session, error) {
	session, ok := ctx.Value(sessionKey).(model.Session)
	if !ok {
		return model.Session{}, model.ErrNoSession
	}
	return session, nil
}
