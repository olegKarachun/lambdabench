package config

type Config struct {
	Region      string `yaml:"region"`
	Function    string `yaml:"function"`
	Payload     string `yaml:"payload"`
	Requests    int    `yaml:"requests"`
	Concurrency int    `yaml:"concurrency"`
	Warmup      bool   `yaml:"warmup"`
}
