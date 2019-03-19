package websocket

import (
    "github.com/gin-gonic/gin"
    "github.com/maxtech/log"
    "net/http"
)

type ws struct {
}

var WS *ws

func (*ws) Handler(ctx *gin.Context) {
    var params wsParams
    err := ctx.ShouldBindQuery(&params)

    if err != nil {
        ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]interface{}{
            "code": http.StatusBadRequest,
            "msg":  "参数解析错误",
        })
        return
    }

    userIdInterface, _ := ctx.Get("user_id")
    if userIdInterface == nil {
        ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]interface{}{
            "code": http.StatusBadRequest,
            "msg":  "登录信息错误",
        })
        return
    }

    params.UserId = userIdInterface.(uint64)

    serveWs(Hub, ctx.Writer, ctx.Request, params)
}

type Message struct {
    Sender   string `json:"sender"`
    Receiver string `json:"receiver"`
    IsDirect bool   `json:"is_direct"`
    Topic    string `json:"topic"`
    IsHost   bool   `json:"is_host"`
    Msg      string `json:"msg"`
}

func InitHub() {
    logger = log.NewLogger("websocket")
    Hub = newHub()
    go Hub.run()
}
