package conn

import "net"

type Conn struct {
	Id uint32
	net.Conn
}
