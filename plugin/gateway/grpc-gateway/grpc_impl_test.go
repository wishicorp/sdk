package grpc_gateway

import (
	"context"
	gwproto "github.com/wishicorp/sdk/plugin/gateway/grpc-gateway/proto"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestNewGateway(t *testing.T) {
	conn, err := grpc.Dial("127.0.0.1:18401", grpc.WithInsecure())
	if nil != err {
		t.Fatal(err)
	}
	defer conn.Close()
	cli := gwproto.NewRpcGatewayClient(conn)
	args := gwproto.RequestArgs{
		Backend:   "admin",
		Namespace: "order",
		Operation: "list",
		Token:     "d6cafbff-e1fa-4669-91ac-d32e3ee18985",
		Data: []byte(`{"size":1,"page":1}`),
	}
	time.Sleep(time.Second)
	reply, err := cli.ExecRequest(context.Background(), &args)
	if nil != err {
		t.Fatal(err)
	}
	t.Log(reply)
}
