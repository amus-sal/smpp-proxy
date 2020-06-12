package main

import (
	"bufio"
	"fmt"
	"net"

	"./connection"

	"./proxy"
)

func main() {
	listner, err := net.Listen("tcp", "127.0.0.1:1234")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listner.Close()
	for {
		c, err := listner.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		smppCon := connection.SmppConn{c, bufio.NewReader(c), bufio.NewWriter(c)}

		proxy := proxy.NewProxy(smppCon)
		go proxy.RunProxy()
	}

}
