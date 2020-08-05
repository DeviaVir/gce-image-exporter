package main

import (
	"flag"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/m-lab/go/flagx"
	"github.com/m-lab/go/prometheusx"
	"github.com/m-lab/go/rtx"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"

	gceimage "github.com/DeviaVir/gce-image-exporter/gce-image"
)

var (
	projects     flagx.StringArray
	collectTimes flagx.DurationArray
)

func init() {
	flag.Var(&projects, "project", "<project-id>")
	flag.Var(&collectTimes, "time", "Run collections at given interval <600s>.")
	log.SetFlags(log.LUTC | log.Lshortfile | log.Ltime | log.Ldate)
}

var (
	mainCtx, mainCancel = context.WithCancel(context.Background())
)

// updateForever runs the gceimage.Update on the given bucket at the given collect time every day.
func updateForever(ctx context.Context, wg *sync.WaitGroup, client *http.Client, bucket string, collect time.Duration) {
	defer wg.Done()

	gceimage.Update(mainCtx, client, bucket)

	for {
		select {
		case <-mainCtx.Done():
			return
		case <-time.After(collect):
			gceimage.Update(mainCtx, client, bucket)
		}
	}
}

var logFatal = log.Fatal

func main() {
	flag.Parse()
	rtx.Must(flagx.ArgsFromEnv(flag.CommandLine), "Failed to parse args")

	if len(projects) != len(collectTimes) {
		logFatal("Must provide same number of projects as collection times.")
	}

	srv := prometheusx.MustServeMetrics()
	defer srv.Close()

	client, err := google.DefaultClient(mainCtx, compute.CloudPlatformScope)
	rtx.Must(err, "Failed to create client")

	wg := sync.WaitGroup{}
	for i, t := range collectTimes {
		wg.Add(1)
		go updateForever(mainCtx, &wg, client, projects[i], t)
	}
	wg.Wait()
}
