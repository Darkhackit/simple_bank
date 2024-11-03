package api

import (
	db "github.com/Darkhackit/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
	"net/http"
)

type transferRequest struct {
	FromAccountID int64 `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64 `json:"to_account_id" binding:"required,gt=1"`
	Amount        int64 `json:"amount" binding:"required"`
	//Currency      string `json:"currency" binding:"required,oneof=USD EUR GHC"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var request transferRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: request.FromAccountID,
		ToAccountID:   request.ToAccountID,
		Amount:        request.Amount,
	}

	account, err := server.s.TransferTx(arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)

}
