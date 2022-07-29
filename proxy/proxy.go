package proxy

import (
	"context"
	"encoding/json"
	"github.com/gdong42/grpc-mate/metadata"
	"github.com/gdong42/grpc-mate/proxy"
)

type Proxy struct {
	conn *proxy.Proxy
}

// New -
func New(conn *proxy.Proxy) *Proxy {
	return &Proxy{conn}
}

// Call a service
func (p *Proxy) Call(ctx context.Context, service, method string, data interface{}) ([]byte, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return p.conn.Invoke(ctx, service, method, b, &metadata.Metadata{})
}
