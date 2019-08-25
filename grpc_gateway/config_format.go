package grpc_gateway

type GRPCConfigFormat struct {
    Address string `json:"address" yaml:"address"`
    Https   bool   `json:"https" yaml:"https"`
}
