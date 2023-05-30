package network

import (
    "net"
    "net/url"
    "net/http"
    "fmt"
    "io"
    "encoding/json"
    "time"
    "bytes"
    "encoding/gob"

    "github.com/lex0c/openet/pkg/config"
    "github.com/lex0c/openet/pkg/security"
)

type Discovery struct {
    gatewayUrl  url.URL
    security    security.Security
}

func (d *Discovery) Sync() (bool, error) {
    ip, err := d.GetMyIp()

    if err != nil {
        return false, fmt.Errorf("An error occurred: %v", err)
    }

    var addr string
    addr, err = security.GetMyAddress()

    if err != nil {
        return false, fmt.Errorf("An error occurred: %v", err)
    }

    t := time.Now()

    ni := NodeInfo{
        Ip: ip,
        Address: addr,
        UpdatedAt: t.String(),
        TTL: config.NodeIPCacheTTL,
    }

    msg := Message{
        Data: ni,
        Signature: "",
    }

    buf := new(bytes.Buffer)
    enc := gob.NewEncoder(buf)

    if err = enc.Encode(msg); err != nil {
        return false, fmt.Errorf("An error occurred: %v", err)
    }

    payload := buf.Bytes()

    req, err := http.NewRequest("POST", "", bytes.NewBuffer(payload))

    if err != nil {
        return false, fmt.Errorf("An error occurred: %v", err)
    }

    req.Header.Set("Content-Type", "application/octet-stream")

    client := &http.Client{}
    resp, err := client.Do(req)

    if err != nil {
        return false, fmt.Errorf("An error occurred: %v", err)
    }

    defer resp.Body.Close()

    return true, nil
}

func (d *Discovery) GetMyIp() (net.IP, error) {
    resp, err := http.Get(fmt.Sprintf("%s/%s", d.gatewayUrl, config.MyIPRoute))

    defer resp.Body.Close()

    if err != nil {
        return net.IP{}, fmt.Errorf("An error occurred: %v", err)
    }

    body, err := io.ReadAll(resp.Body)

    if err != nil {
        return net.IP{}, fmt.Errorf("An error occurred: %v", err)
    }

    var ip net.IP

    if err = json.Unmarshal(body, &ip); err != nil {
        return net.IP{}, fmt.Errorf("An error occurred: %v", err)
    }

    return ip, nil
}

func (d *Discovery) ResolveAddress(address string) (NodeInfo, error) {
    return NodeInfo{}, nil
}

func (d *Discovery) ResolveIp(ip net.IP) (NodeInfo, error) {
    return NodeInfo{}, nil
}

