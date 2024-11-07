package gapi

import (
	"fmt"
	"github.com/Darkhackit/simplebank/token"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
	"strings"
)

const (
	authorizationHeader = "authorization"
	authorizationBearer = "bearer"
)

func (server *Server) authorizeUser(ctx context.Context) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	values := md.Get(authorizationHeader)
	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}
	authHeader := values[0]
	field := strings.Fields(authHeader)
	if len(field) != 2 {
		return nil, fmt.Errorf("Invalid authorization header")
	}

	authType := strings.ToLower(field[0])
	if authType != authorizationBearer {
		return nil, fmt.Errorf("unsupported authorization type")
	}

	accessToken := field[1]

	payload, err := server.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, err
	}
	return payload, nil

}
