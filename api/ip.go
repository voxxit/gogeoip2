package api

import (
	"fmt"
	"net"
)

var db *IPDatabase

type IPConfig struct {
	Addr             string
	RequestLatitude  string
	RequestLongitude string
}

type IP struct {
	IP              string   `json:"ip"`
	Scores          *Scores  `json:"scores"`
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
	CIDRReport      string   `json:"cidr-report"`
	CleanTalkReport string   `json:"cleantalk-report"`
	IPInfoReport    string   `json:"ipinfo-report"`
	PeeringDBReport string   `json:"peering-db-report"`
}

func LookupIP(config *IPConfig) (*IP, error) {
	db = OpenDatabases()

	defer db.isp.Close()
	defer db.city.Close()

	parsedIP := net.ParseIP(config.Addr)

	record, err := db.isp.ISP(parsedIP)
	if err != nil {
		return nil, err
	}

	location, err := db.city.City(parsedIP)
	if err != nil {
		return nil, err
	}

	hostnames, _ := net.LookupAddr(config.Addr)

	distance := distanceTo(
		location.Location.Latitude,
		location.Location.Longitude,
		config.RequestLatitude,
		config.RequestLongitude,
	)

	scores := NewScores(&ScoresInput{
		Distance:    distance,
		CountryCode: location.Country.IsoCode,
	})

	ip := &IP{
		Scores:     scores,
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

	return ip, nil
}
