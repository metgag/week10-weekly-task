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
		utils.PrintError("UNABLE TO BIND REGISTER BODY", 8, err)
		ctx.JSON(http.StatusInternalServerError, newRegisterResponse(
			"server unable to bind request", false, "",
		))
		return
	}

	p := pkg.NewHashParams()
	p.UseRecommended()
	encodedHash, err := p.GenerateFromPassword(body.Password)
	if err != nil {
		utils.PrintError("UNABLE TO ENCODE PASSWORD", 8, err)
		ctx.JSON(http.StatusInternalServerError, newRegisterResponse(
			"server error while encode password", false, "",
		))
		return
	}

	id, err := a.ar.AddNewUser(ctx.Request.Context(), body.Email, encodedHash)
	if err != nil {
		utils.PrintError("EMAIL ALREADY REGISTERED", 12, err)
		ctx.JSON(http.StatusConflict, newRegisterResponse(
			"duplicate email addresses", false, "",
		))
		return
	}

	ctx.JSON(http.StatusOK, newRegisterResponse(
		"", true, fmt.Sprintf("register succesfully w/ ID: %d", id),
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
		utils.PrintError("UNABLE TO BIND LOGIN BODY", 8, err)
		ctx.JSON(http.StatusInternalServerError, newRegisterResponse(
			"server unable to bind request", false, "",
		))
		return
	}

	user, err := a.ar.GetUser(ctx.Request.Context(), body.Email)
	if err != nil {
		utils.PrintError("NO MATCHING USER", 12, err)
		ctx.JSON(http.StatusInternalServerError, newRegisterResponse(
			"server unable to get user", false, "",
		))
		return
	}

	hc := pkg.NewHashParams()
	isMatch, err := hc.ComparePasswordAndHash(body.Password, user.Password)
	if err != nil {
		utils.PrintError("UNABLE TO COMPARE LOGIN PASSWORD", 8, err)
		ctx.JSON(http.StatusInternalServerError, newRegisterResponse(
			"server unable to compare password", false, "",
		))
		return
	}
	if !isMatch {
		utils.PrintError("INVALID CREDENTIALS", 12, nil)
		ctx.JSON(http.StatusBadRequest, newRegisterResponse(
			"invalid email or password", false, "",
		))
		return
	}

	claims := pkg.NewJWTClaims(user.ID, user.Password, user.Role)
	token, err := claims.GenAccessToken()
	if err != nil {
		utils.PrintError("FAIL GENERATE ACCESS TOKEN", 12, err)
		ctx.JSON(http.StatusInternalServerError, newRegisterResponse(
			"server unable to generate access token", false, "",
		))
		return
	}

	ctx.JSON(http.StatusOK, newLoginResponse(
		fmt.Sprintf("logged in as UID %d", claims.UserID), token, true,
	))
}
