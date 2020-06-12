package server

import (
	"bytes"
	"io"
	"log"

	"../connection"

	"github.com/fiorix/go-smpp/smpp/pdu"
)

type (
	//SmppServer ...
	SmppServer struct {
		connection connection.SmppConn
		ServerSub  chan<- pdu.Body
		ServerRec  <-chan pdu.Body
	}
)

//NewServer ...
func NewServer(connection connection.SmppConn, ServerSub chan<- pdu.Body, ServerRec <-chan pdu.Body) *SmppServer {
	return &SmppServer{
		connection,
		ServerSub,
		ServerRec,
	}
}

//RunServer ...
func (server *SmppServer) RunServer() {
	log.Println("Start Server initiation for SMPP Session")
	for {

		packet, err := pdu.Decode(server.connection.R)
		if err != nil {
			log.Println("[Server] Get From Client Close Error EOF")
			println(err)
			// TODO add Handler if connection down to send to global channel to release proxy
			break
		}
		log.Println("[Server] Get From Client")
		server.ServerSub <- packet
		log.Println("[Server] Send to Proxy")

		proxyPacket, ok := <-server.ServerRec
		if !ok {
			log.Println("[Server] Get Channel Close  from Proxy")
			break
		} else {
			log.Println("[Server] Get From Proxy")
			var reader bytes.Buffer
			proxyPacket.SerializeTo(&reader)
			io.Copy(server.connection.W, &reader)
			log.Println("[Server] Send to Client")
			server.connection.W.Flush()
		}

	}
}
