package http_gateway

import (
    "github.com/gin-gonic/gin"
    "log"
    "net/http"
)

type serverUtil struct {
}

var ServerUtil *serverUtil

func (s *serverUtil) Init(mode string) *gin.Engine {
    gin.SetMode(mode)
    return gin.Default()
}

func (s *serverUtil) Start(mode string, https bool, address, certFile, keyFile string) {
    gin.SetMode(mode)
    engine := gin.Default()
    var err error
    if https {
        err = http.ListenAndServeTLS(address, certFile, keyFile, engine)
    } else {
        err = http.ListenAndServe(address, engine)
    }
    if err != nil {
        log.Fatal(err)
    }
    return
}

func (s *serverUtil) StartByConfig(mode string, httpConfigFormat HttpConfigFormat) {
    var err error
    gin.SetMode(mode)
    engine := gin.Default()
    if httpConfigFormat.Https {
        err = http.ListenAndServeTLS(httpConfigFormat.Address, httpConfigFormat.CertFile, httpConfigFormat.KeyFile, engine)
    } else {
        err = http.ListenAndServe(httpConfigFormat.Address, engine)
    }
    if err != nil {
        log.Fatal(err)
    }
    return
}
