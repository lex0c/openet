package network

import (
    "net"
    "net/url"
    "net/http"
    "fmt"
    "io"
    "encoding/json"

    "github.com/lex0c/openet/pkg/config"
    "github.com/lex0c/openet/pkg/security"
)

type Discovery struct {
    gatewayUrl  url.URL
    security    security.Security
}

type NodeInfo struct {
    Ip          net.IP
    Sign        string
    UpdatedAt   string
}

func (d *Discovery) SyncGateway() (bool, error) {
    ni, err := d.GetMyIp()

    if err != nil {
        return false, fmt.Errorf("An error occurred: %v", err)
    }

    // ...

    return true, nil
}

func (d *Discovery) GetMyIp() (NodeInfo, error) {
    resp, err := http.Get(fmt.Sprintf("%s/%s", d.gatewayUrl, config.MyIPRoute))

    defer resp.Body.Close()

    if err != nil {
        return NodeInfo{}, fmt.Errorf("An error occurred: %v", err)
    }

    body, err := io.ReadAll(resp.Body)

    if err != nil {
        return NodeInfo{}, fmt.Errorf("An error occurred: %v", err)
    }

    var ni NodeInfo

    if err = json.Unmarshal(body, &ni); err != nil {
        return NodeInfo{}, fmt.Errorf("An error occurred: %v", err)
    }

    return ni, nil
}

func (d *Discovery) FindIp(address string) (net.IP, error) {
    return nil, nil
}

