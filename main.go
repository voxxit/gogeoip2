package main

import (
	"encoding/json"
	"net/http"
	"os"

	"gitlab.com/voxxit/gogeoip2/api"
	"gitlab.com/voxxit/gogeoip2/ingest"

	"github.com/codegangsta/cli"
	"github.com/gin-gonic/gin"
)

func ipLookupHandler(db *api.IPDatabase, c *gin.Context) {
	ip, err := api.LookupIP(db, &api.IPConfig{
		Addr:             c.Param("addr"),
		RequestLatitude:  c.Query("lat"),
		RequestLongitude: c.Query("lng"),
	})

	if err != nil {
		c.Error(err)
	}

	data, err := json.Marshal(&ip)
	if err != nil {
		c.Error(err)
	}

	c.Data(http.StatusOK, "application/json", data)
}

func main() {
	db := api.OpenDatabases()

	defer db.Isp.Close()
	defer db.City.Close()

	app := cli.NewApp()
	app.Name = "gogeoip2"
	app.Usage = "Lookup useful data on any IP"
	app.Commands = []cli.Command{
		{
			Name:    "bulk",
			Aliases: []string{"b"},
			Usage:   "ingest a CSV of IPs, and get data back in CSV",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file",
					Usage: "CSV file to ingest",
				},
			},
			Action: func(c *cli.Context) {
				ingest.NewIngestion(db, c)
			},
		},
		{
			Name:    "server",
			Aliases: []string{"s"},
			Usage:   "Serve up requests for IP data",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "port",
					Value:  "5000",
					Usage:  "port to listen on",
					EnvVar: "PORT",
				},
			},
			Action: func(c *cli.Context) {
				router := gin.New()
				router.Use(gin.Logger())

				router.GET("/ip/:addr", func(g *gin.Context) {
					ipLookupHandler(db, g)
				})

				router.Run(":" + c.String("port"))
			},
		},
	}

	app.Run(os.Args)
}
