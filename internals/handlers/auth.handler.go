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

func newRegisterResponse(err string, success bool, result string) models.RegisterResponse {
	return models.RegisterResponse{Error: err, Success: success, Result: result}
}

// HandleRegister godoc
// @Summary user register handler func
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.Register true "register body json content"
// @Success 200 {object} models.RegisterResponse "successful register response"
// @Router /auth/register [post]
func (a *AuthHandler) AddUser(ctx *gin.Context) {
	var body = models.Register{}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusInternalServerError, newRegisterResponse(
			"server unable to bind request", false, "",
		))
		return
	}

	p := pkg.NewHashParams()
	p.UseRecommended()
	encodedHash, err := p.GenerateFromPassword(body.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newRegisterResponse(
			"server error while encoding password", false, "",
		))
		return
	}

	id, err := a.ar.AddNewUser(ctx.Request.Context(), body.Email, encodedHash)
	if err != nil {
		ctx.JSON(http.StatusConflict, newRegisterResponse(
			"duplicate email addresses", false, "",
		))
		return
	}

	ctx.JSON(http.StatusOK, newRegisterResponse(
		"", true, fmt.Sprintf("register succesfully w/ ID: %d", id),
	))
}

func newLoginResponse(res, bearer string, success bool) models.LoginResponse {
	return models.LoginResponse{Result: res, Success: success, Bearer: bearer}
}

// HandleLogin godoc
// @Summary user login handler func
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.Login true "login body json content"
// @Success 200 {object} models.LoginResponse "successful login response"
// @Router /auth/login [post]
func (a *AuthHandler) Login(ctx *gin.Context) {
	var body = models.Login{}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusInternalServerError, newRegisterResponse(
			"server unable to bind request", false, "",
		))
		return
	}

	user, err := a.ar.GetUser(ctx.Request.Context(), body.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newRegisterResponse(
			"server unable to get user", false, "",
		))
		return
	}

	hc := pkg.NewHashParams()
	isMatch, err := hc.ComparePasswordAndHash(body.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newRegisterResponse(
			"server unable to compare password", false, "",
		))
		return
	}
	if !isMatch {
		ctx.JSON(http.StatusBadRequest, newRegisterResponse(
			"invalid email or password", false, "",
		))
		return
	}

	claims := pkg.NewJWTClaims(user.ID, user.Role)
	token, err := claims.GenAccessToken()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newRegisterResponse(
			"server unable to generate access token", false, "",
		))
		return
	}

	ctx.JSON(http.StatusOK, newLoginResponse(
		fmt.Sprintf("logged in as UID %d", claims.UserID), token, true,
	))
}
