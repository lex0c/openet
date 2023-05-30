package network

import (
    "net"
)

type Message struct {
    Data        interface{}
    Signature   string
}

type NodeInfo struct {
    Ip          net.IP
    Address     string
    UpdatedAt   string
    TTL         int
}

