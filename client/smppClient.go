package client

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/fiorix/go-smpp/smpp/pdu"
)

type (
	//SmppClient ...
	SmppClient struct {
		ClientSub chan<- pdu.Body
		ClientRec <-chan pdu.Body
	}
)

//NewClient ...
func NewClient(ClientSub chan<- pdu.Body, ClientRec <-chan pdu.Body) *SmppClient {
	return &SmppClient{
		ClientSub,
		ClientRec,
	}
}

//RunClient ...
func (client *SmppClient) RunClient() {
	log.Println("Start Client  initiation for SMPP Session")
	address := os.Getenv("SERVER_ADDRESS")
	if address == "" {
		log.Fatal("There ia a problem with Server Address")
	}
	c, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println(err)
		// TODO add Handler if connection dowb to send to global channel to release proxy
		return
	}
	defer c.Close()

	for {

		packet, ok := <-client.ClientRec
		log.Println("[Client] Get From Proxy")
		if !ok {
			log.Println("[Client] Get From Channel Close Proxy")
			break
		}

		var b bytes.Buffer
		packet.SerializeTo(&b)

		writer := bufio.NewWriter(c)
		io.Copy(writer, &b)
		writer.Flush()
		log.Println("[Client]  Send to Operator")

		reader := bufio.NewReader(c)
		log.Println("[Client] Get From Operator")
		resp, _ := pdu.Decode(reader)
		if resp == nil {
			// TODO add Handler if connection dowb to send to global channel to release proxy
			break
		}
		client.ClientSub <- resp
		log.Println("[Client] Send to proxy")
	}

}
