package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"net/http"
	"time"
)

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken string    `json:"access_token"`
	ExpiresIn   time.Time `json:"expires_in"`
}

func (server *Server) renewAccessToken(c *gin.Context) {
	var request renewAccessTokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	token, err := server.tokenMaker.VerifyToken(request.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	uuid := pgtype.UUID{
		Bytes: token.ID,
		Valid: true,
	}
	session, err := server.q.GetSession(c, uuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	if session.IsBlocked {
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	if session.Username != token.Username {
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	if session.RefreshToken != request.RefreshToken {
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	accessToken, payloadToken, err := server.tokenMaker.CreateToken(session.Username, server.config.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	response := renewAccessTokenResponse{
		AccessToken: accessToken,
		ExpiresIn:   payloadToken.ExpiredAt,
	}

	c.JSON(http.StatusOK, response)

}
