package api

import (
	db "github.com/Darkhackit/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	_ "github.com/lib/pq"
	"net/http"
	"strings"
)

type createActionRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Balance  int64  `json:"balance"`
	Currency string `json:"currency" binding:"required,oneof=USD GHC"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var request createActionRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	owner := pgtype.Text{
		String: request.Owner,
		Valid:  true,
	}
	currency := pgtype.Text{
		String: request.Currency,
		Valid:  true,
	}
	balance := pgtype.Int8{
		Int64: request.Balance,
		Valid: true,
	}

	arg := db.CreateAccountParams{
		Owner:    owner,
		Balance:  balance,
		Currency: currency,
	}

	acct, err := server.q.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, acct)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	account, err := server.q.GetAccount(ctx, req.ID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

type listAccountsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var req listAccountsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListAccountParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := server.q.ListAccount(ctx, arg)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
