package grpc_gateway

import (
    "github.com/maxtech/log"
    "google.golang.org/grpc"
    "net"
    "sync"
    "time"
)

type GRPCServer struct {
    Logger log.AppLogger
    mu     sync.Mutex
    lis    net.Listener
    server *grpc.Server
    status chan int
}

func Init(_initLogger log.AppLogger, _address string) *GRPCServer {
    var err error
    server := new(GRPCServer)
    server.status = make(chan int, 1)
    server.Logger = _initLogger
    server.lis, err = net.Listen("tcp", _address)
    if err != nil {
        _initLogger.Error("failed to listen:", err.Error())
    }
    server.server = grpc.NewServer()

    return server
}

func (gs *GRPCServer) GetGRPCServer() *grpc.Server {
    return gs.server
}

func (gs *GRPCServer) Start() {
    for {
        if len(gs.status) < 1 {
            gs.mu.Lock()
            gs.status <- 1
            gs.mu.Unlock()
            go gs.run(gs.server, gs.lis)
        }
        time.Sleep(time.Second)
    }
}

func (gs *GRPCServer) run(_server *grpc.Server, _lis net.Listener) {
    defer func() {
        if err := recover(); err != nil {
            gs.Logger.Error("recover error:", err)
            gs.mu.Lock()
            <-gs.status
            gs.mu.Unlock()
        }
    }()

    if err := _server.Serve(_lis); err != nil {
        gs.Logger.Error("failed to serve:", err.Error())
    }
}
