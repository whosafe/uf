package sqlite

import (
	"time"

	"github.com/whosafe/uf/uerror"
)

// Scan 扫描单行(零反射)
func (q *QueryBuilder) Scan(dest Scanner) error {
	sqlStr, args := q.BuildSQL()

	// 记录日志
	startTime := time.Now()

	if q.Connection.config.Log != nil && q.Connection.config.Log.Enabled {
		if q.Connection.config.Log.LogParams {
			q.Connection.logger.DebugCtx(q.ctx, "执行查询", "sql", sqlStr, "args", args)
		} else {
			q.Connection.logger.DebugCtx(q.ctx, "执行查询", "sql", sqlStr)
		}
	}

	// 执行查询
	rows, err := q.Connection.db.QueryContext(q.ctx, sqlStr, args...)
	if err != nil {
		// 记录错误日志
		if q.Connection.config.Log != nil && q.Connection.config.Log.Enabled {
			q.Connection.logger.ErrorCtx(q.ctx, "查询失败", "sql", sqlStr, "error", err.Error())
		}
		return uerror.Wrap(err, "查询失败")
	}
	defer rows.Close()

	if !rows.Next() {
		return ErrNoRows
	}

	// 获取列名
	columns, err := rows.Columns()
	if err != nil {
		return uerror.Wrap(err, "获取列名失败")
	}

	// 获取值
	values := make([]any, len(columns))
	valuePtrs := make([]any, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	if err := rows.Scan(valuePtrs...); err != nil {
		return uerror.Wrap(err, "获取行数据失败")
	}

	// 调用用户定义的 Scan 方法 - 逐个字段调用
	for i, col := range columns {
		if err := dest.Scan(col, values[i]); err != nil {
			return uerror.Wrap(err, "扫描字段失败")
		}
	}

	// 计算执行时间
	duration := time.Since(startTime)

	// 记录成功日志
	if q.Connection.config.Log != nil && q.Connection.config.Log.Enabled {
		q.Connection.logger.InfoCtx(q.ctx, "查询成功", "duration", duration)
	}

	// 记录慢查询
	if q.Connection.config.Log != nil && q.Connection.config.Log.SlowQuery {
		if duration > q.Connection.config.Query.SlowQueryThreshold {
			q.Connection.logger.WarnCtx(q.ctx, "慢查询", "sql", sqlStr, "duration", duration)
		}
	}

	return nil
}

// ScanAll 扫描多行
func (q *QueryBuilder) ScanAll(newScanner func() Scanner) ([]Scanner, error) {
	sqlStr, args := q.BuildSQL()

	// 记录日志
	startTime := time.Now()

	if q.Connection.config.Log != nil && q.Connection.config.Log.Enabled {
		if q.Connection.config.Log.LogParams {
			q.Connection.logger.DebugCtx(q.ctx, "执行查询", "sql", sqlStr, "args", args)
		} else {
			q.Connection.logger.DebugCtx(q.ctx, "执行查询", "sql", sqlStr)
		}
	}

	// 执行查询
	rows, err := q.Connection.db.QueryContext(q.ctx, sqlStr, args...)
	if err != nil {
		// 记录错误日志
		if q.Connection.config.Log != nil && q.Connection.config.Log.Enabled {
			q.Connection.logger.ErrorCtx(q.ctx, "批量查询失败", "sql", sqlStr, "error", err.Error())
		}
		return nil, uerror.Wrap(err, "批量查询失败")
	}
	defer rows.Close()

	// 获取列名
	columns, err := rows.Columns()
	if err != nil {
		return nil, uerror.Wrap(err, "获取列名失败")
	}

	var results []Scanner
	for rows.Next() {
		values := make([]any, len(columns))
		valuePtrs := make([]any, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, uerror.Wrap(err, "获取行数据失败")
		}

		scanner := newScanner()
		// 逐个字段调用 Scan
		for i, col := range columns {
			if err := scanner.Scan(col, values[i]); err != nil {
				return nil, uerror.Wrap(err, "扫描字段失败")
			}
		}
		results = append(results, scanner)
	}

	// 计算执行时间
	duration := time.Since(startTime)

	// 记录成功日志
	if q.Connection.config.Log != nil && q.Connection.config.Log.Enabled {
		q.Connection.logger.InfoCtx(q.ctx, "批量查询成功", "count", len(results), "duration", duration)
	}

	// 记录慢查询
	if q.Connection.config.Log != nil && q.Connection.config.Log.SlowQuery {
		if duration > q.Connection.config.Query.SlowQueryThreshold {
			q.Connection.logger.WarnCtx(q.ctx, "慢查询", "sql", sqlStr, "duration", duration, "count", len(results))
		}
	}

	if err := rows.Err(); err != nil {
		return results, uerror.Wrap(err, "遍历结果集失败")
	}
	return results, nil
}
