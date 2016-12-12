package api

import (
	"fmt"
	"log"
	"net"
)

type IPConfig struct {
	Addr             string
	RequestLatitude  string
	RequestLongitude string
}

type IP struct {
	IP              string   `json:"ip"`
	ASNum           uint     `json:"asnum,omitempty"`
	City            string   `json:"city"`
	Country         string   `json:"country"`
	CountryISO      string   `json:"country-iso"`
	Continent       string   `json:"continent"`
	Latitude        float64  `json:"lat"`
	Longitude       float64  `json:"lng"`
	TimeZone        string   `json:"time-zone"`
	Distance        float64  `json:"distance,omitempty"`
	Hostnames       []string `json:"hostnames,omitempty"`
	ASOrg           string   `json:"asorg,omitempty"`
	ISP             string   `json:"isp,omitempty"`
	Org             string   `json:"org,omitempty"`
	CIDRReport      string   `json:"cidr-report,omitempty"`
	CleanTalkReport string   `json:"cleantalk-report,omitempty"`
	IPInfoReport    string   `json:"ipinfo-report,omitempty"`
	PeeringDBReport string   `json:"peering-db-report,omitempty"`
}

func LookupIP(db *IPDatabase, config *IPConfig) (*IP, error) {
	parsedIP := net.ParseIP(config.Addr)

	location, err := db.City.City(parsedIP)
	if err != nil {
		return nil, err
	}

	hostnames, _ := net.LookupAddr(config.Addr)

	distance, err := distanceTo(
		location.Location.Latitude,
		location.Location.Longitude,
		config.RequestLatitude,
		config.RequestLongitude,
	)
	if err != nil {
		log.Println("Distance could not be calculated")
	}

	ip := &IP{
		Latitude:   location.Location.Latitude,
		Longitude:  location.Location.Longitude,
		IP:         config.Addr,
		Distance:   distance,
		TimeZone:   location.Location.TimeZone,
		City:       location.City.Names["en"],
		Country:    location.Country.Names["en"],
		CountryISO: location.Country.IsoCode,
		Continent:  location.Continent.Code,
		Hostnames:  hostnames,
	}

	if db.Isp != nil {
		record, err := db.Isp.ISP(parsedIP)
		if err != nil {
			return nil, err
		}

		ip.ASNum = record.AutonomousSystemNumber
		ip.ASOrg = record.AutonomousSystemOrganization
		ip.ISP = record.ISP
		ip.Org = record.Organization

		ip.CIDRReport = fmt.Sprintf(
			"http://www.cidr-report.org/cgi-bin/as-report?as=AS%d&view=2.0",
			record.AutonomousSystemNumber,
		)

		ip.CleanTalkReport = fmt.Sprintf(
			"https://cleantalk.org/blacklists/AS%d",
			record.AutonomousSystemNumber,
		)

		ip.IPInfoReport = fmt.Sprintf(
			"http://ipinfo.io/AS%d",
			record.AutonomousSystemNumber,
		)

		ip.PeeringDBReport = fmt.Sprintf(
			"https://beta.peeringdb.com/api/net?asn=%d",
			record.AutonomousSystemNumber,
		)
	}

	return ip, nil
}
