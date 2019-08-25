package http_gateway

type HttpConfigFormat struct {
    Address  string `json:"address" yaml:"address"`
    Https    bool   `json:"https" yaml:"https"`
    CertFile string `json:"cert_file" yaml:"cert_file"`
    KeyFile  string `json:"key_file" yaml:"key_file"`
}
