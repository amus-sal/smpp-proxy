package proxy

import (
	"log"

	"../client"
	"../connection"
	"../server"

	"github.com/fiorix/go-smpp/smpp/pdu"
)

//Proxy ...
type Proxy struct {
	connection connection.SmppConn
}

//NewProxy ...
func NewProxy(connection connection.SmppConn) *Proxy {
	return &Proxy{
		connection,
	}
}

//RunProxy ...
func (proxy *Proxy) RunProxy() {
	log.Println("Start Proxy initiation for SMPP Session")
	serverRec := make(chan pdu.Body)
	serverSub := make(chan pdu.Body)
	clientRec := make(chan pdu.Body)
	clientSub := make(chan pdu.Body)

	smppServer := server.NewServer(proxy.connection, serverRec, serverSub)
	client := client.NewClient(clientRec, clientSub)
	go smppServer.RunServer()
	go client.RunClient()
	for {
		select {
		case serverPacket := <-serverRec:
			if serverPacket.Header().ID == pdu.UnbindRespID {
				log.Println("[Proxy]  Get From Server", serverPacket.Header().ID)
				clientSub <- serverPacket
				close(serverSub)
				close(clientSub)
				close(clientRec)
				close(serverRec)
				log.Println("[Proxy]  Send to client")
				break
			} else {
				clientSub <- serverPacket
				log.Println("[Proxy]  Send to client")
			}
		case clientPack := <-clientRec:
			log.Println("[Proxy]  Get From Client", clientPack.Header().ID)
			if clientPack.Header().ID == pdu.UnbindRespID {
				log.Println("[Proxy]  Send UNbind to Server")
				serverSub <- clientPack
				close(clientSub)
				close(serverSub)
				close(serverRec)
				close(clientRec)
				break
			} else {
				serverSub <- clientPack
				log.Println("[Proxy]  Send to Server")
			}
		default:
			log.Println("Idel process")
		}
	}

}
