package api

import (
	"log"

	"github.com/oschwald/geoip2-golang"
)

type IPDatabase struct {
	isp  *geoip2.Reader
	city *geoip2.Reader
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

	db.isp = isp
	db.city = city

	return db
}
