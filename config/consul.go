package config

import (
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var cfgConsul *viper.Viper

var ConsulConfig =new (Consul)

type Consul struct {
	Addr string;          // 地址
	ID  string;           // id 唯一标识
	Name string;         // 服务名称
	Tags string;         // 标签
	CheckPort int;       // 检查端口号
	RegistPort int;      // 注册端口号
	Timeout int          // 超时时间
	IntervalTime int;    // 健康检查间隔
	DeregisterTime int  // 服务注销时间，相当于过期时间
	HealthUrl string // 健康检查url
}

//载入配置文件
func SetupConsulConfig(path string) {
	viper.SetConfigFile(path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(fmt.Sprintf("Read config file fail: %s", err.Error()))
	}
	err = viper.ReadConfig(strings.NewReader(os.ExpandEnv(string(content))))
	if err != nil {
		log.Fatal(fmt.Sprintf("Parse config file fail: %s", err.Error()))
	}
	cfgConsul=viper.Sub("config.consul")
	if cfgConsul==nil {
		panic("No found settings.consul in the configuration")
	}
	ConsulConfig=initConsul(cfgConsul)
}


func initConsul(cfg *viper.Viper) *Consul {
	config:=&Consul{
		Addr:	cfg.GetString("addr"),
		ID :	cfg.GetString("id"),
		Name :	cfg.GetString("name"),
		Tags :	cfg.GetString("tags"),
		CheckPort :	cfg.GetInt("check_port"),
		RegistPort :	cfg.GetInt("regist_port"),
		Timeout :	cfg.GetInt("time_out"),
		IntervalTime :	cfg.GetInt("interval_time"),
		DeregisterTime :	cfg.GetInt("deregister_time"),
		HealthUrl: cfg.GetString("health_url"),
	}
	return  config
}



