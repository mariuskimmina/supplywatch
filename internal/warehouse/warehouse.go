package warehouse

import (
	"fmt"
	"net"
	"os"

	log "github.com/mariuskimmina/supplywatch/pkg/log"
)

type warehouse struct {
	logger *log.Logger
}

func NewWarehouse(logger *log.Logger) *warehouse {
	return &warehouse{
		logger: logger,
	}
}

var (
	address = net.UDPAddr{
		Port: 4444,
		IP:   net.ParseIP("0.0.0.0"),
	}
)

const (
	maxBufferSize = 1024
)

func (w *warehouse) Start() {
	//ctx := context.Background()
	w.logger.Info("Warehouse Starting")
	listen, err := net.ListenUDP("udp", &address)
	if err != nil {
		return
	}
	defer listen.Close()
	go w.recvDataFromSensor(listen)
	ln, err := net.Listen("tcp", "0.0.0.0:8000")
	if err != nil {
		w.logger.Error(err.Error())
		return
	}
	defer ln.Close()
	for {
		c, err := ln.Accept()
		if err != nil {
			w.logger.Error(err.Error())
			return
		}
		go w.handleConnection(c)
	}
}


func (w *warehouse) recvDataFromSensor(listen *net.UDPConn) {
    f, err := os.Create("/tmp/log2")
    if err != nil {
        return
    }
    defer f.Close()
	for {
		p := make([]byte, maxBufferSize)
		_, remoteaddr, err := listen.ReadFromUDP(p)
		if err != nil {
			// logger.Error(err.Error())
			return
		}
		fmt.Printf("Read a message from %v %s \n", remoteaddr, p)
        f.Write(p)
        f.Write([]byte("\n"))
	}
}

