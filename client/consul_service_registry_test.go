package client

//import (
//	"fmt"
//	"net"
//	"testing"
//
//	"github.com/gin-gonic/gin"
//)
//
//// 获取本地IP
//func getLocalIP() string {
//	netInterfaces, err := net.Interfaces()
//	if err != nil {
//		return ""
//	}
//	for i := 0; i < len(netInterfaces); i++ {
//		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
//			addrs, _ := netInterfaces[i].Addrs()
//			for _, address := range addrs {
//				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
//					if ipnet.IP.To4() != nil {
//						return ipnet.IP.String()
//					}
//				}
//			}
//		}
//	}
//	return ""
//}
//
//func TestConsulServiceRegistry(t *testing.T) {
//	registryClient, _ = NewConsulServiceRegistryAddress("192.168.5.110:8500", "")
//	hostAddress := getLocalIP()
//	serviceInstance := DefaultServiceInstance{
//		InstanceId:     "123456789",
//		ServiceName:    "test123",
//		Host:           hostAddress,
//		Port:           8026,
//		Metadata:       []string{"new"},
//		Timeout:        5,
//		Interval:       15,
//		DeregisterTime: 20,
//		CheckHTTP:      fmt.Sprintf("http://%s:%d/%s/%s", hostAddress, 8026, "consulhealth", "123456789"),
//	}
//	serviceInstanceInfo, _ := NewDefaultServiceInstance(&serviceInstance)
//	registryClient.Register(serviceInstanceInfo)
//	r := gin.Default()
//	// 健康检测接口，其实只要是 200 就认为成功了
//	r.GET("/consulhealth/*any", func(c *gin.Context) {
//		c.JSON(200, gin.H{
//			"message": "ok",
//		})
//	})
//	err := r.Run(":8026")
//	if err != nil {
//		registryClient.Deregister()
//	}
//}
