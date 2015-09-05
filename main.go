package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kellydunn/golang-geo"
	"github.com/oschwald/geoip2-golang"
)

// IPAddrResponse ...
type IPAddrResponse struct {
	IP              string  `json:"ip"`
	ASNum           uint    `json:"asnum,omitempty"`
	City            string  `json:"city"`
	Country         string  `json:"country"`
	CountryISO      string  `json:"country-iso"`
	Latitude        float64 `json:"lat"`
	Longitude       float64 `json:"lng"`
	Distance        float64 `json:"distance,omitempty"`
	Hostname        string  `json:"hostname,omitempty"`
	ASOrg           string  `json:"asorg,omitempty"`
	ISP             string  `json:"isp,omitempty"`
	Org             string  `json:"org,omitempty"`
	CIDRReport      string  `json:"cidr-report"`
	CleanTalkReport string  `json:"cleantalk-report"`
	IPInfoReport    string  `json:"ipinfo-report"`
	PeeringDBReport string  `json:"peering-db-report"`
}

func distanceFrom(lat1 float64, lng1 float64, lat2 string, lng2 string) float64 {
	ipGeoPoint := geo.NewPoint(lat1, lng1)

	parsedLat, err := strconv.ParseFloat(lat2, 64)
	if err != nil {
		log.Fatal(err)
	}

	parsedLng, err := strconv.ParseFloat(lng2, 64)
	if err != nil {
		log.Fatal(err)
	}

	lookupGeoPoint := geo.NewPoint(parsedLat, parsedLng)

	return ipGeoPoint.GreatCircleDistance(lookupGeoPoint)
}

func lookupIP(db *geoip2.Reader, citydb *geoip2.Reader, addr string, lat string, lng string) ([]byte, error) {
	var distance float64

	parsedIP := net.ParseIP(addr)

	record, err := db.ISP(parsedIP)
	if err != nil {
		return nil, err
	}

	location, err := citydb.City(parsedIP)
	if err != nil {
		return nil, err
	}

	if lat != "" && lng != "" {
		distance = distanceFrom(location.Location.Latitude, location.Location.Longitude, lat, lng)
	}

	hostnames, err := net.LookupAddr(addr)
	if err != nil {
		hostnames = []string{""}
	}

	resp := IPAddrResponse{
		Latitude:   location.Location.Latitude,
		Longitude:  location.Location.Longitude,
		Distance:   distance,
		IP:         addr,
		City:       location.City.Names["en"],
		Country:    location.Country.Names["en"],
		CountryISO: location.Country.IsoCode,
		Hostname:   hostnames[0],
		ASNum:      record.AutonomousSystemNumber,
		ASOrg:      record.AutonomousSystemOrganization,
		ISP:        record.ISP,
		Org:        record.Organization,
		CIDRReport: fmt.Sprintf(
			"http://www.cidr-report.org/cgi-bin/as-report?as=AS%d&view=2.0",
			record.AutonomousSystemNumber,
		),
		CleanTalkReport: fmt.Sprintf(
			"https://cleantalk.org/blacklists/AS%d",
			record.AutonomousSystemNumber,
		),
		IPInfoReport: fmt.Sprintf(
			"http://ipinfo.io/AS%d",
			record.AutonomousSystemNumber,
		),
		PeeringDBReport: fmt.Sprintf(
			"https://beta.peeringdb.com/api/asn/%d",
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
		port = "5000"
	}

	db, err := geoip2.Open("GeoIP2-ISP.mmdb")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	citydb, err := geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}

	defer citydb.Close()

	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/ip/:addr", func(c *gin.Context) {
		json, err := lookupIP(db, citydb, c.Param("addr"), c.Query("lat"), c.Query("lng"))
		if err != nil {
			c.Error(err)
		}

		c.Data(http.StatusOK, "application/json", json)
	})

	router.Run(":" + port)
}
