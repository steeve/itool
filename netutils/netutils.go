package netutils

import (
	"context"
	"net"
	"net/url"
	"sync"
	"time"
)

type URLDialerFunc func(context.Context, *url.URL) (net.Conn, error)

var (
	dialersMu sync.RWMutex
	dialers   = map[string]URLDialerFunc{
		"unix": func(ctx context.Context, u *url.URL) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, "unix", u.Path)
		},
	}
)

func URLDialContext(ctx context.Context, host string) (net.Conn, error) {
	u, err := url.Parse(host)
	if err != nil {
		return nil, err
	}
	df := URLDialerFuncForScheme(u.Scheme)
	return df(ctx, u)
}

func RegisterURLScheme(scheme string, df URLDialerFunc) {
	dialersMu.Lock()
	defer dialersMu.Unlock()
	dialers[scheme] = df
}

func URLDialerFuncForScheme(scheme string) URLDialerFunc {
	dialersMu.RLock()
	defer dialersMu.RUnlock()
	if df, ok := dialers[scheme]; ok {
		return df
	}
	return func(ctx context.Context, u *url.URL) (net.Conn, error) {
		return (&net.Dialer{}).DialContext(ctx, u.Scheme, u.Host)
	}
}

type connCtx struct {
	net.Conn
	ctx context.Context
}

func ConnCtx(ctx context.Context, conn net.Conn) net.Conn {
	return connCtx{
		Conn: conn,
		ctx:  ctx,
	}
}

func (c connCtx) Read(p []byte) (int, error) {
	if t, ok := c.ctx.Deadline(); ok {
		if err := c.Conn.SetReadDeadline(t); err != nil {
			return 0, err
		}
		defer c.Conn.SetReadDeadline(time.Time{})
	}
	return c.Conn.Read(p)
}

func (c connCtx) Write(p []byte) (int, error) {
	if t, ok := c.ctx.Deadline(); ok {
		if err := c.Conn.SetWriteDeadline(t); err != nil {
			return 0, err
		}
		defer c.Conn.SetWriteDeadline(time.Time{})
	}
	return c.Conn.Write(p)
}
