package connman

import (
	"errors"
	"sync"

	"github.com/EwanValentine/grpc-proxy-framework/config"
	"github.com/EwanValentine/grpc-proxy-framework/proxy"

	pproxy "github.com/gdong42/grpc-mate/proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Connections map[string]string

type Connection struct {
	Proxy      *proxy.Proxy
	connection *grpc.ClientConn
}

type ActiveConnections map[string]Connection

type ConnectionManager struct {
	connections       Connections
	activeConnections ActiveConnections
	config            *config.Config
	mu                sync.RWMutex
}

// New ConnectionManager
func New(conf *config.Config) *ConnectionManager {
	return &ConnectionManager{
		connections:       make(Connections),
		activeConnections: make(ActiveConnections),
		config:            conf,
		mu:                sync.RWMutex{},
	}
}

func (cm *ConnectionManager) Remove(name string) {
	cm.mu.Lock()
	delete(cm.connections, name)
	cm.mu.Unlock()
}

// Start all connections
func (cm *ConnectionManager) Start() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	for name, addr := range cm.config.Connections {
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return err
		}

		proxyConn := proxy.New(pproxy.NewProxy(conn))

		cm.activeConnections[name] = Connection{
			Proxy:      proxyConn,
			connection: conn,
		}
	}

	return nil
}

var (
	ErrNotFound = errors.New("connection not found")
)

// GetByName get connection by name
func (cm *ConnectionManager) GetByName(name string) (*proxy.Proxy, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	conn, ok := cm.activeConnections[name]
	if !ok {
		return nil, ErrNotFound
	}
	return conn.Proxy, nil
}

// Close all connections
func (cm *ConnectionManager) Close() {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	for _, conn := range cm.activeConnections {
		conn.connection.Close()
	}
}
