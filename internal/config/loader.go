package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func Load(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, fmt.Errorf("parse yaml: %w", err)
	}

	if err := validate(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func validate(cfg *Config) error {
	cfg.Listen = strings.TrimSpace(cfg.Listen)
	if cfg.Listen == "" {
		return errors.New("config error: 'listen' wajib diisi, contoh ':8080'")
	}
	if len(cfg.APIs) == 0 {
		return errors.New("config error: 'apis' minimal 1 item")
	}

	for i, api := range cfg.APIs {
		if strings.TrimSpace(api.Path) == "" {
			return fmt.Errorf("config error: apis[%d].path wajib diisi", i)
		}
		if !strings.HasPrefix(api.Path, "/") {
			return fmt.Errorf("config error: apis[%d].path harus diawali '/'", i)
		}
		if strings.TrimSpace(api.ForwardTo) == "" {
			return fmt.Errorf("config error: apis[%d].forward_to wajib diisi", i)
		}
		// Validasi ringan: forward_to harus http(s)
		if !strings.HasPrefix(api.ForwardTo, "http://") && !strings.HasPrefix(api.ForwardTo, "https://") {
			return fmt.Errorf("config error: apis[%d].forward_to harus diawali http:// atau https://", i)
		}
	}

	return nil
}
