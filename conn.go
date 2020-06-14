package pool

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Conn interface {
	GetC() interface{}
	Close() error

	SetUseTime(time.Time)
	GetUseTime() time.Time
}

type DbWrapConn struct {
	c *gorm.DB
	t time.Time
}

func NewDbWrapConn(conn *gorm.DB) *DbWrapConn {
	return &DbWrapConn{
		c: conn,
		t: time.Now(),
	}
}

func (conn *DbWrapConn) GetC() interface{} {
	return conn.c
}

func (conn *DbWrapConn) Close() error {
	dbconn := conn.c
	conn.c = nil
	return dbconn.Close()
}

func (conn *DbWrapConn) SetUseTime(t time.Time) {
	conn.t = t
}

func (conn *DbWrapConn) GetUseTime() time.Time {
	return conn.t
}
