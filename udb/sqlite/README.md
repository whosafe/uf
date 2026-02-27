# udb/sqlite - SQLite 数据库层

基于标准库 `database/sql` 和纯 Go 实现 `modernc.org/sqlite` 的零反射、高性能 SQLite 数据库层。

该模块提供了与 `udb/postgresql` 完全一致的查询构建器、事务和 CRUD 操作 API，方便项目在 PostgreSQL 与 SQLite 间无缝切换。

## ✨ 核心特性

- **纯 Go 实现**：无需 CGO，无需在 Windows 编译环境等繁琐配置，直接开箱即用。
- **与 PostgreSQL 无缝兼容**：构建器及其链式调用方式和接口与 `udb/postgresql` 保持结构一致（移除了 `$1, $2` 替换逻辑，直接使用原生的 `?`）。
- **零反射设计**：通过 `Scanner` 接口提取数据，提供原生级别的性能体验。
- **链路追踪**：继承 `ulogger` 追踪系统记录日志和监控慢查询。

## 📦 安装

使用以下命令获取：

```bash
go get github.com/whosafe/uf/udb/sqlite
```

同时请确保你安装了对应依赖：

```bash
go get modernc.org/sqlite
```

## 🚀 快速开始

### 1. 配置文件

SQLite 不再需要主机、端口、用户名和密码，通常只需配置路径即可。这在轻量级部署和测试环境中非常有用。
这里是一个简单的连接配置：

```yaml
database:
  sqlite:
    # 数据库路径 (如需在内存建立数据库可以使用 ":memory:")
    path: "./data.db"
    mode: "rwc"  # ro (读), rw (读写), rwc (读写且不存在则创建), memory
    
    # 连接池配置
    pool:
      max_conns: 10
      max_conn_lifetime: "1h"
      max_conn_idle_time: "30m"
    
    # 查询配置
    query:
      default_timeout: "30s"
      slow_query_threshold: "1s"
    
    # 日志配置
    log:
      enabled: true
      level: "info"
      format: "text"
      output: "stdout"
      slow_query: true
      log_params: false
```

### 2. 使用方法

与 `postgresql` 相比，创建连接实例的包名改为 `sqlite`，其他后续的 `queryBuilder`、操作和事务等用法完全保持一致。

```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/whosafe/uf/uconfig"
    "github.com/whosafe/uf/ucontext"
    "github.com/whosafe/uf/udb/sqlite"
    "github.com/whosafe/uf/uconv"
)

type User struct {
    ID        int64
    Username  string
    Email     string
    Age       int
    CreatedAt time.Time
}

// Scan 实现 Scanner 接口
func (u *User) Scan(key string, value any) error {
    switch key {
    case "id":
        u.ID = uconv.ToInt64Def(value, 0)
    case "username":
        u.Username = uconv.ToString(value)
    case "email":
        u.Email = uconv.ToString(value)
    case "age":
        u.Age = uconv.ToIntDef(value, 0)
    case "created_at":
        u.CreatedAt = uconv.ToTimeDef(value, time.Time{})
    }
    return nil
}

func main() {
    uconfig.Load("config.yaml")

    // 创建 SQLite 连接
    conn, err := sqlite.New(sqlite.DefaultConfig())
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    ctx := ucontext.NewContext(context.Background())

    // 查询单条记录
    var user User
    err = conn.Query(ctx).
        Table("users").
        Where("id = ?", 1).
        Scan(&user)

    if err != nil {
        if err == sqlite.ErrNoRows {
            log.Println("用户不存在")
        } else {
            log.Fatal(err)
        }
    }

    log.Printf("用户: %+v\n", user)
}
```

关于所有的 API `Select`, `Where`, `OrderBy`, `GroupBy` 以及 `Insert`, `Update`, `Delete` 与事务的使用方式，详见 `postgresql` 的文档，所有的用法如出一辙！
