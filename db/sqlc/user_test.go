package db

import (
	"context"
	"github.com/Darkhackit/simplebank/util"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func CreateRandomUser(t *testing.T) Users {
	name := pgtype.Text{
		String: util.RandomOwner(),
		Valid:  true,
	}
	active := pgtype.Bool{
		Bool:  true,
		Valid: true,
	}

	password, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	rs := util.RandomOwner()

	arg := CreateUserParams{
		Name:     name,
		Active:   active,
		Email:    rs + "@gmail.com",
		Password: string(password),
		Username: rs,
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, name, user.Name)
	require.Equal(t, active, user.Active)
	require.Equal(t, string(password), user.Password)
	require.Equal(t, rs, user.Username)
	return user
}

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user := CreateRandomUser(t)
	ctx := context.Background()

	getUser, err := testQueries.GetUser(ctx, user.Username)
	require.NoError(t, err)
	require.Equal(t, user.Name, getUser.Name)
	require.Equal(t, user.Active, getUser.Active)
	require.Equal(t, user.Email, getUser.Email)
	require.Equal(t, user.Password, getUser.Password)
	require.Equal(t, user.Username, getUser.Username)
}
