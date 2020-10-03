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
		Backend:   "account",
		Namespace: "user",
		Operation: "home",
		Token:     "062ad5fc-c643-44e6-a916-b29f8bc4384c",
		Data:      []byte{},
	}
	time.Sleep(time.Second)
	reply, err := cli.ExecRequest(context.Background(), &args)
	if nil != err {
		t.Fatal(err)
	}
	t.Log(string(reply.Result.Data))
}
