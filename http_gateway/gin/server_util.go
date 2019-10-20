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

func (su *serverUtil) InitGin(_mode string) *gin.Engine {
    gin.SetMode(_mode)
    return gin.Default()
}

func (su *serverUtil) StartGin(_engine *gin.Engine, _https bool, _address, _certFile, _keyFile string) {
    var err error
    if _https {
        err = http.ListenAndServeTLS(_address, _certFile, _keyFile, _engine)
    } else {
        err = http.ListenAndServe(_address, _engine)
    }
    if err != nil {
        log.Fatal(err)
    }
    return
}

func (su *serverUtil) StartGinByConfig(_engine *gin.Engine, _httpConfigFormat http_gateway.HttpConfigFormat) {
    var err error
    if _httpConfigFormat.Https {
        err = http.ListenAndServeTLS(_httpConfigFormat.HttpAddress, _httpConfigFormat.CertFile, _httpConfigFormat.KeyFile, _engine)
    } else {
        err = http.ListenAndServe(_httpConfigFormat.HttpAddress, _engine)
    }
    if err != nil {
        log.Fatal(err)
    }
    return
}
