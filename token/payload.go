package token

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)
import "github.com/google/uuid"

var ErrExpiredToken = errors.New("token expired all")

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (payload *Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{Time: payload.ExpiredAt}, nil
}

func (payload *Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{Time: payload.IssuedAt}, nil
}

func (payload *Payload) GetNotBefore() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{Time: payload.IssuedAt}, nil
}

func (payload *Payload) GetIssuer() (string, error) {
	return payload.Username, nil
}

func (payload *Payload) GetSubject() (string, error) {
	return payload.Username, nil
}

func (payload *Payload) GetAudience() (jwt.ClaimStrings, error) {
	//TODO implement me
	panic("implement me")
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
