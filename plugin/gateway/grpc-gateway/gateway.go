package grpc_gateway

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/wishicorp/sdk/plugin/gateway"
	proto "github.com/wishicorp/sdk/plugin/gateway/grpc-gateway/proto"
	"github.com/wishicorp/sdk/plugin/pluginregister"
	"github.com/wishicorp/sdk/pool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"math"
	"net"
	"strings"
	"time"
)

var _ gateway.Gateway = (*GRPCGateway)(nil)

type GRPCGateway struct {
	pm          *pluginregister.PluginManager
	logger      hclog.Logger
	workerPool  *pool.WorkerPool
	ctx         context.Context
	cancel      context.CancelFunc
	running     chan bool
	workerSize  int
	security    gateway.Security
	grpcServer  *grpc.Server
	tcpListen   *net.TCPListener
	impl        proto.RpcGatewayServer
	authMethod  *gateway.Method
	authEnabled bool
}

func (m *GRPCGateway) SetAuthEnabled() {
	m.authEnabled = true
}

func (m *GRPCGateway) SetAuthDisabled() {
	m.authEnabled = false
}

func (m *GRPCGateway) SetAuthMethod(method string) error {
	defer func() {
		proto.RegisterRpcGatewayServer(m.grpcServer, m.NewGRPCGatewayImpl())
		reflection.Register(m.grpcServer)
	}()

	if method == "" {
		return nil
	}
	methods := strings.Split(method, ".")[:]
	if len(methods) != 3 {
		return errors.New("auth method error")
	}
	m.authMethod = &gateway.Method{
		Backend:   methods[0],
		Namespace: methods[1],
		Operation: methods[2],
	}

	return nil
}

func (m *GRPCGateway) SetSecurity(security gateway.Security) {
	m.security = security
}

func NewGateway(m *pluginregister.PluginManager, workerSize int, logger hclog.Logger) *GRPCGateway {
	ctx, cancel := context.WithCancel(context.Background())
	gw := &GRPCGateway{
		pm:          m,
		logger:      logger.Named("rpc-gateway"),
		ctx:         ctx,
		cancel:      cancel,
		running:     make(chan bool, 1),
		workerSize:  workerSize,
		authEnabled: true,
	}
	var opts []grpc.ServerOption
	opts = append(opts, grpc.MaxRecvMsgSize(math.MaxInt32))
	opts = append(opts, grpc.MaxSendMsgSize(math.MaxInt32))
	gw.grpcServer = grpc.NewServer(opts...)

	return gw
}
func (m *GRPCGateway) Listen(addr string, port uint) error {

	tcp, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", addr, port))
	if nil != err {
		return err
	}
	tcpListen, err := net.ListenTCP("tcp", tcp)
	if nil != err {
		return err
	}
	m.tcpListen = tcpListen

	return nil
}

func (m *GRPCGateway) Serve() (err error) {
	go func() {
		err = m.grpcServer.Serve(m.tcpListen)
	}()
	return err
}

//关闭网关
func (m *GRPCGateway) Shutdown() {
	defer func() {
		if m.logger.IsTrace() {
			m.logger.Trace("exited")
		}
	}()

	m.workerPool.Shutdown()
	m.grpcServer.Stop()
	select {
	case <-m.workerPool.Running():
	case <-time.After(time.Second * 1):
	}
	close(m.running)
}

//网关是否在运行(阻塞等待)
func (m *GRPCGateway) Running() <-chan bool {
	return m.running
}
