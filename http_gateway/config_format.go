package http_gateway

type HttpConfigFormat struct {
    HttpAddress  string `json:"http_address" yaml:"http_address"`
    HttpsAddress string `json:"https_address" yaml:"https_address"`
    Https        bool   `json:"https" yaml:"https"`
    CertFile     string `json:"cert_file" yaml:"cert_file"`
    KeyFile      string `json:"key_file" yaml:"key_file"`
}
