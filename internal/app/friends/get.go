package friends

import (
	"net/http"

	apimodel "github.com/ShmaykhelDuo/battler/internal/model/api"
	"github.com/ShmaykhelDuo/battler/internal/model/social"
	"github.com/ShmaykhelDuo/battler/internal/pkg/api"
	"github.com/google/uuid"
)

type Profile struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

type ProfileFriendshipStatus struct {
	ID               uuid.UUID `json:"id"`
	Username         string    `json:"username"`
	FriendshipStatus int       `json:"friendshipstatus"`
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

func (h *Handler) FriendshipStatus(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		api.HandleError(w, apimodel.Error{
			Kind:    apimodel.KindNotFound,
			Message: "invalid user id",
		})
		return
	}

	profileStatus, err := h.s.FriendshipStatus(r.Context(), id)
	if err != nil {
		api.HandleError(w, err)
		return
	}

	res := ProfileFriendshipStatus{
		ID:               profileStatus.ID,
		Username:         profileStatus.Username,
		FriendshipStatus: int(profileStatus.FriendshipStatus),
	}
	api.WriteJSONResponse(w, http.StatusOK, res)
}
