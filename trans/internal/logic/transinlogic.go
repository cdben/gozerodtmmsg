package logic

import (
	"context"
	"database/sql"
	"github.com/dtm-labs/dtmgrpc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"msg/trans/internal/svc"
	"msg/trans/trans"

	"github.com/zeromicro/go-zero/core/logx"
	status "google.golang.org/grpc/status"
)

type TransInLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTransInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TransInLogic {
	return &TransInLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TransInLogic) TransIn(in *trans.AdjustInfo) (*trans.Response, error) {

	db, err := sqlx.NewMysql(l.svcCtx.Config.Mysql.DataSource).RawDB()
	if err != nil {
		return nil, status.Error(500, err.Error())
	}
	// 获取子事务屏障
	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	if err != nil {
		return nil, status.Error(500, err.Error())
	}
	err = barrier.CallWithDB(db, func(tx *sql.Tx) error {
		// 更新余额
		_, err := l.svcCtx.UserAccountModel.TxAdjustBalance(tx, in.UserID, in.Amount)
		return err
	})

	return &trans.Response{}, nil
}
