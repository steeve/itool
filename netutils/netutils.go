package netutils

import (
	"context"
	"net"
	"net/url"
)

func DialURLContext(ctx context.Context, host string) (net.Conn, error) {
	u, err := url.Parse(host)
	if err != nil {
		return nil, err
	}
	dialer := net.Dialer{}
	if u.Scheme == "unix" {
		return dialer.DialContext(ctx, u.Scheme, u.Path)
	}
	return dialer.DialContext(ctx, u.Scheme, u.Host)
}
