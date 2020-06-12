package connection

import (
	"bufio"
	"net"
)

//SmppConn ..
type SmppConn struct {
	Rwc net.Conn
	R   *bufio.Reader
	W   *bufio.Writer
}
