package config

type Config struct {
	Listen string    `yaml:"listen"`
	APIs   []APIItem `yaml:"apis"`
}

type APIItem struct {
	Host      string   `yaml:"host,omitempty"`
	Path      string   `yaml:"path"`
	ForwardTo string   `yaml:"forward_to"`
	Protect   *Protect `yaml:"protect,omitempty"`
	Errors    map[int]ErrorRule `yaml:"errors,omitempty"`
}

type Protect struct {
	APIKey    bool        `yaml:"api_key,omitempty"`
	RateLimit string      `yaml:"rate_limit,omitempty"`
	BasicAuth *BasicAuth  `yaml:"basic_auth,omitempty"`
}

type BasicAuth struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type ErrorRule struct {
	Message string `yaml:"message,omitempty"`
}
