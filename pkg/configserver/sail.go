package configserver

import (
	"encoding/json"
	"github.com/HYY-yu/sail-client"
)

type Config struct {
	ETCDEndpoints  string `toml:"etcd_endpoints"`
	ProjectKey     string `toml:"project_key"`
	Namespace      string `toml:"namespace"`
	Configs        string `toml:"configs"`
	ConfigFilePath string `toml:"config_file_path"`
	LogLevel       string `toml:"log_level"`
}

type Sail struct {
	*sail.Sail
}

func NewSail(cfg *Config) *Sail {
	s := sail.New(&sail.MetaConfig{
		ETCDEndpoints:  cfg.ETCDEndpoints,
		ProjectKey:     cfg.ProjectKey,
		Namespace:      cfg.Namespace,
		Configs:        cfg.Configs,
		ConfigFilePath: cfg.ConfigFilePath,
		LogLevel:       cfg.LogLevel,
	})
	return &Sail{Sail: s}
}

func (s *Sail) FromJsonBytes() ([]byte, error) {
	if err := s.Pull(); err != nil {
		return nil, err
	}
	return s.fromJsonBytes(s.Sail)
}

func (s *Sail) fromJsonBytes(sail *sail.Sail) ([]byte, error) {
	v, err := sail.MergeVipers()
	if err != nil {
		return nil, err
	}
	data := v.AllSettings()
	return json.Marshal(data)
}
