module github.com/go-godin/ticket-service

go 1.12

require (
	contrib.go.opencensus.io/exporter/zipkin v0.1.1
	github.com/apache/thrift v0.12.0 // indirect
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/codahale/hdrhistogram v0.0.0-20161010025455-3a0bb77429bd // indirect
	github.com/go-godin/log v0.0.0-20190716173926-b62a2fca0801
	github.com/go-kit/kit v0.9.0
	github.com/golang/protobuf v1.3.2
	github.com/micro/go-micro v1.8.2
	github.com/oklog/run v1.0.0
	github.com/opentracing-contrib/go-observer v0.0.0-20170622124052-a52f23424492 // indirect
	github.com/opentracing/opentracing-go v1.1.0
	github.com/openzipkin-contrib/zipkin-go-opentracing v0.3.5
	github.com/openzipkin/zipkin-go v0.2.0
	github.com/pkg/errors v0.8.1
	github.com/rs/xid v1.2.1
	github.com/uber-go/atomic v1.4.0 // indirect
	github.com/uber/jaeger-client-go v2.16.0+incompatible
	github.com/uber/jaeger-lib v2.0.0+incompatible // indirect
	go.opencensus.io v0.22.0
	go.uber.org/atomic v1.4.0 // indirect
	golang.org/x/net v0.0.0-20190724013045-ca1201d0de80 // indirect
	golang.org/x/sys v0.0.0-20190726091711-fc99dfbffb4e // indirect
	google.golang.org/grpc v1.22.0
)

replace github.com/go-godin/log => ../log
