package grpc_gateway

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/wishicorp/sdk/plugin/gateway"
	"github.com/wishicorp/sdk/plugin/gateway/grpc-gateway/proto"
	"github.com/wishicorp/sdk/plugin/logical"
	"github.com/wishicorp/sdk/plugin/pluginregister"
	"github.com/wishicorp/sdk/pool"
	"google.golang.org/grpc"
	"math"
	"net"
	"time"
)

var _ gateway.Gateway = (*GRPCGateway)(nil)

type GRPCGateway struct {
	pm            *pluginregister.PluginManager
	logger        hclog.Logger
	workerPool    *pool.WorkerPool
	ctx           context.Context
	cancel        context.CancelFunc
	running       chan bool
	workerSize    int
	security      gateway.Security
	authenticator logical.PluginAuthenticator
	grpcServer    *grpc.Server
	tcpListen     *net.TCPListener
	impl          proto.RpcGatewayServer
}

func (m *GRPCGateway) SetSecurity(security gateway.Security) {
	panic("implement me")
}

func NewGateway(m *pluginregister.PluginManager, workerSize int, logger hclog.Logger) *GRPCGateway {
	ctx, cancel := context.WithCancel(context.Background())
	gw := &GRPCGateway{
		pm:         m,
		logger:     logger.Named("rpc-gateway"),
		ctx:        ctx,
		cancel:     cancel,
		running:    make(chan bool, 1),
		workerSize: workerSize,
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

func (m *GRPCGateway) SetPluginAuthorized(authenticator logical.PluginAuthenticator) {
	m.authenticator = authenticator
}

func (m *GRPCGateway) Serve() error {
	m.startWorkerPool(m.workerSize)
	proto.RegisterRpcGatewayServer(m.grpcServer, &proto.UnimplementedRpcGatewayServer{})
	m.grpcServer.Serve(m.tcpListen)
	return nil
}

//关闭网关
func (m *GRPCGateway) Shutdown() {
	defer func() {
		if m.logger.IsTrace() {
			m.logger.Trace("exited")
		}
	}()

	m.workerPool.Shutdown()

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
