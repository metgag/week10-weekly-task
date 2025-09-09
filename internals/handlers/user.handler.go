package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/metgag/koda-weekly10/internals/models"
	"github.com/metgag/koda-weekly10/internals/repositories"
	"github.com/metgag/koda-weekly10/pkg"
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

// HandleGetUserProfile godoc
//
//	@Summary	get user profile info based from ID
//	@Tags		users
//	@Accept		json
//	@Produce	json
//	@Param		Authorization	header		string	true	"Bearer token"
//	@Success	200				{object}	models.UserinfResponse
//	@Router		/users/ [get]
func (u *UserHandler) HandleUserinf(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	user, _ := claims.(pkg.Claims)

	userinf, err := u.ur.GetUserinf(ctx.Request.Context(), user.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newUserinfResponse(
			models.UserInf{}, false, fmt.Sprintf("server unable to reach user w/ ID %d", user.UserID),
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

// HandleUpdateUserProfile godoc
//
//	@Summary	update user profile info based from ID
//	@Tags		users
//	@Accept		json
//	@Produce	json
//	@Param		Authorization	header		string			true	"Bearer token"
//	@Param		request			body		models.NewInf	true	"profile update body json content"
//	@Success	200				{object}	models.UpdateResponse
//	@Router		/users/ [patch]
func (u *UserHandler) HandleUpdateUserInf(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	user, _ := claims.(pkg.Claims)

	var newUserInf models.NewInf
	if err := ctx.ShouldBindJSON(&newUserInf); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, newUpdateResponse(
			"", false, "server unable to bind input",
		))
		return
	}

	ctag, err := u.ur.UpdateUserinf(newUserInf, ctx, user.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newUpdateResponse(
			"", false, "server unable to update request",
		))
		return
	}
	if ctag.RowsAffected() == 0 {
		if _, err := u.ur.InitUpdateUserinf(newUserInf, ctx, user.UserID); err == nil {
			ctx.JSON(http.StatusCreated, newUpdateResponse(
				fmt.Sprintf("user update created w/ ID %d", user.UserID), true, "",
			))
			return
		}

		ctx.JSON(http.StatusBadRequest, newUpdateResponse(
			"", false, fmt.Sprintf("there is no user w/ ID %d", user.UserID),
		))
		return
	}

	ctx.JSON(http.StatusOK, newUpdateResponse(
		fmt.Sprintf("updated user w/ ID %d", user.UserID), true, "",
	))
}

func newHistoryResponse(res models.UserOrder, success bool, err string) models.HistoryResponse {
	return models.HistoryResponse{Result: res, Success: success, Error: err}
}

// HandleGetUserOrderHistory godoc
//
//	@Summary	get user order history info based from ID
//	@Tags		users
//	@Accept		json
//	@Produce	json
//	@Param		Authorization	header		string	true	"Bearer token"
//	@Success	200				{object}	models.UpdateResponse
//	@Router		/users/orders [get]
func (u *UserHandler) HandleUserOrderHistory(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	user, _ := claims.(pkg.Claims)

	history, err := u.ur.GetUserOrderHistory(ctx.Request.Context(), user.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newHistoryResponse(
			models.UserOrder{}, false, "server unable to get user order history",
		))
		return
	}

	ctx.JSON(http.StatusOK, newHistoryResponse(
		history, true, "",
	))
}
