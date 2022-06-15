package logic

import (
	"context"
	"github.com/dtm-labs/dtmgrpc"
	"google.golang.org/grpc/status"
	"msg/trans/transclient"

	"msg/api/internal/svc"
	"msg/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TransLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTransLogic(ctx context.Context, svcCtx *svc.ServiceContext) TransLogic {
	return TransLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TransLogic) Trans(req types.TransRequest) (*types.TransResponse, error) {
	transRpcBusiServer, err := l.svcCtx.Config.TransRpc.BuildTarget()
	if err != nil {
		return nil, status.Error(100, "程序异常")
	}

	// dtm 服务的 etcd 注册地址
	var dtmServer = "etcd://localhost:2379/dtmservice"

	// 创建一个gid
	gid := dtmgrpc.MustGenGid(dtmServer)

	msg := dtmgrpc.NewMsgGrpc(dtmServer, gid).
		Add(transRpcBusiServer+"/transclient.Trans/TransOut", &transclient.AdjustInfo{
			UserID: req.UserId,
			Amount: req.Amount,
		}).
		Add(transRpcBusiServer+"/transclient.Trans/TransIn", &transclient.AdjustInfo{
			UserID: req.ToUserId,
			Amount: req.Amount,
		})

	if err = msg.Submit(); err != nil {
		return nil, err
	}

	return &types.TransResponse{}, nil
}
