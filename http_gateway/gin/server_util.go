package gin_http_gateway

import (
    "github.com/gin-gonic/gin"
    "github.com/maxtech/gateway/http_gateway"
    "log"
    "net/http"
)

type serverUtil struct {
}

var ServerUtil *serverUtil

func (s *serverUtil) InitGin(mode string) *gin.Engine {
    gin.SetMode(mode)
    return gin.Default()
}

func (s *serverUtil) StartGin(engine *gin.Engine, https bool, address, certFile, keyFile string) {
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

func (s *serverUtil) StartGinByConfig(engine *gin.Engine, httpConfigFormat http_gateway.HttpConfigFormat) {
    var err error
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
