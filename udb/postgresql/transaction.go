package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/whosafe/uf/uerror"
)

// Transaction 事务
type Transaction struct {
	tx         pgx.Tx
	ctx        context.Context
	Connection *Connection
}

// Query 创建查询构建器
func (t *Transaction) Query(ctx context.Context) *TxQueryBuilder {
	return &TxQueryBuilder{
		ctx:    ctx,
		tx:     t.tx,
		config: t.Connection.config,
	}
}

// Insert 创建插入构建器(事务中)
func (t *Transaction) Insert(ctx context.Context) *TxInsertBuilder {
	return &TxInsertBuilder{
		ctx: ctx,
		tx:  t.tx,
	}
}

// Update 创建更新构建器(事务中)
func (t *Transaction) Update(ctx context.Context) *TxUpdateBuilder {
	return &TxUpdateBuilder{
		ctx: ctx,
		tx:  t.tx,
	}
}

// Delete 创建删除构建器(事务中)
func (t *Transaction) Delete(ctx context.Context) *TxDeleteBuilder {
	return &TxDeleteBuilder{
		ctx: ctx,
		tx:  t.tx,
	}
}

// Exec 执行 SQL
func (t *Transaction) Exec(ctx context.Context, sql string, args ...any) (int64, error) {
	result, err := t.tx.Exec(ctx, sql, args...)
	if err != nil {
		return 0, uerror.Wrap(err, "事务中执行 SQL 失败")
	}
	return result.RowsAffected(), nil
}

// Commit 提交事务
func (t *Transaction) Commit() error {
	if t.Connection.config.Log != nil && t.Connection.config.Log.Enabled {
		t.Connection.logger.InfoCtx(t.ctx, "提交事务")
	}

	err := t.tx.Commit(t.ctx)
	if err != nil {
		if t.Connection.config.Log != nil && t.Connection.config.Log.Enabled {
			t.Connection.logger.ErrorCtx(t.ctx, "提交事务失败", "error", err.Error())
		}
		return uerror.Wrap(err, "提交事务失败")
	}
	return nil
}

// Rollback 回滚事务
func (t *Transaction) Rollback() error {
	if t.Connection.config.Log != nil && t.Connection.config.Log.Enabled {
		t.Connection.logger.WarnCtx(t.ctx, "回滚事务")
	}

	err := t.tx.Rollback(t.ctx)
	if err != nil {
		if t.Connection.config.Log != nil && t.Connection.config.Log.Enabled {
			t.Connection.logger.ErrorCtx(t.ctx, "回滚事务失败", "error", err.Error())
		}
		return uerror.Wrap(err, "回滚事务失败")
	}
	return nil
}
