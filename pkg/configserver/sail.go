package configserver

import (
	"encoding/json"
	"fmt"
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
	//🔥注意：此处定义的不是configserver中的OnChange
	//        而是sail client配置变更时真正要执行的方法
	sail.OnConfigChange
	c *Config
}

func NewSail(cfg *Config) *Sail {
	//s := sail.New(&sail.MetaConfig{
	//	ETCDEndpoints:  cfg.ETCDEndpoints,
	//	ProjectKey:     cfg.ProjectKey,
	//	Namespace:      cfg.Namespace,
	//	Configs:        cfg.Configs,
	//	ConfigFilePath: cfg.ConfigFilePath,
	//	LogLevel:       cfg.LogLevel,
	//})
	return &Sail{c: cfg}
}

func (s *Sail) Build() error {
	var opts []sail.Option
	if s.OnConfigChange != nil {
		opts = append(opts, sail.WithOnConfigChange(s.OnConfigChange))
	}
	s.Sail = sail.New(&sail.MetaConfig{
		ETCDEndpoints:  s.c.ETCDEndpoints,
		ProjectKey:     s.c.ProjectKey,
		Namespace:      s.c.Namespace,
		Configs:        s.c.Configs,
		ConfigFilePath: s.c.ConfigFilePath,
		LogLevel:       s.c.LogLevel,
	}, opts...)
	return s.Err()
}

// SetOnChange 配置更新时需要执行的方法
func (s *Sail) SetOnChange(f OnChange) {
	s.OnConfigChange = func(configFileKey string, sail *sail.Sail) {
		// 重新获取配置信息
		data, err := s.fromJsonBytes(sail)
		if err != nil {
			fmt.Println(err)
			return
		}
		// 调用传入的 f 方法进行加载
		if err = f(data); err != nil {
			fmt.Println("OnChange err:", err)
		}
	}
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
