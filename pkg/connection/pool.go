package connection

import (
    "encoding/gob"
    "fmt"
    "log"
    "net"
    "sync"
)

type Pool struct {
    connections []net.Conn
    mx          sync.Mutex
}

func (p *Pool) Add(conn net.Conn) error {
    p.mx.Lock()
    defer p.mx.Unlock()

    if len(p.connections) >= 4 {
        return fmt.Errorf("Connection pool is full")
    }

    p.connections = append(p.connections, conn)
    log.Println("Added new connection", conn.RemoteAddr())
    return nil
}

func (p *Pool) Remove(conn net.Conn) {
    p.mx.Lock()
    defer p.mx.Unlock()

    for i, c := range p.connections {
        if c == conn {
            p.connections = append(p.connections[:i], p.connections[i+1:]...)
            log.Println("Removed connection", conn.RemoteAddr())
            return
        }
    }
}

func (p *Pool) Broadcast(msg Message) {
    p.mx.Lock()
    defer p.mx.Unlock()

    for _, conn := range p.connections {
        encoder := gob.NewEncoder(conn)

        if err := encoder.Encode(msg); err != nil {
            log.Println(err)
            continue
        }

		    log.Println("Broadcasted message", msg.Kind, "to:", conn.RemoteAddr())
    }
}

func (p *Pool) ListConnections() []net.Conn {
    return p.connections
}

func NewPool(peers []string) *Pool {
    pool := &Pool{
        connections: make([]net.Conn, 0),
    }

    for _, peer := range peers {
        conn, err := net.Dial("tcp", peer)

        if err != nil {
            log.Println(err)
        } else {
            pool.connections = append(pool.connections, conn)
            log.Println("Connected to", peer)
        }
    }

    return pool
}

func HandleConnection(pool *Pool, conn net.Conn, callback func(message string)) {
    decoder := gob.NewDecoder(conn)

    for {
        var msg Message

        if err := decoder.Decode(&msg); err != nil {
            log.Println("Failed to decode message from", conn.RemoteAddr(), " | ", err)
            pool.Remove(conn)
            conn.Close()
            return
        }

        callback(msg.Kind)
        pool.Broadcast(msg)
    }
}

