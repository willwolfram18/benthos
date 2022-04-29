package tracer

type OtlpConfig struct {
	Endpoint string            `json:"endpoint" yaml:"endpoint"`
	Protocol string            `json:"protocol" yaml:"protocol"`
	Tags     map[string]string `json:"tags" yaml:"tags"`
	Timeout  string            `json:"timeout" yaml:"timeout"`
}

func NewOtlpConfig() OtlpConfig {
	return OtlpConfig{
		Endpoint: "",
		Protocol: "",
		Tags:     map[string]string{},
		Timeout:  "",
	}
}
