package grpc_gateway

import "fmt"

var Version string

func init() {
    Version = fmt.Sprintf(
        "|- %s module:\t\t\t%s",
        "grpc gateway",
        "0.0.1")
}
