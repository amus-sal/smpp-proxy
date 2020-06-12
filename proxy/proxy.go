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
		case serverPacket, ok := <-serverRec:
			if ok {
				log.Println("[Proxy]  Get From Server", serverPacket.Header().ID)
				select {
				case clientSub <- serverPacket:
				default:
				}
				log.Println("[Proxy]  Send to client")
				if serverPacket.Header().ID == pdu.UnbindRespID {
					close(serverSub)
					close(clientSub)
					close(clientRec)
					close(serverRec)
					break
				}
			}
		case clientPack, ok := <-clientRec:
			if ok {
				log.Println("[Proxy]  Get From Client", clientPack.Header().ID)
				select {
				case serverSub <- clientPack:
				default:
				}
				log.Println("[Proxy]  Send to Server")
				if clientPack.Header().ID == pdu.UnbindRespID {
					close(clientSub)
					close(serverSub)
					close(serverRec)
					close(clientRec)
					break
				}
			}

		default:
			// log.Println("Idle process")
		}
	}

}
