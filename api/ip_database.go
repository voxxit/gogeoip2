package api

import (
	"log"
	"os"

	"github.com/oschwald/geoip2-golang"
)

type IPDatabase struct {
	Isp  *geoip2.Reader
	City *geoip2.Reader
}

func OpenDatabases() *IPDatabase {
	db := &IPDatabase{}

	city, err := geoip2.Open(getenv("GEOIP_CITY_DBPATH", "GeoLite2-City.mmdb"))
	if err != nil {
		log.Fatal(err)
	}

	db.City = city

	isp, err := geoip2.Open(getenv("GEOIP_ISP_DBPATH", "GeoIP2-ISP.mmdb"))
	if err != nil {
		log.Println("[WARN] Could not load the GeoIP2 ISP database! For ISP data, please download it from the MaxMind website, place it in the root directory and restart.")
	} else {
		db.Isp = isp
	}

	return db
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
