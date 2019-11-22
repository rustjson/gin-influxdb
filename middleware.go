package gininflux

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	client "github.com/influxdata/influxdb1-client/v2"
)

type GinInflux struct {
	bp             client.BatchPoints
	database       string
	conn           client.Client
	pointName      string
	writeThreshold int
}

func New(addr, database, pointName string, writeThreshold int) GinInflux {

	conn, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: addr,
	})
	if err != nil {
		panic(err)
	}
	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  database,
		Precision: "s",
	})
	if err != nil {
		panic(err)
	}

	return GinInflux{
		bp:             bp,
		conn:           conn,
		database:       database,
		pointName:      pointName,
		writeThreshold: writeThreshold,
	}
}

func (g *GinInflux) write(bp *client.Point) {
	g.bp.AddPoint(bp)
	if len(g.bp.Points()) >= g.writeThreshold {
		err := g.conn.Write(g.bp)
		if err != nil {
			fmt.Errorf("Write to InfluxDB error, err=%v", err)
		} else {
			bp, err := client.NewBatchPoints(client.BatchPointsConfig{
				Database:  g.database,
				Precision: "ms",
			})

			if err != nil {
				fmt.Errorf("Create batch points error, err=%v", err)
			}
			g.bp = bp
		}
	}
}

func (g *GinInflux) HandlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		status := strconv.Itoa(c.Writer.Status())

		c.Next()
		elapsed := float64(time.Since(start)) / float64(time.Second)

		go func() {
			fields := map[string]interface{}{
				"request_uri": c.Request.RequestURI,
				"elapsed":     elapsed,
			}
			tags := map[string]string{
				"method":      c.Request.Method,
				"path":        c.FullPath(),
				"status":      status,
				"request_uri": c.Request.RequestURI,
			}
			pt, err := client.NewPoint(g.pointName, tags, fields, time.Now())
			if err != nil {
				log.Fatal(err)
			}
			g.write(pt)
		}()
	}
}
