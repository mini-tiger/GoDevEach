package services

import (
	"context"
	"sync"
)

/**
 * @Author: Tao Jun
 * @Since: 2023/3/20
 * @Desc: entry.go
**/

type RegSrv struct {
	StopSrvFn  func(ctx context.Context, wg *sync.WaitGroup)
	StartSrvFn func(ctx context.Context)
}

func (r *RegSrv) RegStartSrv(fn func(ctx context.Context)) {
	r.StartSrvFn = fn
}

func (r *RegSrv) StartSrv(ctx context.Context) {
	r.StartSrvFn(ctx)
}

func (r *RegSrv) RegStopSrv(fn func(ctx context.Context, wg *sync.WaitGroup)) {
	r.StopSrvFn = fn
}

func (r *RegSrv) StopSrv(ctx context.Context, wg *sync.WaitGroup) {
	r.StopSrvFn(ctx, wg)
}

var _ RegSrvInterface = new(RegSrv)

type RegSrvInterface interface {
	StopSrv(ctx context.Context, wg *sync.WaitGroup)
	StartSrv(ctx context.Context)
}
