package grpcx

import (
	"context"
	"time"

	"github.com/yakaa/grpcx/config"
	"github.com/yakaa/grpcx/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	rpcResolver "google.golang.org/grpc/resolver"
)

type (
	GrpcxClient struct {
		resolver *resolver.Resolver
		timeOut  time.Duration
	}
)

func MustNewGrpcxClient(conf *config.ClientConf) (*GrpcxClient, error) {
	r, err := resolver.NewResolver(conf)  // xxx resolver的 build方法 监控etcd
	if err != nil {
		return nil, err
	}
	rpcResolver.Register(r)

	return &GrpcxClient{resolver: r, timeOut: config.GrpcxDialTimeout}, nil
}

func (c *GrpcxClient) GetConnection(options ...grpc.DialOption) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBalancerName(roundrobin.Name)} //随机

	if c.resolver.WithBlock() {
		opts = append(opts, grpc.WithBlock())
	}
	opts = append(opts, options...)
	ctx, _ := context.WithTimeout(context.Background(), c.timeOut)
	return grpc.DialContext(ctx, c.resolver.Target(), opts...)
}

// SetTimout ctx content
func (c *GrpcxClient) SetTimeOut(timeOut time.Duration) {
	c.timeOut = timeOut
}
