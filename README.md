# gogeoip2

[![Gitter](https://badges.gitter.im/gogeoip2/Lobby.svg)](https://gitter.im/gogeoip2/Lobby?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

**gogeoip2** is a Go-based IP lookup tool which provides geographic metadata for an IPv4 address in it's simplest form. It can _also_ be extended (with a [premium database](https://www.maxmind.com/en/geoip2-isp-database)) to provide ISP information (i.e. ASN and ISP name)


This data allows one to make informed decisions about who owns IP addresses/network blocks. This tool was used at [HotelTonight](https://www.hoteltonight.com) to make decisions about automatically blocking traffic using a scoring algorithm based on network reputation.

A demo service is currently hosted (and rate-limitted) at http://geo.srv.im — hosted on Docker by [Hyper](https://www.hyper.sh)!

#### Example Output

```shell
~  % curl -s "http://geo.srv.im/ip/173.247.196.18?lat=35.0&lng=-120.0" | jq .
{
  "ip": "173.247.196.18",
  "asnum": 19165,
  "city": "San Francisco",
  "country": "United States",
  "country-iso": "US",
  "continent": "NA",
  "lat": 37.7758,
  "lng": -122.4128,
  "time-zone": "America/Los_Angeles",
  "distance": 376.6809639080358,
  "hostnames": [
    "x.196.247.173.web-pass.com."
  ],
  "asorg": "Webpass Inc.",
  "isp": "Webpass",
  "org": "Webpass",
  "cidr-report": "http://www.cidr-report.org/cgi-bin/as-report?as=AS19165&view=2.0",
  "cleantalk-report": "https://cleantalk.org/blacklists/AS19165",
  "ipinfo-report": "http://ipinfo.io/AS19165",
  "peering-db-report": "https://beta.peeringdb.com/api/asn/19165"
}
```

#### Running Locally

1. Download the MaxMind [GeoLite2 City database](http://geolite.maxmind.com/download/geoip/database/GeoLite2-City.mmdb.gz) and place it in the root directory.
2. (optional) Purchase and download [the GeoIP2 ISP database](https://www.maxmind.com/en/geoip2-isp-database) and put it in the root directory.
3. Grab the build dependencies using `go get -v ./...` from the root directory.
4. Build the package using `go build`
5. Run the server using `./gogeoip2 server`


#### Legal Notices

This product includes GeoLite2 data created by MaxMind, available from http://www.maxmind.com.
