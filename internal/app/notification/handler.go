package notification

import (
	"log/slog"
	"net/http"

	apimodel "github.com/ShmaykhelDuo/battler/internal/model/api"
	"github.com/ShmaykhelDuo/battler/internal/pkg/api"
	"github.com/ShmaykhelDuo/battler/internal/pkg/auth"
	service "github.com/ShmaykhelDuo/battler/internal/service/notification"
	"github.com/gorilla/websocket"
)

type Handler struct {
	svc      *service.Service
	upgrader websocket.Upgrader
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) ReceiveNotifications(w http.ResponseWriter, r *http.Request) {
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

	conn := NewConn(s.UserID, c, h.svc)
	err = conn.Handle(r.Context())
	if err != nil {
		slog.Error("notification connection error", "err", err)
	}
}
