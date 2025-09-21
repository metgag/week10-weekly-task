package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/metgag/koda-weekly10/internals/models"
	"github.com/metgag/koda-weekly10/internals/repositories"
	"github.com/metgag/koda-weekly10/internals/utils"
	"github.com/metgag/koda-weekly10/pkg"
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
//	@Summary		create user order handler func
//	@Description	create a new user's order
//	@Tags			orders
//	@Accept			json
//	@Produce		json
//
//	@Param			order	body		models.CinemaOrder		true	"Order body"	example({"seats": [3, 4, 5]})
//
//	@Success		200		{object}	models.OrderResponse	"Order created successfully"
//	@Failure		400		{object}	models.OrderResponse	"Invalid request payload"
//	@Failure		401		{object}	models.OrderResponse	"Unauthorized: invalid or missing token"
//	@Failure		500		{object}	models.OrderResponse	"Internal server error"
//	@Security		BearerAuth
//	@Router			/orders [post]
func (o *OrderHandler) HandleCreateOrder(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	user, _ := claims.(pkg.Claims)

	var body models.CinemaOrder

	if err := ctx.ShouldBindJSON(&body); err != nil {
		utils.PrintError("UNABLE TO BIND ORDER BODY", 12, err)
		ctx.JSON(http.StatusInternalServerError, newOrderResponse(
			"", false, "server error while binding order",
		))
		return
	}

	res, err := o.or.CreateOrder(ctx, body, user.UserID)
	if err != nil {
		utils.PrintError("UNABLE CREATE ORDER", 12, err)
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
		utils.PrintError("UNABLE GET ORDER HISTORIES", 12, err)
		ctx.JSON(http.StatusInternalServerError, newOrderResponse(
			"", false, "server unable to get order histories",
		))
		return
	}

	ctx.JSON(http.StatusOK, models.OrderHistoriesResponse{
		Result:  histories,
		Success: true,
		Error:   "",
	})
}
