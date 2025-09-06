package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/metgag/koda-weekly10/internals/models"
	"github.com/metgag/koda-weekly10/internals/repositories"
	"github.com/metgag/koda-weekly10/pkg"
)

type AuthHandler struct {
	ar *repositories.AuthRepository
}

func NewAuthHandler(ar *repositories.AuthRepository) *AuthHandler {
	return &AuthHandler{ar: ar}
}

func newAuthResponse(err string, success bool, result string) models.AuthResponse {
	return models.AuthResponse{Error: err, Success: success, Result: result}
}

func (a *AuthHandler) AddUser(ctx *gin.Context) {
	var body = models.Register{}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusInternalServerError, newAuthResponse(
			"server unable to bind request", false, "",
		))
		return
	}

	p := pkg.NewHashParams()
	p.UseRecommended()
	encodedHash, err := p.GenerateFromPassword(body.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newAuthResponse(
			"server error while encoding password", false, "",
		))
		return
	}

	id, err := a.ar.AddNewUser(ctx.Request.Context(), body.Email, encodedHash)
	if err != nil {
		ctx.JSON(http.StatusConflict, newAuthResponse(
			"duplicate email addresses", false, "",
		))
		return
	}

	ctx.JSON(http.StatusOK, newAuthResponse(
		"", true, fmt.Sprintf("register succesfully w/ ID: %d", id),
	))
}

func (a *AuthHandler) Login(ctx *gin.Context) {
	var body = models.Login{}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusInternalServerError, newAuthResponse(
			"server unable to bind request", false, "",
		))
		return
	}

	user, err := a.ar.GetUser(ctx.Request.Context(), body.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newAuthResponse(
			"server unable to get user", false, "",
		))
		return
	}

	hc := pkg.NewHashParams()
	isMatch, err := hc.ComparePasswordAndHash(body.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newAuthResponse(
			"server unable to compare password", false, "",
		))
		return
	}
	if !isMatch {
		ctx.JSON(http.StatusBadRequest, newAuthResponse(
			"invalid email or password", false, "",
		))
		return
	}

	ctx.JSON(http.StatusOK, newAuthResponse(
		"", true, fmt.Sprintf("logged in as UID: %d", user.ID),
	))
}
