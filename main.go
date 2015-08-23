package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/oschwald/geoip2-golang"
)

// IPAddrResponse ...
type IPAddrResponse struct {
	ASNum      uint   `json:"asnum,omitempty"`
	ASOrg      string `json:"asorg,omitempty"`
	ISP        string `json:"isp,omitempty"`
	Org        string `json:"org,omitempty"`
	CIDRReport string `json:"cidr-report"`
}

func lookupIP(db *geoip2.Reader, addr string) ([]byte, error) {
	parsedIP := net.ParseIP(addr)

	record, err := db.ISP(parsedIP)
	if err != nil {
		return nil, err
	}

	resp := IPAddrResponse{
		ASNum: record.AutonomousSystemNumber,
		ASOrg: record.AutonomousSystemOrganization,
		ISP:   record.ISP,
		Org:   record.Organization,
		CIDRReport: fmt.Sprintf(
			"http://www.cidr-report.org/cgi-bin/as-report?as=AS%d&view=2.0",
			record.AutonomousSystemNumber,
		),
	}

	b, err := json.Marshal(&resp)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	db, err := geoip2.Open("GeoIP2-ISP.mmdb")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/ip/:addr", func(c *gin.Context) {
		json, err := lookupIP(db, c.Param("addr"))
		if err != nil {
			c.Error(err)
		}

		c.Data(http.StatusOK, "application/json", json)
	})

	router.Run(":" + port)
}
