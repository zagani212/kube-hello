package main

import (
    "net/http"
	"os"
	"fmt"
	"net"
    "github.com/gin-gonic/gin"
)

type instance struct {
    IP     string  `json:"ip"`
    Message  string  `json:"message"`
    Hostname string  `json:"hostname"`
    Version  string `json:"version"`
}

type info struct {
    App     string  `json:"app"`
    Version  string  `json:"version"`
    Environment string  `json:"environment"`
    StartedAt  float64 `json:"startedAt"`
}

func GetLocalIP() net.IP {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        
    }
    defer conn.Close()

    localAddress := conn.LocalAddr().(*net.UDPAddr)

    return localAddress.IP
}

func getInstance(c *gin.Context){
	var i instance
	
	i.IP = fmt.Sprintf("%s",GetLocalIP())
	i.Message = os.Getenv("MESSAGE")
	i.Hostname, _ = os.Hostname()
	i.Version = os.Getenv("APP_VERSION")
	c.IndentedJSON(http.StatusOK, i)
}


func main() {
	router := gin.Default()
	router.GET("/", getInstance)

	router.Run("localhost:8080")
}