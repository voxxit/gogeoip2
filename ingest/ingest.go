package ingest

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"gitlab.com/voxxit/gogeoip2/api"

	"github.com/codegangsta/cli"
)

func ingestionWorker(id int, jobs <-chan []string, results chan<- []string, db *api.IPDatabase) {
	for j := range jobs {
		ip, err := api.LookupIP(db, &api.IPConfig{
			Addr:             j[4],
			RequestLatitude:  j[2],
			RequestLongitude: j[3],
		})
		if err != nil {
			log.Fatal(err)
		}

		var hostname string

		if len(ip.Hostnames) == 0 {
			hostname = ""
		} else {
			hostname = strings.Join(ip.Hostnames, ",")
		}

		results <- []string{
			ip.IP,
			fmt.Sprintf("%d", ip.ASNum),
			ip.City,
			ip.Country,
			ip.CountryISO,
			ip.Continent,
			fmt.Sprintf("%f", ip.Latitude),
			fmt.Sprintf("%f", ip.Longitude),
			ip.TimeZone,
			fmt.Sprintf("%f", ip.Distance),
			hostname,
			ip.ASOrg,
			ip.ISP,
			ip.Org,
			ip.CIDRReport,
			ip.CleanTalkReport,
			ip.IPInfoReport,
			ip.PeeringDBReport,
		}
	}
}

func NewIngestion(db *api.IPDatabase, c *cli.Context) {
	fileName := c.String("file")
	if fileName == "" {
		log.Fatal("Please specify a file to ingest")
	}

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	rawData, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	jobs := make(chan []string, len(rawData))
	results := make(chan []string, len(rawData))
	writer := csv.NewWriter(os.Stdout)

	for w := 1; w <= 8; w++ {
		go ingestionWorker(w, jobs, results, db)
	}

	for _, elem := range rawData {
		jobs <- elem
	}
	close(jobs)

	for a := 1; a <= len(rawData); a++ {
		writer.Write(<-results)
	}
}
