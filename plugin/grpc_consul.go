package plugin

import (
	"context"
	"github.com/hashicorp/consul/api"
	"github.com/wishicorp/sdk/library/consul"
	"github.com/wishicorp/sdk/plugin/logical"
	"github.com/wishicorp/sdk/plugin/proto"
	"time"

	"google.golang.org/grpc"
)

var _ logical.Consul = (*GRPCConsulClient)(nil)

func newGRPCConsulClient(conn *grpc.ClientConn) *GRPCConsulClient {
	return &GRPCConsulClient{
		client: proto.NewConsulClient(conn),
	}
}

// GRPCConsulClient is an implementation of logical.Consul that communicates
// over RPC.
type GRPCConsulClient struct {
	client proto.ConsulClient
}

func (s *GRPCConsulClient) KVList(ctx context.Context, prefix string) (api.KVPairs, error) {
	resp, err := s.client.KVList(ctx, &proto.KVListArgs{Prefix: prefix})
	if nil != err {
		return nil, err
	}
	var kvps api.KVPairs
	for _, pair := range resp.KvPairs {
		kvp := api.KVPair{
			Key:         pair.Key,
			CreateIndex: pair.CreateIndex,
			ModifyIndex: pair.ModifyIndex,
			LockIndex:   pair.LockIndex,
			Flags:       pair.Flags,
			Value:       pair.Value,
			Session:     pair.Session,
		}
		kvps = append(kvps, &kvp)
	}
	return kvps, nil
}

func (s *GRPCConsulClient) KVCreate(ctx context.Context, p *api.KVPair) error {
	req := proto.KVCasArgs{Kvpair: &proto.KVPair{
		Key:     p.Key,
		Flags:   p.Flags,
		Value:   p.Value,
		Session: p.Session,
	}}
	_, err := s.client.KVCreate(ctx, &req)
	return err
}

func (s *GRPCConsulClient) GetConfig(ctx context.Context, key, version string, sandbox bool) ([]byte, error) {
	req := proto.GetConfigArgs{Key: key, Version: version, Sandbox: sandbox}
	resp, err := s.client.GetConfig(ctx, &req)
	if nil != err {
		return nil, err
	}
	return resp.Value, nil
}

func (s *GRPCConsulClient) GetService(ctx context.Context, id, tags string) (*api.AgentService, error) {
	req := proto.GetServiceArgs{
		Name: id,
		Tags: tags,
	}
	resp, err := s.client.GetService(ctx, &req, largeMsgGRPCCallOpts...)
	if nil != err {
		return nil, err
	}
	address := map[string]api.ServiceAddress{}
	for k, a := range resp.TaggedAddresses {
		addr := api.ServiceAddress{
			Address: a.Address,
			Port:    int(a.Port),
		}
		address[k] = addr
	}
	return &api.AgentService{
		Kind:            api.ServiceKind(resp.Kind),
		ID:              resp.Id,
		Service:         resp.Service,
		Tags:            resp.Tags,
		Meta:            resp.Meta,
		Port:            int(resp.Port),
		Address:         resp.Address,
		TaggedAddresses: address,
		Weights: api.AgentWeights{
			Passing: int(resp.Weights.Passing),
			Warning: int(resp.Weights.Warning),
		},
		EnableTagOverride: resp.EnableTagOverride,
	}, nil
}

func (s *GRPCConsulClient) NewSession(ctx context.Context, name string, ttl time.Duration, behavior consul.SessionBehavior) (string, error) {
	req := proto.NewSessionArgs{
		Name:     name,
		Ttl:      ttl.String(),
		Behavior: string(behavior),
	}
	resp, err := s.client.NewSession(ctx, &req, largeMsgGRPCCallOpts...)
	if nil != err {
		return "", err
	}
	return resp.Name, nil
}

func (s *GRPCConsulClient) SessionInfo(ctx context.Context, id string) (*api.SessionEntry, error) {
	req := proto.SessionInfoArgs{
		Id: id,
	}
	resp, err := s.client.SessionInfo(ctx, &req, largeMsgGRPCCallOpts...)
	if nil != err {
		return nil, err
	}
	duration, _ := time.ParseDuration(resp.Entry.LockDelay)
	var checks []api.ServiceCheck
	for _, check := range resp.Entry.ServiceChecks {
		c := api.ServiceCheck{
			ID:        check.ID,
			Namespace: check.Namespace,
		}
		checks = append(checks, c)
	}
	return &api.SessionEntry{
		CreateIndex:   resp.Entry.CreateIndex,
		ID:            resp.Entry.ID,
		Name:          resp.Entry.Name,
		Node:          resp.Entry.Node,
		LockDelay:     duration,
		Behavior:      resp.Entry.Behavior,
		TTL:           resp.Entry.TTL,
		Checks:        resp.Entry.Checks,
		NodeChecks:    resp.Entry.NodeChecks,
		ServiceChecks: checks,
	}, err
}

func (s *GRPCConsulClient) DestroySession(ctx context.Context, id string) error {
	req := proto.DestroySessionArgs{
		Id: id,
	}
	_, err := s.client.DestroySession(ctx, &req)
	return err
}

func (s *GRPCConsulClient) KVAcquire(ctx context.Context, key, session string) (success bool, err error) {
	req := proto.KVAcquireArgs{
		Key:     key,
		Session: session,
	}
	resp, err := s.client.KVAcquire(ctx, &req)
	if nil != err {
		return false, err
	}
	return resp.Success, nil
}

func (s *GRPCConsulClient) KVRelease(ctx context.Context, key string) error {
	req := proto.KVReleaseArgs{
		Key: key,
	}
	_, err := s.client.KVRelease(ctx, &req)
	return err
}

func (s *GRPCConsulClient) KVInfo(ctx context.Context, key string) (*api.KVPair, error) {
	req := proto.KVInfoArgs{
		Key: key,
	}
	resp, err := s.client.KVInfo(ctx, &req)
	if nil != err {
		return nil, err
	}
	return &api.KVPair{
		Key:         resp.Kvpair.Key,
		CreateIndex: resp.Kvpair.CreateIndex,
		ModifyIndex: resp.Kvpair.ModifyIndex,
		LockIndex:   resp.Kvpair.LockIndex,
		Flags:       resp.Kvpair.Flags,
		Value:       resp.Kvpair.Value,
		Session:     resp.Kvpair.Session,
	}, nil
}

func (s *GRPCConsulClient) KVCas(ctx context.Context, p *api.KVPair) (bool, error) {
	req := proto.KVCasArgs{
		Kvpair: &proto.KVPair{
			Key:         p.Key,
			CreateIndex: p.CreateIndex,
			ModifyIndex: p.ModifyIndex,
			LockIndex:   p.LockIndex,
			Flags:       p.Flags,
			Value:       p.Value,
			Session:     p.Session,
		},
	}
	resp, err := s.client.KVCas(ctx, &req)
	if nil != err {
		return false, err
	}
	return resp.Success, nil
}

// ConsulServer is a net/rpc compatible structure for serving
type GRPCConsulServer struct {
	impl logical.Consul
}

func (s *GRPCConsulServer) KVCreate(ctx context.Context, args *proto.KVCasArgs) (*proto.ConsulEmpty, error) {
	err := s.impl.KVCreate(ctx, &api.KVPair{
		Key:     args.Kvpair.Key,
		Flags:   args.Kvpair.Flags,
		Value:   args.Kvpair.Value,
		Session: args.Kvpair.Session,
	})
	return &proto.ConsulEmpty{}, err
}

func (s *GRPCConsulServer) KVList(ctx context.Context, args *proto.KVListArgs) (*proto.KVListReply, error) {
	resp, err := s.impl.KVList(ctx, args.Prefix)
	if nil != err {
		return nil, err
	}
	var kvps []*proto.KVPair
	for _, pair := range resp {
		kvp := proto.KVPair{
			Key:         pair.Key,
			CreateIndex: pair.CreateIndex,
			ModifyIndex: pair.ModifyIndex,
			LockIndex:   pair.LockIndex,
			Flags:       pair.Flags,
			Value:       pair.Value,
			Session:     pair.Session,
		}
		kvps = append(kvps, &kvp)
	}
	return &proto.KVListReply{KvPairs: kvps}, nil
}

func (s *GRPCConsulServer) GetConfig(ctx context.Context, args *proto.GetConfigArgs) (*proto.GetConfigReply, error) {
	resp, err := s.impl.GetConfig(ctx, args.Key, args.Version, args.Sandbox)
	if nil != err {
		return nil, err
	}
	return &proto.GetConfigReply{
		Value: resp,
	}, nil
}

func (s *GRPCConsulServer) GetService(ctx context.Context, args *proto.GetServiceArgs) (*proto.GetServiceReply, error) {
	resp, err := s.impl.GetService(ctx, args.Name, args.Tags)
	if nil != err {
		return nil, err
	}
	address := map[string]*proto.ServiceAddress{}
	for k, addr := range resp.TaggedAddresses {
		a := proto.ServiceAddress{
			Address: addr.Address,
			Port:    int32(addr.Port),
		}
		address[k] = &a
	}
	weights := &proto.AgentWeights{
		Passing: int32(resp.Weights.Passing),
		Warning: int32(resp.Weights.Warning),
	}
	return &proto.GetServiceReply{
		Kind:              string(resp.Kind),
		Id:                resp.ID,
		Service:           resp.Service,
		Tags:              resp.Tags,
		Meta:              resp.Meta,
		Port:              int32(resp.Port),
		Address:           resp.Address,
		TaggedAddresses:   address,
		Weights:           weights,
		EnableTagOverride: resp.EnableTagOverride,
	}, nil
}

func (s *GRPCConsulServer) NewSession(ctx context.Context, args *proto.NewSessionArgs) (*proto.NewSessionReply, error) {
	duration, _ := time.ParseDuration(args.Ttl)
	resp, err := s.impl.NewSession(ctx, args.Name, duration, consul.SessionBehavior(args.Behavior))
	if nil != err {
		return nil, err
	}
	return &proto.NewSessionReply{Name: resp}, nil

}

func (s *GRPCConsulServer) SessionInfo(ctx context.Context, args *proto.SessionInfoArgs) (*proto.SessionInfoReply, error) {
	resp, err := s.impl.SessionInfo(ctx, args.Id)
	if nil != err {
		return nil, err
	}
	var checks []*proto.ServiceCheck
	for _, check := range resp.ServiceChecks {
		c := &proto.ServiceCheck{
			ID:        check.ID,
			Namespace: check.Namespace,
		}
		checks = append(checks, c)
	}
	return &proto.SessionInfoReply{Entry: &proto.SessionEntry{
		CreateIndex:   resp.CreateIndex,
		ID:            resp.ID,
		Name:          resp.Name,
		Node:          resp.Node,
		LockDelay:     resp.LockDelay.String(),
		Behavior:      resp.Behavior,
		TTL:           resp.TTL,
		Checks:        resp.Checks,
		NodeChecks:    resp.NodeChecks,
		ServiceChecks: checks,
	}}, nil
}

func (s *GRPCConsulServer) DestroySession(ctx context.Context, args *proto.DestroySessionArgs) (*proto.ConsulEmpty, error) {
	err := s.impl.DestroySession(ctx, args.Id)
	return &proto.ConsulEmpty{}, err

}

func (s *GRPCConsulServer) KVAcquire(ctx context.Context, args *proto.KVAcquireArgs) (*proto.KVAcquireReply, error) {
	resp, err := s.impl.KVAcquire(ctx, args.Key, args.Session)
	if nil != err {
		return nil, err
	}
	return &proto.KVAcquireReply{Success: resp}, nil

}

func (s *GRPCConsulServer) KVRelease(ctx context.Context, args *proto.KVReleaseArgs) (*proto.ConsulEmpty, error) {
	err := s.impl.KVRelease(ctx, args.Key)
	return &proto.ConsulEmpty{}, err
}

func (s *GRPCConsulServer) KVInfo(ctx context.Context, args *proto.KVInfoArgs) (*proto.KVInfoReply, error) {
	resp, err := s.impl.KVInfo(ctx, args.Key)
	if nil != err {
		return nil, err
	}
	return &proto.KVInfoReply{Kvpair: &proto.KVPair{
		Key:         resp.Key,
		CreateIndex: resp.CreateIndex,
		ModifyIndex: resp.ModifyIndex,
		LockIndex:   resp.ModifyIndex,
		Flags:       resp.Flags,
		Value:       resp.Value,
		Session:     resp.Session,
	}}, nil
}

func (s *GRPCConsulServer) KVCas(ctx context.Context, args *proto.KVCasArgs) (*proto.KVCasReply, error) {
	kvp := api.KVPair{
		Key:         args.Kvpair.Key,
		CreateIndex: args.Kvpair.CreateIndex,
		ModifyIndex: args.Kvpair.ModifyIndex,
		LockIndex:   args.Kvpair.LockIndex,
		Flags:       args.Kvpair.Flags,
		Value:       args.Kvpair.Value,
		Session:     args.Kvpair.Session,
	}
	resp, err := s.impl.KVCas(ctx, &kvp)
	if nil != err {
		return nil, err
	}
	return &proto.KVCasReply{Success: resp}, nil
}
