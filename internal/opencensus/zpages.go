package opencensus

import (
	"github.com/go-godin/log"
	"net"
	// stdlib
	"net/http"

	// external
	"github.com/oklog/run"
	"go.opencensus.io/zpages"
)

// ZPages handling setup
func ZPages(g *run.Group, logger log.Logger) {
	var (
		bindIP      = "0.0.0.0"
		listener, _ = net.Listen("tcp", bindIP+":3000") // dynamic port assignment
		addr        = listener.Addr().String()
	)


	g.Add(func() error {
		logger.Info( "zpages started", "addr", "http://"+addr)
		return http.Serve(listener, zpages.Handler)
	}, func(error) {
		listener.Close()
	})
}
