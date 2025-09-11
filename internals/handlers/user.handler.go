package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/metgag/koda-weekly10/internals/models"
	"github.com/metgag/koda-weekly10/internals/repositories"
	"github.com/metgag/koda-weekly10/internals/utils"
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
		utils.PrintError(fmt.Sprintf("UNABLE TO REACH PROFILE WITH USER ID %d", user.UserID), 8, err)
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
		utils.PrintError("UNABLE TO BIND PROFILE UPDATE BODY", 12, err)
		ctx.JSON(http.StatusInternalServerError, newUpdateResponse(
			"", false, "server unable to bind input",
		))
		return
	}

	ctag, err := u.ur.UpdateUserinf(newUserInf, ctx, user.UserID)
	if err != nil {
		utils.PrintError("UNABLE TO MAKE UPDATE REQUEST", 12, err)
		ctx.JSON(http.StatusInternalServerError, newUpdateResponse(
			"", false, "server unable to update request",
		))
		return
	}
	if ctag.RowsAffected() == 0 {
		// if _, err := u.ur.InitUpdateUserinf(newUserInf, ctx, user.UserID); err == nil {
		// 	ctx.JSON(http.StatusCreated, newUpdateResponse(
		// 		fmt.Sprintf("user update created w/ ID %d", user.UserID), true, "",
		// 	))
		// 	return
		// }

		utils.PrintError(fmt.Sprintf("UNABLE TO UPDATE PROFILE WITH USER ID %d", user.UserID), 8, err)
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
		utils.PrintError("UNABLE TO GET USER ORDER HISTORIES", 12, err)
		ctx.JSON(http.StatusInternalServerError, newHistoryResponse(
			models.UserOrder{}, false, "server unable to get user order history",
		))
		return
	}
	if len(history.OrderHistory) == 0 {
		utils.PrintError("USER HAVE NO ORDER HISTORIES", 12, nil)
		ctx.JSON(http.StatusNoContent, newHistoryResponse(
			models.UserOrder{}, false, "user have no order yet",
		))
		return
	}

	ctx.JSON(http.StatusOK, newHistoryResponse(
		history, true, "",
	))
}

func newEditPasswordResponse(res, err string, success bool) models.EditPasswordResponse {
	return models.EditPasswordResponse{Result: res, Error: err, Success: success}
}

func (u *UserHandler) HandlePasswordEdit(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	user, _ := claims.(pkg.Claims)

	var newPwdBody models.PasswordBody
	if err := ctx.ShouldBindJSON(&newPwdBody); err != nil {
		utils.PrintError("UNABLE BIND PASSWORD TO BODY", 12, err)
		ctx.JSON(http.StatusBadRequest, newEditPasswordResponse(
			"", "server unable to bind request", false,
		))
		return
	}

	p := pkg.NewHashParams()
	p.UseRecommended()
	newEncodedHash, err := p.GenerateFromPassword(newPwdBody.Password)
	if err != nil {
		utils.PrintError("UNABLE TO HASH NEW PASSWORD", 8, err)
		ctx.JSON(http.StatusInternalServerError, newEditPasswordResponse(
			"", "server error while encode new password", false,
		))
		return
	}

	ctag, err := u.ur.UpdateUserPassword(ctx.Request.Context(), newEncodedHash, user.UserID)
	if err != nil {
		utils.PrintError("SERVER UNABLE TO MAKE PASSWORD UPDATE", 8, err)
		ctx.JSON(http.StatusInternalServerError, newEditPasswordResponse(
			"", "server error while updating new password", false,
		))
		return
	}

	if ctag.RowsAffected() == 0 {
		utils.PrintError(fmt.Sprintf("UNABLE TO UPDATE PASSWORD, NO USER ID %d", user.UserID), 8, err)
		ctx.JSON(http.StatusBadRequest, newEditPasswordResponse(
			"", fmt.Sprintf("there is no user w/ ID %d", user.UserID), false,
		))
		return
	}

	ctx.JSON(http.StatusOK, newEditPasswordResponse(
		fmt.Sprintf("succesfully update user's password w/ ID %d", user.UserID), "", true,
	))
}
