package game

import (
	"log/slog"
	"net/http"

	"github.com/ShmaykhelDuo/battler/internal/app/game/ws"
	apimodel "github.com/ShmaykhelDuo/battler/internal/model/api"
	"github.com/ShmaykhelDuo/battler/internal/pkg/api"
	"github.com/ShmaykhelDuo/battler/internal/pkg/auth"
)

func (h *Handler) StartMatch(w http.ResponseWriter, r *http.Request) {
	s, err := auth.Session(r.Context())
	if err != nil {
		api.HandleError(w, apimodel.Error{
			Kind: apimodel.KindUnauthenticated,
		})
		return
	}

	c, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("websocket upgrade failed", "err", err)
		return
	}

	conn := ws.NewConn(s.UserID, c, h.matchSvc)
	err = conn.Handle(r.Context())
	if err != nil {
		slog.Error("match connection error", "err", err)
	}
}
