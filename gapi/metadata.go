package gapi

import (
	"context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	xForwardedFHeader          = "x-forwarded-for"
	userAgentHeader            = "user-agent"
)

type Metadata struct {
	UserAgent string
	ClientIp  string
}

func (server *Server) extractMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if userAgent := md.Get(grpcGatewayUserAgentHeader); len(userAgent) > 0 {
			mtdt.UserAgent = userAgent[0]
		}
		if userAgent := md.Get(userAgentHeader); len(userAgent) > 0 {
			mtdt.UserAgent = userAgent[0]
		}
		if xFor := md.Get(xForwardedFHeader); len(xFor) > 0 {
			mtdt.ClientIp = xFor[0]
		}
	}
	if p, ok := peer.FromContext(ctx); ok {
		mtdt.ClientIp = p.Addr.String()
	}
	return mtdt
}
