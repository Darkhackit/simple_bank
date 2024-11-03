package token

import (
	"fmt"
	"github.com/Darkhackit/simplebank/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPasetoMaker(t *testing.T) {
	secret := util.RandomString(32)
	maker, err := NewPasetoToken(secret)
	require.NoError(t, err)

	fmt.Println("secret", secret)

	username := util.RandomOwner()

	duration := time.Minute

	issuedAt := time.Now()

	expiredAt := time.Now().Add(duration)

	token, err := maker.CreateToken(username, duration)

	fmt.Println("token", token)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)

	fmt.Println("payload", payload)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.Equal(t, payload.Username, username)
	require.NotZero(t, payload.ID)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)

}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoToken(util.RandomString(32))
	require.NoError(t, err)

	token, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.Nil(t, payload)

}
