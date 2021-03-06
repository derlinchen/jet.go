package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

type Setting struct {
	vp *viper.Viper
}

func NewSetting(configs []string) (*Setting, error) {
	vp := viper.New()
	// 设置配置文件名，配置文件路径，配置文件格式
	vp.SetConfigName("config")
	for _, config := range configs {
		if config != "" {
			vp.AddConfigPath(config)
		}
	}
	vp.AddConfigPath("./")
	vp.AddConfigPath("config/")
	vp.SetConfigType("yaml")
	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}
	s := &Setting{vp}
	s.WatchSettingChange()
	return s, nil
}

// 监听文件变化，配置热更新
func (s *Setting) WatchSettingChange() {
	go func() {
		s.vp.WatchConfig()
		s.vp.OnConfigChange(func(in fsnotify.Event) {
			log.Println("配置文件热更新", in.Name)
			_ = s.ReloadAllSection()
		})
	}()
}
