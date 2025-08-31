package profile

import (
	"database/sql"
	"net/http"
	"questionbasket/frame"
	"questionbasket/models"
)

type ProfileAPI struct {
	profileModel models.ProfileModel
}

func NewProfileAPI() ProfileAPI {
	pAPI := ProfileAPI{
		profileModel: models.ProfileModel{},
	}
	frame.DatabaseBind(&pAPI.profileModel)
	return pAPI
}

func (pAPI *ProfileAPI) GetProfile(ctx frame.APIContext) {
	profile, err := pAPI.profileModel.GetProfile()
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.ReturnError("not_found", "Profile data not found", http.StatusNotFound)
		} else {
			ctx.ReturnError("db_error", "Failed to get profile data", http.StatusInternalServerError)
		}
		return
	}

	resp := ProfileResponse{
		Name:        profile.Name,
		Description: profile.Description,
		BasedOn:     profile.BasedOn,
	}
	ctx.ReturnJSON(resp)
}
