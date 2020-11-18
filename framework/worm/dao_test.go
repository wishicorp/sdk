package worm

import (
	"github.com/hashicorp/go-hclog"
	"testing"
)

type TestTrans struct {
	Id int `json:"id" xorm:"id"`
}

func (m *TestTrans)TableName()string  {
	return "test_trans"
}
func TestEntity_Create(t *testing.T) {
	logger := hclog.Default()
	c := Config{
		ShowSql:        true,
		Master:         "root:123456@tcp(127.0.0.1)/virtual_coin?charset=utf8mb4&parseTime=true&loc=Local",
	}
	conn, _ := NewConn(&c, logger)
	createTable := `create table if not exists test_trans(id int)engine=innodb`;
	conn.Query(createTable)
	dao := NewBaseDao(conn)
	session := dao.NewSession()
	defer session.Close()


	t1 := TestTrans{Id: 1}
	t2 := TestTrans{Id: 2}
	t3 := TestTrans{Id: 3}
	dao.InsertOne(&t1)
	session.Begin()
	session.InsertOne(&t2)
	if _, err := dao.InsertOne(&t3); err != nil {
		t.Fatal(err)
		return
	}
}
