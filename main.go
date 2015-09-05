package main

import (
	"encoding/json"
	"net/http"
	"os"

	"gitlab.com/voxxit/gogeoip2/api"

	"github.com/gin-gonic/gin"
)

func ipLookupHandler(c *gin.Context) {
	ip, err := api.LookupIP(&api.IPConfig{
		Addr:             c.Param("addr"),
		RequestLatitude:  c.Query("lat"),
		RequestLongitude: c.Query("lng"),
	})

	if err != nil {
		c.Error(err)
	}

	data, err := json.Marshal(&ip)
	if err != nil {
		c.Error(err)
	}

	c.Data(http.StatusOK, "application/json", data)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/ip/:addr", ipLookupHandler)

	router.Run(":" + port)
}
