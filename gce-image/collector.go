package gceimage

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"google.golang.org/api/compute/v1"
)

var (
	promLastUpdateDuration = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gce_image_update_time_seconds",
			Help: "Most recent time to update metrics",
		},
		[]string{"project"},
	)
	promErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gce_image_update_errors_total",
			Help: "Number of update errors",
		},
		[]string{"project", "type"},
	)
	promImageTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gce_image_total",
			Help: "GCE image count",
		},
		[]string{"project"},
	)
	promImageFiles = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gce_image_files_timestamp",
			Help: "GCE image files creation timestamp",
		},
		[]string{"project", "name", "family", "status"},
	)
	promImageBytes = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gce_image_files_bytes",
			Help: "GCE image files archive bytes",
		},
		[]string{"project", "name", "family", "status"},
	)
)

// Update runs the collector query and atomically updates the cached metrics.
func Update(ctx context.Context, client *http.Client, project string) error {
	start := time.Now()
	log.Println("Starting to walk:", start.Format("2006/01/02"))

	total, err := listImages(ctx, client, project)
	promImageTotal.WithLabelValues(project).Set(float64(total))

	log.Println("Total time to Update:", time.Since(start))
	promLastUpdateDuration.WithLabelValues(project).Set(time.Since(start).Seconds())
	return err
}

func listImages(ctx context.Context, client *http.Client, project string) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	computeService, err := compute.New(client)
	if err != nil {
		promErrors.WithLabelValues(project, "compute-new").Inc()
		return 0, fmt.Errorf("Image(%q).compute.New: %v", project, err)
	}

	counter := 0
	req := computeService.Images.List(project)
	if err := req.Pages(ctx, func(page *compute.ImageList) error {
		for _, image := range page.Items {
			layout := "2006-01-02T15:04:05.000Z07:00"
			t, _ := time.Parse(layout, image.CreationTimestamp)
			promImageFiles.WithLabelValues(project, image.Name, image.Family, image.Status).Set(float64(t.Unix()))
			promImageBytes.WithLabelValues(project, image.Name, image.Family, image.Status).Set(float64(image.ArchiveSizeBytes))
			counter++
		}
		return nil
	}); err != nil {
		promErrors.WithLabelValues(project, "compute-list").Inc()
		return 0, fmt.Errorf("Image(%q).compute.New: %v", project, err)
	}

	return counter, nil
}
