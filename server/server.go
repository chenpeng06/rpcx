package server

import (
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrServerClosed  = errors.New("http:	Server closed")
	ErrReqReachLimit = errors.New("request reach rate limit")
)

const (
	// ReaderBuffsize is used for bufio reader.
	ReaderBufSize = 1024
	WriteBufSize  = 1024
)

type contextKey struct {
	name string
}

func (c *contextKey) String() string {
	return "rpcx context value " + c.name
}

var (
	// RemoteConnContextKey is a context key. It can be used in
	// services with context package to access to the connection arrived on.
	RemoteConnContextKey = &contextKey{"remote-conn"}
	// StartRequestContextKey records the start time
	StartRequestContextKey = &contextKey{"start-parse-request"}
	// StartSendRequestContextKey records the start time
	StartSendRequestContextKey = &contextKey{"start-send-request"}
	// TagsContextKey is used to record extra info in handling services. Its value is a map[string]interface{}.
	TagsContextKey = &contextKey{"service-tag"}
	// HttpConnContextKey is used to store http connection.
	HttpConnContextKey = &contextKey{"http-conn"}
)

type Handler func(ctx *Context) error

type WorkerPool interface {
	Submit(task func())
	StopAndWaitFor(deadline time.Duration)
	Stop()
	StopAndWait()
}

type Server struct {
	ln                 net.Listener
	readTimeout        time.Duration
	writeTimeout       time.Duration
	gatewayHTTPServer  *http.Server
	jsonrpcHTTPServer  *http.Server
	DisableHTTPGateway bool // disable http invoke or not.
	DisableJSONRPC     bool // disable jsonrpc invoke or not.
	EnableProfile      bool // enable profile and statsview or not
	AsyncWrite         bool // set true if your server only serves few clients
	pool               WorkerPool

	serviceMapMu sync.RWMutex
	serviceMap   map[string]*service

	router map[string]Handler

	mu         sync.Mutex
	activeConn map[net.Conn]struct{}
	doneChan   chan struct{}
	seq        atomic.Uint64

	inShutdown int32
	onShutdown []func(s *Server)
	onRestart  []func(s *Server)

	// TLSConfig for creating tls tcp connection
	tlsConfig *tls.Config
	// BlockCrypt for kcp.BlockCrypt
	options map[string]interface{}
	// CORS options
	corsOptions *CORSOptions
}
