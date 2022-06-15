// Code generated by goctl. DO NOT EDIT!
// Source: trans.proto

package transclient

import (
	"context"

	"msg/trans/trans"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	AdjustInfo = trans.AdjustInfo
	Response   = trans.Response

	Trans interface {
		TransOut(ctx context.Context, in *AdjustInfo, opts ...grpc.CallOption) (*Response, error)
		TransIn(ctx context.Context, in *AdjustInfo, opts ...grpc.CallOption) (*Response, error)
	}

	defaultTrans struct {
		cli zrpc.Client
	}
)

func NewTrans(cli zrpc.Client) Trans {
	return &defaultTrans{
		cli: cli,
	}
}

func (m *defaultTrans) TransOut(ctx context.Context, in *AdjustInfo, opts ...grpc.CallOption) (*Response, error) {
	client := trans.NewTransClient(m.cli.Conn())
	return client.TransOut(ctx, in, opts...)
}

func (m *defaultTrans) TransIn(ctx context.Context, in *AdjustInfo, opts ...grpc.CallOption) (*Response, error) {
	client := trans.NewTransClient(m.cli.Conn())
	return client.TransIn(ctx, in, opts...)
}
