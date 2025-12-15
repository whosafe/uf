package main

import (
	"fmt"
	"log"

	"iutime.com/utime/uf/uconfig"
	"iutime.com/utime/uf/uconv"
)

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string
	Port []int
}

func (s *ServerConfig) UnmarshalYAML(key string, value *uconfig.Node) error {

	switch key {
	case "host":
		s.Host = value.String()
	case "port":
		// 处理数组
		s.Port = make([]int, 0)
		if err := value.Iter(func(i int, v *uconfig.Node) error {
			s.Port = append(s.Port, uconv.ToIntDef(v, 0))
			return nil
		}); err != nil {
			return err
		}
	}
	return nil
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level string
	Path  string
}

func (l *LoggerConfig) UnmarshalYAML(key string, value *uconfig.Node) error {

	switch key {
	case "level":
		l.Level = value.String()
	case "path":
		l.Path = value.String()
	}
	return nil
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	DSN     string
	MaxOpen int
	Logger  LoggerConfig
}

func (d *DatabaseConfig) UnmarshalYAML(key string, value *uconfig.Node) error {

	switch key {
	case "dsn":
		d.DSN = value.String()
	case "max_open":
		d.MaxOpen = uconv.ToIntDef(value, 0)
	case "logger":
		// 解析嵌套结构
		if err := value.Decode(&d.Logger); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	var srvCfg ServerConfig
	var dbCfg DatabaseConfig

	// 注册配置解析
	uconfig.Register("server", srvCfg.UnmarshalYAML)
	uconfig.Register("database", dbCfg.UnmarshalYAML)

	// 加载配置
	fmt.Println("开始加载配置...")
	if err := uconfig.Load("example/uconfig/config.yaml"); err != nil {
		log.Fatalf("加载配置失败: %+v", err)
	}

	fmt.Printf("Server Config: %+v\n", srvCfg)
	fmt.Printf("Database Config: %+v\n", dbCfg)
	fmt.Println("配置加载成功")
}
