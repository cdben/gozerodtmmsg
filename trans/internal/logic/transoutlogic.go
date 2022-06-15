package logic

import (
	"context"
	"database/sql"
	"github.com/dtm-labs/dtmcli"
	"github.com/dtm-labs/dtmgrpc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"msg/trans/internal/svc"
	"msg/trans/trans"

	"github.com/zeromicro/go-zero/core/logx"
)

type TransOutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTransOutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TransOutLogic {
	return &TransOutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TransOutLogic) TransOut(in *trans.AdjustInfo) (*trans.Response, error) {
	db, err := sqlx.NewMysql(l.svcCtx.Config.Mysql.DataSource).RawDB()
	if err != nil {
		return nil, status.Error(codes.Aborted, dtmcli.ResultFailure)
	}
	// 获取子事务屏障
	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	if err != nil {
		return nil, status.Error(codes.Aborted, dtmcli.ResultFailure)
	}
	err = barrier.CallWithDB(db, func(tx *sql.Tx) error {
		// 更新余额
		result, err := l.svcCtx.UserAccountModel.TxAdjustBalance(tx, in.UserID, -in.Amount)
		if err != nil {
			return err
		}
		affected, err := result.RowsAffected()
		// 余额不足，返回子事务失败
		if err == nil && affected == 0 {
			return dtmcli.ErrFailure
		}
		return err
	})

	// 这种情况是余额不足，不再重试，走回滚
	if err == dtmcli.ErrFailure {
		return nil, status.Error(codes.Aborted, dtmcli.ResultFailure)
	}
	return &trans.Response{}, nil
}
