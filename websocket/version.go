package websocket

import "fmt"

var Version string

func init() {
    Version = fmt.Sprintf(
        "|- %s module:\t\t\t%s",
        "websocket",
        "0.0.1")
}
