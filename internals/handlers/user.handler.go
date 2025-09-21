package handlers

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

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
//	@Summary		get user profile info based from ID
//	@Description	get user's profile details
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.UserinfResponse	"User profile retrieved successfully"
//	@Failure		401	{object}	models.UserinfResponse	"Unauthorized or invalid token"
//	@Failure		500	{object}	models.UserinfResponse	"Internal server error"
//	@Security		BearerAuth
//	@Router			/users [get]
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
	userinf.Role = user.Role

	ctx.JSON(http.StatusOK, newUserinfResponse(
		userinf, true, "",
	))
}

func newUpdateResponse(res string, success bool, err string) models.UpdateResponse {
	return models.UpdateResponse{Result: res, Success: success, Error: err}
}

// HandleUpdateUserProfile godoc
//
//	@Summary		Update user profile info by ID
//	@Description	Update user's profile details via multipart form
//	@Tags			users
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			first_name		formData	string					false	"First name"
//	@Param			last_name		formData	string					false	"Last name"
//	@Param			phone_number	formData	string					false	"Phone number (e.g., 08667728761)"
//	@Param			point_count		formData	number					false	"User's point (e.g., 4.8)"
//	@Param			avatar			formData	file					false	"Avatar image file"
//	@Success		200				{object}	models.UpdateResponse	"User profile updated successfully"
//	@Failure		400				{object}	models.UpdateResponse	"Invalid user ID or no user found"
//	@Failure		500				{object}	models.UpdateResponse	"Server error while updating profile"
//	@Security		BearerAuth
//	@Router			/users [patch]
func (u *UserHandler) HandleUpdateUserInf(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	user, _ := claims.(pkg.Claims)

	var newUserInf models.NewInf
	if err := ctx.ShouldBind(&newUserInf); err != nil {
		utils.PrintError("UNABLE TO BIND PROFILE UPDATE BODY", 12, err)
		ctx.JSON(http.StatusInternalServerError, newUpdateResponse(
			"", false, "server unable to bind input",
		))
		return
	}

	// log.Println("AVAVAVAVA", newUserInf)

	avatar := newUserInf.Avatar
	var avatarName string
	if avatar != nil {
		ext := filepath.Ext(avatar.Filename)
		filename := fmt.Sprintf("avatar_%d_%d%s", user.UserID, time.Now().Unix(), ext)
		location := filepath.Join("public", "user", filename)
		if err := ctx.SaveUploadedFile(avatar, location); err != nil {
			utils.PrintError("INVALID AVATAR IMAGE", 12, err)
			ctx.JSON(http.StatusBadRequest, newUpdateMovieResponse(
				false, "", "unable to upload movie backdrop",
			))
			return
		}
		avatarName = filename
		log.Println("====================", filename)
	}

	ctag, err := u.ur.UpdateUserinf(newUserInf, ctx, user.UserID, avatarName)
	if err != nil {
		utils.PrintError("UNABLE TO MAKE UPDATE REQUEST", 12, err)
		ctx.JSON(http.StatusInternalServerError, newUpdateResponse(
			"", false, "server unable to update request",
		))
		return
	}
	if ctag.RowsAffected() == 0 {
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
//	@Summary		get user order history info based from ID
//
//	@Description	get user's order watch histories
//
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.HistoryResponse	"User order history retrieved successfully"
//	@Failure		204	{object}	models.HistoryResponse	"No order history found for user"
//	@Failure		401	{object}	models.HistoryResponse	"Unauthorized or invalid token"
//	@Failure		500	{object}	models.HistoryResponse	"Internal server error"
//	@Security		BearerAuth
//	@Router			/users/orders [get]
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

// HandlePasswordEdit godoc
//
//	@Summary		update user's password
//	@Description	allows user's to update their password
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.PasswordBody			true	"new password JSON body"
//	@Success		200		{object}	models.EditPasswordResponse	"password updated successfully"
//	@Failure		400		{object}	models.EditPasswordResponse	"bad request - validation failed or user not found"
//	@Failure		500		{object}	models.EditPasswordResponse	"server error while updating password"
//	@Security		BearerAuth
//	@Router			/users/password [patch]
func (u *UserHandler) HandlePasswordEdit(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	user, _ := claims.(pkg.Claims)

	// validasi apakah menggunakan token terbaru
	tokenIat := user.IssuedAt.Time.Unix()
	lastUpdate, err := u.ur.GetLastUpdated(ctx.Request.Context(), user.UserID)
	if err != nil {
		utils.PrintError("ERR GET LAST UPDATE WHILE EDIT PASSWORD", 12, err)
		ctx.JSON(http.StatusInternalServerError, newEditPasswordResponse(
			"", "internal server error", false,
		))
		return
	}
	if tokenIat < lastUpdate.Unix() {
		ctx.JSON(http.StatusUnauthorized, newEditPasswordResponse(
			"", "access token revoked", false,
		))
		return
	}

	var newPwdBody models.PasswordBody
	if err := ctx.ShouldBindJSON(&newPwdBody); err != nil {
		utils.PrintError("UNABLE BIND PASSWORD TO BODY", 12, err)
		ctx.JSON(http.StatusBadRequest, newEditPasswordResponse(
			"", "unable to bind request", false,
		))
		return
	}

	p := pkg.NewHashParams()
	p.UseRecommended()

	isMatch, err := p.ComparePasswordAndHash(newPwdBody.OldPassword, user.Password)
	if err != nil {
		utils.PrintError("UNABLE TO COMPARING PASSWORD", 12, err)
		ctx.JSON(http.StatusInternalServerError, newEditPasswordResponse(
			"", "server error while compare old password", false,
		))
		return
	}
	if !isMatch {
		utils.PrintError("UNABLE TO EDIT PWD, OLD PASSWORD MISMATCH", 12, nil)
		ctx.JSON(http.StatusBadRequest, newEditPasswordResponse(
			"", "old password mismatch", false,
		))
		return
	}

	now := time.Now()
	newEncodedHash, err := p.GenerateFromPassword(newPwdBody.NewPassword)
	if err != nil {
		utils.PrintError("UNABLE TO HASH NEW PASSWORD", 8, err)
		ctx.JSON(http.StatusInternalServerError, newEditPasswordResponse(
			"", "server error while encode new password", false,
		))
		return
	}

	ctag, err := u.ur.UpdateUserPassword(ctx.Request.Context(), newEncodedHash, user.UserID, now)
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
