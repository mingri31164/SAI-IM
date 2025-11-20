package configserver

import (
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/conf"
)

var ErrNotSetConfig = errors.New("none config info...")

type ConfigServer interface {
	FromJsonBytes() ([]byte, error) // 数据来源(json)
}

// 代理类
type configSever struct {
	ConfigServer
	configFile string
}

func NewConfigServer(configFile string, s ConfigServer) *configSever {
	return &configSever{
		ConfigServer: s,
		configFile:   configFile,
	}
}

// 解析配置加载
func (s *configSever) MustLoad(v any) error {
	if s.configFile == "" && s.ConfigServer == nil {
		return ErrNotSetConfig
	}
	if s.ConfigServer == nil {
		// 使用go-zero默认方式
		conf.MustLoad(s.configFile, v)
		return nil
	}

	// 使用配置中心的加载方式
	data, err := s.ConfigServer.FromJsonBytes()
	if err != nil {
		return err
	}

	return LoadFromJsonBytes(data, v)
}

func LoadFromJsonBytes(data []byte, v any) error {
	return conf.LoadFromJsonBytes(data, v)
}
