package grpc_gateway

type GRPCConfigFormat struct {
    Address string `yaml:"address"`
    Https   bool   `yaml:"https"`
}
