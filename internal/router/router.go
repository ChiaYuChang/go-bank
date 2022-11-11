package router

import (
	"github.com/gin-gonic/gin"
	swaggerFile "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"
	v1 "gitlab.com/gjerry134679/bank/internal/router/api/v1"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.GET("swagger/*any", swagger.WrapHandler(swaggerFile.Handler))

	account := v1.NewAccount()

	apiv1 := r.Group("api/v1")
	{
		// New Account
		apiv1.POST("/account", account.Create)
		// Delete Account
		apiv1.DELETE("/account/:id", account.Delete)
		// Update Account Balance
		apiv1.PUT("/account/:id", account.Update)
		// Withdraw Money from Balance
		apiv1.PATCH("/account/withdraw/:id/amount", account.Withdraw)
		// Deposit Mondy from Balance
		apiv1.PATCH("/account/deposit/:id/amount", account.Deposit)
	}

	return r
}
