// Code generated by goctl. DO NOT EDIT!
// Source: trans.proto

package server

import (
	"context"

	"msg/trans/internal/logic"
	"msg/trans/internal/svc"
	"msg/trans/trans"
)

type TransServer struct {
	svcCtx *svc.ServiceContext
	trans.UnimplementedTransServer
}

func NewTransServer(svcCtx *svc.ServiceContext) *TransServer {
	return &TransServer{
		svcCtx: svcCtx,
	}
}

func (s *TransServer) TransOut(ctx context.Context, in *trans.AdjustInfo) (*trans.Response, error) {
	l := logic.NewTransOutLogic(ctx, s.svcCtx)
	return l.TransOut(in)
}

func (s *TransServer) TransIn(ctx context.Context, in *trans.AdjustInfo) (*trans.Response, error) {
	l := logic.NewTransInLogic(ctx, s.svcCtx)
	return l.TransIn(in)
}
