package friends

import (
	"net/http"

	apimodel "github.com/ShmaykhelDuo/battler/internal/model/api"
	"github.com/ShmaykhelDuo/battler/internal/pkg/api"
	"github.com/google/uuid"
)

func (h *Handler) CreateFriendLink(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		api.HandleError(w, apimodel.Error{
			Kind:    apimodel.KindNotFound,
			Message: "invalid user id",
		})
		return
	}

	profile, err := h.s.CreateFriendLink(r.Context(), id)
	if err != nil {
		api.HandleError(w, err)
		return
	}

	res := Profile{
		ID:       profile.ID,
		Username: profile.Username,
	}

	api.WriteJSONResponse(w, http.StatusCreated, res)
}

func (h *Handler) RemoveFriendLink(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		api.HandleError(w, apimodel.Error{
			Kind:    apimodel.KindNotFound,
			Message: "invalid user id",
		})
		return
	}

	err = h.s.RemoveFriendLink(r.Context(), id)
	if err != nil {
		api.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
