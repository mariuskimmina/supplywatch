package warehouse

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/mariuskimmina/supplywatch/pkg/logger"
)

type warehouse struct {

}

func NewWarehouse() *warehouse {
    return &warehouse{}
}

var (
    address = "0.0.0.0:1234"
)

const (
    maxBufferSize = 1024
)

func (w *warehouse) Start() {
    ctx := context.Background()
    logger.Debug("Starting warehouse")
    listen, err := net.ListenPacket("udp", address)
    if err != nil {
        return
    }

    defer listen.Close()

    doneChan := make(chan error, 1)
    buffer := make([]byte, maxBufferSize)

    go func() {
        for {
            n, addr, err := listen.ReadFrom(buffer)
            if err != nil {
                doneChan <- err
                return
            }

            fmt.Printf("packet-received: bytes=%d from=%s\n",
            n, addr.String())

            writeDeadline := time.Now().Add(15 * time.Second)
            err = listen.SetWriteDeadline(writeDeadline)
            if err != nil {
                doneChan <- err
                return
            }
            readDeadline := time.Now().Add(15 * time.Second)
            err = listen.SetReadDeadline(readDeadline)
            if err != nil {
                doneChan <- err
                return
            }

            n, err = listen.WriteTo(buffer[:n], addr)
            if err != nil {
                doneChan <- err
                return
            }

            fmt.Printf("packet-written: bytes=%d to=%s\n", n, addr.String())
        }
    }()

    select {
    case <-ctx.Done():
        fmt.Println("cancelled")
        err = ctx.Err()
    case err = <-doneChan:
    }

    return
}
