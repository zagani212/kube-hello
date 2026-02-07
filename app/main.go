package main

import (
    "net/http"
	"os"
	"context"
	"os/signal"
	"syscall"
	"fmt"
	"time"
	"net"
	"github.com/semihalev/gin-stats"
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
	i.Version = "1.0.0"
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
	i.Version = "1.0.0"
	i.Environment = os.Getenv("APP_ENV")
	i.StartedAt = now.Format("01-02-2006 15:04:05")
	c.IndentedJSON(http.StatusOK, i)
}
var now time.Time
func main() {
	now = time.Now()
	router := gin.Default()
	router.Use(stats.RequestStats())
	router.GET("/", getInstance)
	router.GET("/health", getHealth)
	router.GET("/info", getInfo)
	router.GET("/metrics", func(c *gin.Context) {
		c.JSON(http.StatusOK, stats.Report())
	})

	router.Run(":8080")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router.Handler(),
	}
	signalChan := make(chan os.Signal, 1)
 	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)


	 // Wait for SIGTERM or SIGINT
	<-signalChan
	fmt.Println("\nSIGTERM received. Waiting 2 seconds to drain requests...")

	// Wait 2 seconds to allow ongoing requests to complete
	time.Sleep(2 * time.Second)

	// Gracefully shut down the server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Error during shutdown: %s\n", err)
	}

	fmt.Println("Server stopped gracefully.")

}