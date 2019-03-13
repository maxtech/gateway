package grpc_gateway

type GrpcConfigFormat struct {
    Address  string `yaml:"address"`
    Https    bool   `yaml:"https"`
}
