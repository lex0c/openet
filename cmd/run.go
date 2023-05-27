package main

import (
    "flag"
    "fmt"
    "log"
    "net"
    "strings"
    "github.com/lex0c/openet/pkg/network"
    "github.com/lex0c/openet/pkg/config"
)

func main() {
    var incomingPort int
    var outgoingPorts string

    flag.IntVar(&incomingPort, "incoming", 8080, "Port for incoming connections")
    flag.StringVar(&outgoingPorts, "outgoing", "", "Comma-separated list of ports for outgoing addresses (IP:Port)")
    flag.Parse()

    outgoingPortList := strings.Split(outgoingPorts, ",")

    incomingPool := network.NewPool(nil, config.MaxIncomingConnections)
    outgoingPool := network.NewPool(outgoingPortList, config.MaxOutgoingConnections)

    ln, err := net.Listen("tcp", fmt.Sprintf(":%d", incomingPort))

    if err != nil {
        log.Println(err)
    } else {
        log.Println("Listening on port", incomingPort)

        go func() {
            for {
                conn, err := ln.Accept()

                if err != nil {
                    log.Println("Failed to accept incoming connection:", err)
                    continue
                }

                err = incomingPool.Add(conn)

                if err != nil {
                    log.Println("Failed to add connection:", conn.RemoteAddr(), " | ", err)
                    conn.Close()
                    continue
                }

                go network.HandleConnection(outgoingPool, conn, func(message string) {
                    fmt.Println("Received: ", message)
                })
			      }
        }()
    }

    for _, conn := range outgoingPool.ListConnections() {
        go network.HandleConnection(incomingPool, conn, func(message string) {
            log.Println("Received: ", message)
        })
    }

    select {}
}
