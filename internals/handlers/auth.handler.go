package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/metgag/koda-weekly10/internals/models"
	"github.com/metgag/koda-weekly10/internals/repositories"
	"github.com/metgag/koda-weekly10/internals/utils"
	"github.com/metgag/koda-weekly10/pkg"
)

type AuthHandler struct {
	ar *repositories.AuthRepository
}

func NewAuthHandler(ar *repositories.AuthRepository) *AuthHandler {
	return &AuthHandler{ar: ar}
}

func newRegisterResponse(err string, success bool, result string) models.RegisterResponse {
	return models.RegisterResponse{Error: err, Success: success, Result: result}
}

// HandleRegister godoc
//
//	@Summary		Register new user
//	@Description	Register a new user with email and password
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.Register			true	"User registration request body"
//	@Success		200		{object}	models.RegisterResponse	"Successfully registered"
//	@Failure		409		{object}	models.RegisterResponse	"Email already registered"
//	@Failure		500		{object}	models.RegisterResponse	"Internal server error"
//	@Router			/auth/register [post]
func (a *AuthHandler) HandleRegister(ctx *gin.Context) {
	var body = models.Register{}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		utils.LogCtxError(
			ctx,
			"UNABLE BINDING REGISTER BODY",
			"Internal server error",
			err,
			http.StatusInternalServerError,
		)
		return
	}

	p := pkg.NewHashParams()
	p.UseRecommended()
	encodedHash, err := p.GenerateFromPassword(body.Password)
	if err != nil {
		utils.LogCtxError(
			ctx,
			"UNABLE TO HASH PASSWORD",
			"Internal server error",
			err,
			http.StatusInternalServerError,
		)
		return
	}

	regisEmail, err := a.ar.AddNewUser(ctx.Request.Context(), body.Email, encodedHash)
	if err != nil {
		utils.LogCtxError(
			ctx,
			"EMAIL ALREADY REGISTERED",
			"Please log in or use a different email",
			err,
			http.StatusConflict,
		)
		return
	}

	ctx.JSON(http.StatusCreated, models.NewFullfilledResponse(
		http.StatusCreated,
		fmt.Sprintf("Succefully registered %s", regisEmail),
	))
}

func newLoginResponse(res, token string, success bool) models.LoginResponse {
	return models.LoginResponse{Result: res, Success: success, Token: token}
}

// HandleLogin godoc
//
//	@Summary		User login
//	@Description	Authenticates user by verifying email and password. Returns a JWT access token upon success.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.Login			true	"User login request body"
//	@Success		200		{object}	models.LoginResponse	"Login successful with JWT token"
//	@Failure		400		{object}	models.LoginResponse	"Invalid email or password"
//	@Failure		500		{object}	models.LoginResponse	"Internal server error"
//	@Router			/auth/login [post]
func (a *AuthHandler) HandleLogin(ctx *gin.Context) {
	var body = models.Login{}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		utils.LogCtxError(
			ctx,
			"UNABLE BINDING LOGIN BODY",
			"Internal server error",
			err,
			http.StatusInternalServerError,
		)
		return
	}

	user, err := a.ar.GetUser(ctx.Request.Context(), body.Email)
	if err != nil {
		utils.LogCtxError(
			ctx,
			"NO MATCHING USER",
			"Invalid email or password",
			err,
			http.StatusNotFound,
		)
		return
	}

	hc := pkg.NewHashParams()
	isMatch, err := hc.ComparePasswordAndHash(body.Password, user.Password)
	if err != nil {
		utils.LogCtxError(
			ctx,
			"UNABLE TO COMPARE LOGIN PASSWORD",
			"Internal server error",
			err,
			http.StatusInternalServerError,
		)
		return
	}
	if !isMatch {
		utils.LogCtxError(
			ctx,
			"INVALID LOGIN PASSWORD",
			"Invalid email or password",
			errors.New("invalid login password"),
			http.StatusBadRequest,
		)
		return
	}

	claims := pkg.NewJWTClaims(user.ID, user.Email, user.Password, user.Role)
	token, err := claims.GenAccessToken()
	if err != nil {
		utils.LogCtxError(
			ctx,
			"UNABLE TO GENERATE ACCESS TOKEN",
			"Internal server error",
			err,
			http.StatusInternalServerError,
		)
		return
	}

	ctx.JSON(http.StatusOK, newLoginResponse(
		fmt.Sprintf("Logged in as %s", user.Email), token, true,
	))
}

func newLogoutResponse(res, err string, success bool) models.LogoutResponse {
	return models.LogoutResponse{Result: res, Success: success, Error: err}
}

// HandleLogout godoc
//
//	@Summary		Logout user
//	@Description	Logout user with blacklist to redis
//	@Tags			auth
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.LogoutResponse
//
// // @Failure      500  {object}  models.LogoutResponse
//
//	@Router			/auth/logout [delete]
func (a *AuthHandler) HandleLogout(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	claims, _ := ctx.Get("claims")
	user, _ := claims.(pkg.Claims)

	// token = strings.TrimPrefix(token, "Bearer ")
	token = strings.Split(token, " ")[1]

	err := a.ar.SetLogoutCache(ctx.Request.Context(), token, user.IssuedAt.Time)
	if err != nil {
		log.Println("SetLogoutCache error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, newLogoutResponse(
			"", "SERVER ERROR WHILE LOG OUT", false,
		))
		return
	}

	ctx.JSON(http.StatusOK, newLogoutResponse(
		"logout succesfully", "", true,
	))
}
