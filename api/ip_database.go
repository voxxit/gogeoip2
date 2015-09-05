package api

import (
	"log"

	"github.com/oschwald/geoip2-golang"
)

type IPDatabase struct {
	Isp  *geoip2.Reader
	City *geoip2.Reader
}

func OpenDatabases() *IPDatabase {
	db := &IPDatabase{}

	isp, err := geoip2.Open("GeoIP2-ISP.mmdb")
	if err != nil {
		log.Fatal(err)
	}

	city, err := geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}

	db.Isp = isp
	db.City = city

	return db
}
