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

func (s *serverUtil) Start(engine *gin.Engine, https bool, address, certFile, keyFile string) {
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
