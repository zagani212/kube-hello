package main

import (
    "net/http"
	"os"
	"context"
	"os/signal"
	"syscall"
	"fmt"
	"strconv"
	"time"
	"net"
    "github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metrics struct {
	http_requests_total *prometheus.CounterVec
	http_request_duration_seconds_bucket prometheus.Histogram
}

type instance struct {
    IP     string  `json:"ip"`
    Message  string  `json:"message"`
    Hostname string  `json:"hostname"`
    Version  string `json:"version"`
}

type response struct {
	Status	string `json:"status"` 	
}

type workStruct struct {
	Worked	bool `json:"worked"` 	
	DurationMs	int `json:"surationMs"` 	
	Hostname	string `json:"hostname"` 	
}

type info struct {
    App     string  `json:"app"`
    Version  string  `json:"version"`
    Environment string  `json:"environment"`
    StartedAt  string `json:"startedAt"`
}

func newMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		http_requests_total: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "No of request handled the server",
			},
			[]string{"method", "path"},
		),
		http_request_duration_seconds_bucket: promauto.With(reg).NewHistogram(
			prometheus.HistogramOpts{
				Name: "http_request_duration_seconds_bucket",
				Help: "The duration of requests in seconds",
				Buckets: prometheus.LinearBuckets(0.1, 0.1, 10),
			}),
	}
	return m
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

func work(c *gin.Context){
	var res workStruct
	res.Worked = true
	res.DurationMs, _ = strconv.Atoi(c.DefaultQuery("duration", "500"))
	res.Hostname, _ = os.Hostname()
	time.Sleep(time.Duration(res.DurationMs) * time.Millisecond)
	c.IndentedJSON(http.StatusOK, res)
}

func durationMiddleware(m *metrics) gin.HandlerFunc {
	return func (c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start).Seconds()
		path := c.FullPath()
		method := c.Request.Method
		m.http_request_duration_seconds_bucket.Observe(duration)
		m.http_requests_total.WithLabelValues(method, path).Inc()
	}
}

var now time.Time
func main() {
	reg := prometheus.NewRegistry()
	m := newMetrics(reg)
	now = time.Now()
	router := gin.Default()
	router.Use(durationMiddleware(m))
	router.GET("/", getInstance)
	router.GET("/health", getHealth)
	router.GET("/info", getInfo)
	router.GET("/work", work)
	router.GET("/metrics", gin.WrapH(promhttp.HandlerFor(reg, promhttp.HandlerOpts{})))

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