package api

import (
	"log"
	"strconv"

	"github.com/kellydunn/golang-geo"
)

var (
	Countries = map[string]bool{
		"AE": true, "AR": true, "AU": true, "BG": true, "BM": true, "BR": true, "CA": true, "CH": true, "CL": true,
		"CO": true, "CR": true, "CZ": true, "DK": true, "DO": true, "GB": true, "HU": true, "IS": true, "MA": true,
		"MX": true, "NO": true, "NZ": true, "PA": true, "PE": true, "PL": true, "RU": true, "SA": true, "SE": true,
		"TR": true, "US": true, "UY": true, "ZW": false, "UA": false, "YE": false, "TH": false, "CN": false,
		"RO": false, "VN": false,
	}
)

func distanceTo(lat1 float64, lng1 float64, lat2 string, lng2 string) float64 {
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
