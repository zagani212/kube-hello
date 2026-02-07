package main

import (
    "net/http"
	"os"
	"fmt"
	"time"
	"net"
    "github.com/gin-gonic/gin"
)

type instance struct {
    IP     string  `json:"ip"`
    Message  string  `json:"message"`
    Hostname string  `json:"hostname"`
    Version  string `json:"version"`
}

type response struct {
	Status	string `json:"status"` 	
}

type info struct {
    App     string  `json:"app"`
    Version  string  `json:"version"`
    Environment string  `json:"environment"`
    StartedAt  string `json:"startedAt"`
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

func getHealth(c *gin.Context){
	var status response
	status.Status = "ok"
	c.IndentedJSON(http.StatusOK, status)
}

func getInfo(c *gin.Context){
	var i info
	
	i.App = os.Getenv("APP_NAME")
	i.Version = os.Getenv("APP_VERSION")
	i.Environment = os.Getenv("APP_ENV")
	i.StartedAt = now.Format("01-02-2006 15:04:05")
	c.IndentedJSON(http.StatusOK, i)
}
var now time.Time
func main() {
	now = time.Now()
	router := gin.Default()
	router.GET("/", getInstance)
	router.GET("/health", getHealth)
	router.GET("/info", getInfo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router.Run("localhost:"+port)
}