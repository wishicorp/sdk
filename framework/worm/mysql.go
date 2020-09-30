//worm 数据库连接配置
package worm

import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/hashicorp/go-hclog"
	"log"
	"xorm.io/core"
)

//数据库连接配置
type Config struct {
	ShowSql        bool     `hcl:"show_sql"`
	MaxIdle        int      `hcl:"max_idle"`
	MaxConn        int      `hcl:"max_conn"`
	Master         string   `hcl:"master"`
	Slaves         []string `hcl:"slaves"`
	UseMasterSlave bool     `hcl:"use_master_slave"`
}

type DBInterface interface {
	xorm.EngineInterface
}

func NewConn(c *Config, logger hclog.Logger) (DBInterface, error) {
	if c.UseMasterSlave {
		return NewMSConn(c, logger)
	}
	return NewSingleConn(c, logger)
}

//初始化数据库连接
//root:123456@tcp(127.0.0.1)/insight_reader?charset=utf8
func NewSingleConn(c *Config, logger hclog.Logger) (DBInterface, error) {
	if nil == c || "" == c.Master {
		return nil, errors.New("config or config.Url can not be null")
	}

	conn, err := xorm.NewEngine("mysql", c.Master)
	if nil != err || nil == conn {
		log.Println("failed to initializing db connection:", err)
		return nil, err
	}
	conn.SetLogger(NewSQLLogger(logger, c.ShowSql))
	conn.SetMapper(core.GonicMapper{})
	conn.ShowSQL(c.ShowSql)
	conn.SetLogLevel(core.LOG_INFO)
	conn.SetMaxIdleConns(c.MaxIdle)
	conn.SetMaxOpenConns(c.MaxConn)
	return conn, nil
}

//初始化主从数据库连接, master不能为空，slaves可以为空
func NewMSConn(c *Config, logger hclog.Logger) (DBInterface, error) {
	if nil == c || "" == c.Master {
		return nil, errors.New("config or config.Url can not be null")
	}
	conns := make([]string, len(c.Slaves)+1)
	conns[0] = c.Master
	for i, v := range c.Slaves {
		conns[i+1] = v
		if "" == v {
			return nil, errors.New("config or config.Url can not be null")
		}
	}

	group, err := xorm.NewEngineGroup("xorm", conns)

	if nil != err || nil == group {
		log.Printf("failed to initializing db connection: %s\n", err)
		return nil, err
	}

	group.SetLogger(NewSQLLogger(logger, c.ShowSql))
	group.SetMapper(core.GonicMapper{})
	group.ShowSQL(c.ShowSql)
	group.SetMaxIdleConns(c.MaxIdle)
	group.SetMaxOpenConns(c.MaxConn)
	return group, nil
}
