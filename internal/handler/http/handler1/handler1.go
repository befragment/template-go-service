package handler1

import (
	"encoding/json"
	"errors"
	"net/http"

	jwtctx "github.com/befragment/template-go/internal/domain/lib/auth/jwtprovider"
)

type Handler1 struct {
	// actually more usecases can be defined here
	uc1 usecase1
}

func NewHandler1(uc1 usecase1) *Handler1 {
	return &Handler1{uc1: uc1}
}

func (h *Handler1) SomeMethodHandler1(w http.ResponseWriter, r *http.Request) {
	 // always get context from http request and pass it to usecases
	ctx := r.Context()

	// example of getting principal
	// then its going to be passed to usecases
	p, ok := jwtctx.PrincipalFromContext(ctx) 
	if !ok || p == nil {
		utils.RespondUnauthorized(w, "unauthorized")
		return
	}

	integer, err := h.uc1.SomeMethodUC1(ctx, id)
	
	if err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			utils.RespondForbidden(w, err.Error())
			return
		}
		utils.RespondInternalServerError(w, err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, integer)
}
