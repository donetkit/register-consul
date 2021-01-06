package consul

import (
	"fmt"
	"github.com/donetkit/register-consul/config"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var NodeId string = getUUID()

// 注册consulClient
// serviceName 服务名称  serverNode
// tag 服务标签  v1
func RegisterConsulClient()  {
	go registerConsulServer()
	time.Sleep(time.Second*1)
	var lastIndex uint64
	configConsul := api.DefaultConfig()
	configConsul.Address =config.ConsulConfig.Addr // "10.0.0.10:8500" //consul server
	clientConsul, err := api.NewClient(configConsul)
	if err != nil {
		fmt.Println("api new client is failed, err:", err)
		return
	}
	services, metainfo, err := clientConsul.Health().Service(config.ConsulConfig.Name, config.ConsulConfig.Tags, true, &api.QueryOptions{
		WaitIndex: lastIndex, // 同步点，这个调用将一直阻塞，直到有新的更新
	})
	if err != nil {
		logrus.Warn("error retrieving instances from Consul: %v", err)
	}
	lastIndex = metainfo.LastIndex
	addrs := map[string]struct{}{}
	for _, service := range services {
		fmt.Println("service.Service.Address:", service.Service.Address, "service.Service.Port:", service.Service.Port)
		addrs[net.JoinHostPort(service.Service.Address, strconv.Itoa(service.Service.Port))] = struct{}{}
	}
}


func registerConsulServer() {
	configConsul := api.DefaultConfig()
	configConsul.Address = config.ConsulConfig.Addr
	client, err := api.NewClient(configConsul)
	if err != nil {
		log.Fatal("consul client error : ", err)
	}
	registration := new(api.AgentServiceRegistration)
	if config.ConsulConfig.ID == ""{
		config.ConsulConfig.ID = NodeId
	}
	registration.ID = config.ConsulConfig.ID                 // 服务节点的名称
	registration.Name = config.ConsulConfig.Name             // 服务名称
	registration.Port = config.ConsulConfig.RegistPort       //9527 // 服务端口
	registration.Tags = []string{config.ConsulConfig.Tags}   //tag，可以为空
	registration.Address = getLocalIP()  //服务 IP
	url := fmt.Sprintf("http://%s:%d/%s/%s", registration.Address, config.ConsulConfig.CheckPort, "consulhealth",NodeId)
	if config.ConsulConfig.HealthUrl != "" {
		url = config.ConsulConfig.HealthUrl
	}
	registration.Check = &api.AgentServiceCheck{ // 健康检查
		HTTP:                      url,
		Timeout:                   fmt.Sprintf("%d%s",config.ConsulConfig.Timeout,"s") ,
		Interval:                  fmt.Sprintf("%d%s",config.ConsulConfig.IntervalTime,"s") ,  // 健康检查间隔
		DeregisterCriticalServiceAfter: fmt.Sprintf("%d%s",config.ConsulConfig.DeregisterTime,"s") , //check失败后30秒删除本服务，注销时间，相当于过期时间
		// GRPC: fmt.Sprintf("%v:%v/%v", IP, r.Port, r.Service),// grpc 支持，执行健康检查的地址，service 会传到 Health.Check 函数中
	}

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		log.Fatal("register server error : ", err)
	}

	if config.ConsulConfig.HealthUrl == "" {
		http.HandleFunc(fmt.Sprintf("/%s/%s","consulhealth",NodeId), consulCheck)
		http.ListenAndServe(fmt.Sprintf(":%d", config.ConsulConfig.CheckPort), nil)
	}

}

// consul 服务端会自己发送请求，来进行健康检查
func consulCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Println("consul Check remote:" + r.RemoteAddr + r.URL.String())
}

// 获取本地IP
func getLocalIP() string {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()
			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return  ipnet.IP.String()
					}
				}
			}
		}
	}
	return ""
}

// uuid
func getUUID() string {
	u, _ := uuid.NewRandom()
	return strings.ReplaceAll(u.String(),"-","")
}



