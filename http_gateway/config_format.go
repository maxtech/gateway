package http_gateway

type HttpConfigFormat struct {
    Address  string `yaml:"address"`
    Https    bool   `yaml:"https"`
    CertFile string `yaml:"cert_file"`
    KeyFile  string `yaml:"key_file"`
}
