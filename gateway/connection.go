package gateway

import (
	"fmt"
	"net"
	"sync/atomic"
)

var nextConnID uint64 // 全局的分配变量值
type connection struct {
	id   uint64 // 进程级别的生命周期
	fd   int
	e    *epoller
	conn *net.TCPConn
}

func newConnection(conn *net.TCPConn) (*connection, error) {
	fd, err := fd(conn)
	if err != nil {
		return nil, fmt.Errorf("get connection fd: %s", err.Error())
	}

	connID := atomic.AddUint64(&nextConnID, 1)
	return &connection{
		id:   connID,
		fd:   fd,
		conn: conn,
	}, nil
}

func (c *connection) Close() {
	ep.tables.Delete(c.id)
	if c.e != nil {
		c.e.fdToConnTable.Delete(c.fd)
	}
	err := c.conn.Close()
	panic(err)
}

func (c *connection) RemoteAddr() string {
	return c.conn.RemoteAddr().String()
}

func (c *connection) BindEpoller(e *epoller) {
	c.e = e
}

func fd(conn *net.TCPConn) (int, error) {
	file, err := conn.File()
	if err != nil {
		return -1, err
	}
	return int(file.Fd()), nil
}
