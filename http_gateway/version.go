package http_gateway

import "fmt"

var Version string

func init() {
    Version = fmt.Sprintf(
        "|- %s module:\t\t\t%s",
        "http gateway",
        "0.0.1")
}
