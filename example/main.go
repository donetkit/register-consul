package main

import (
	"github.com/donetkit/register-consul/config"
	"github.com/donetkit/register-consul/consul"
	"time"
)

func main() {
	config.SetupConsulConfig("./example/consul.yml")
	consul.RegisterConsulClient()
	time.Sleep(time.Second*120)
}
