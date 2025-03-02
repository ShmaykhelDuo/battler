package friends

import (
	"net/http"

	"github.com/ShmaykhelDuo/battler/internal/model/social"
	"github.com/ShmaykhelDuo/battler/internal/pkg/api"
	"github.com/google/uuid"
)

type Profile struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

func (h *Handler) Friends(w http.ResponseWriter, r *http.Request) {
	profiles, err := h.s.Friends(r.Context())
	if err != nil {
		api.HandleError(w, err)
		return
	}

	api.WriteJSONResponse(w, http.StatusOK, profilesToDto(profiles))
}

func (h *Handler) IncomingFriendRequests(w http.ResponseWriter, r *http.Request) {
	profiles, err := h.s.IncomingFriendRequests(r.Context())
	if err != nil {
		api.HandleError(w, err)
		return
	}

	api.WriteJSONResponse(w, http.StatusOK, profilesToDto(profiles))
}

func (h *Handler) OutgoingFriendRequests(w http.ResponseWriter, r *http.Request) {
	profiles, err := h.s.OutgoingFriendRequests(r.Context())
	if err != nil {
		api.HandleError(w, err)
		return
	}

	api.WriteJSONResponse(w, http.StatusOK, profilesToDto(profiles))
}

func profilesToDto(profiles []social.Profile) []Profile {
	res := make([]Profile, len(profiles))
	for i, p := range profiles {
		res[i] = Profile{
			ID:       p.ID,
			Username: p.Username,
		}
	}
	return res
}
