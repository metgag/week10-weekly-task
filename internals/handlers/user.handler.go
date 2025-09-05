package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/metgag/koda-weekly10/internals/models"
	"github.com/metgag/koda-weekly10/internals/repositories"
)

type UserHandler struct {
	ur *repositories.UserRepository
}

func NewUserHandler(ur *repositories.UserRepository) *UserHandler {
	return &UserHandler{ur: ur}
}

func newUserinfResponse(res models.UserInf, success bool, err string) models.UserinfResponse {
	return models.UserinfResponse{Result: res, Success: success, Error: err}
}

func (u *UserHandler) HandleUserinf(ctx *gin.Context) {
	idParam, err := strconv.Atoi(ctx.Param("uid"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newUserinfResponse(
			models.UserInf{}, false, "invalid uid input",
		))
		return
	}

	userinf, err := u.ur.GetUserinf(ctx.Request.Context(), idParam)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newUserinfResponse(
			models.UserInf{}, false, fmt.Sprintf("server unable to reach user w/ ID %d", idParam),
		))
		return
	}

	ctx.JSON(http.StatusOK, newUserinfResponse(
		userinf, true, "",
	))
}

func newUpdateResponse(res string, success bool, err string) models.UpdateResponse {
	return models.UpdateResponse{Result: res, Success: success, Error: err}
}

func (u *UserHandler) HandleUpdateUserInf(ctx *gin.Context) {
	idParam, err := strconv.Atoi(ctx.Param("uid"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newUpdateResponse(
			"", false, "invalid uid input",
		))
		return
	}

	var newUserInf models.NewInf
	if err := ctx.ShouldBindJSON(&newUserInf); err != nil {
		ctx.JSON(http.StatusInternalServerError, newUpdateResponse(
			"", false, "server unable to bind input",
		))
		return
	}

	ctag, err := u.ur.UpdateUserinf(newUserInf, ctx, idParam)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newUpdateResponse(
			"", false, "server unable to update request",
		))
		return
	}
	if ctag.RowsAffected() == 0 {
		if _, err := u.ur.InitUpdateUserinf(newUserInf, ctx, idParam); err == nil {
			ctx.JSON(http.StatusCreated, newUpdateResponse(
				fmt.Sprintf("user update created w/ ID %d", idParam), true, "",
			))
			return
		}

		ctx.JSON(http.StatusBadRequest, newUpdateResponse(
			"", false, fmt.Sprintf("there is no user w/ ID %d", idParam),
		))
		return
	}

	ctx.JSON(http.StatusOK, newUpdateResponse(
		fmt.Sprintf("updated user w/ ID %d", idParam), true, "",
	))
}
