package v1

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/gjerry134679/bank/internal/router/api"
)

type Account struct{ api.Empty }

func NewAccount() Account {
	return Account{}
}

func (a Account) Withdraw(ctx *gin.Context) {

}

func (a Account) Deposit(ctx *gin.Context) {

}
