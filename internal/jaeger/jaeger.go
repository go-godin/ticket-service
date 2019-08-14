package jaeger

import (
	"fmt"
	"io"
	"time"

	jaegercfg "github.com/uber/jaeger-client-go/config"
)

func New(serviceName string, backendAddress string) io.Closer {
	jaegerCfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:                    "const",
			Param:                   1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  backendAddress,
		},
	}
	closer, err := jaegerCfg.InitGlobalTracer(serviceName)
	if err != nil {
		panic(fmt.Sprintf("could not initialize jaeger tracer: %s", err))
	}
	return closer
}
