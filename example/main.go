package main

import (
	"register-consul/config"
	"register-consul/consul"
	"time"
)

func main() {
	config.SetupConsulConfig("./example/consul.yml")
	consul.RegisterConsulClient()
	time.Sleep(time.Second*120)
}
