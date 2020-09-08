package yml_config

import (
	"goskeleton/app/global/my_errors"
	"goskeleton/app/global/variable"
	"log"

	"github.com/spf13/viper"
)

func CreateYamlFactory() *ymlConfig {
	yamlConfig := viper.New()
	yamlConfig.AddConfigPath(variable.BasePath + "/config")
	// 需要读取的配置文件名
	yamlConfig.SetConfigName("config")
	// 设置配置文件类型
	yamlConfig.SetConfigType("yml")

	if err := yamlConfig.ReadInConfig(); err != nil {
		log.Fatal(my_errors.ErrorsConfigInitFail + err.Error())
	}
	return &ymlConfig{
		yamlConfig,
	}
}

type ymlConfig struct {
	viper *viper.Viper
}

// Get 一个原始值
func (c *ymlConfig)Get(keyName string) interface{} {
	return c.viper.Get(keyName)
}

// GetString
func (c *ymlConfig)GetString(keyName string) string {
	return c.viper.GetString(keyName)
}

// GetInt
func (c *ymlConfig)GetInt(keyName string) int {
	return c.viper.GetInt(keyName)
}