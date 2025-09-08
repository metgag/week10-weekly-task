package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/metgag/koda-weekly10/internals/models"
	"github.com/metgag/koda-weekly10/internals/repositories"
)

type OrderHandler struct {
	or *repositories.OrderRepository
}

func NewOrderHandler(or *repositories.OrderRepository) *OrderHandler {
	return &OrderHandler{or: or}
}

func newOrderResponse(res string, success bool, err string) models.OrderResponse {
	return models.OrderResponse{Result: res, Success: success, Error: err}
}

// HandleCreateOrder godoc
//
//	@Summary	create user order handler func
//	@Tags		orders
//	@Accept		json
//	@Produce	json
//	@Param		Authorization	header		string				true	"Bearer token"
//	@Param		request			body		models.CinemaOrder	true	"order body json content"
//	@Success	200				{object}	models.OrderResponse
//	@Router		/orders [post]
func (o *OrderHandler) HandleCreateOrder(ctx *gin.Context) {
	jwtSecret := os.Getenv("JWT_SECRET")
	token := ctx.GetHeader("Authorization")

	parsedToken, _ := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return []byte(jwtSecret), nil
	})
	claims, _ := parsedToken.Claims.(jwt.MapClaims)
	uid := claims["user_id"].(float64)

	var body models.CinemaOrder

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusInternalServerError, newOrderResponse(
			"", false, "server error while binding order",
		))
		return
	}

	res, err := o.or.CreateOrder(ctx, body, uid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newOrderResponse(
			"", false, "server unable to create order",
		))
		return
	}

	ctx.JSON(http.StatusOK, newOrderResponse(
		res, true, "",
	))
}

func (o *OrderHandler) HandleGetOrderHistory(ctx *gin.Context) {
	histories, err := o.or.GetOrderHistories(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newOrderResponse(
			"", false, "server unable to get order histories",
		))
		return
	}

	ctx.JSON(http.StatusOK, struct {
		Result  []models.OrderHistory `json:"result"`
		Success bool                  `json:"success"`
		Error   string                `json:"error"`
	}{
		histories,
		true,
		"",
	})
}
