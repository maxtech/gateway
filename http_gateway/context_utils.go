package http_gateway

import (
    "github.com/gin-gonic/gin"
    "net"
    "strings"
)

func GetHostAndPortFromContext(context *gin.Context) (host, port string, err error) {
    ipStrings := make([]string, 0)
    ips := context.Request.Header.Get("X-Forwarded-For")
    if ips != "" {
        ipStrings = strings.Split(ips, ",")
    } else {
        ipStrings = append(ipStrings, context.Request.RemoteAddr)
    }

    if len(ipStrings) > 0 && ipStrings[0] != "" {
        host, port, err = net.SplitHostPort(ipStrings[0])
    }

    return host, port, err
}
