package opencensus

import (
	// stdlib
	"io"

	// external
	oczipkin "contrib.go.opencensus.io/exporter/zipkin"
	zipkin "github.com/openzipkin/zipkin-go"
	reporter "github.com/openzipkin/zipkin-go/reporter/http"
	"go.opencensus.io/trace"
	// project
)

const (
	zipkinURL = "http://localhost:9411/api/v2/spans"
)

func setupZipkin(serviceName string) (trace.Exporter, io.Closer) {
	var (
		rep  = reporter.NewReporter(zipkinURL)
		addr = "0.0.0.0"
	)
	localEndpoint, _ := zipkin.NewEndpoint(serviceName, addr)

	return oczipkin.NewExporter(rep, localEndpoint), rep
}
