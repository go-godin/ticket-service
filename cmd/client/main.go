package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/openzipkin/zipkin-go/reporter/http"

	"github.com/go-godin/log"
	"github.com/go-godin/ticket-service/pkg/client"
)

func main() {
	fs := flag.NewFlagSet("ticket-cli", flag.ExitOnError)
	var (
		zipkinURL = fs.String("zipkin-url", "http://localhost:9411/api/v2/spans", "Enable Zipkin v2 tracing (zipkin-go) via a HTTP reporter URL e.g. http://localhost:9411/api/v2/spans")
		method    = fs.String("method", "create", "create")
	)
	_ = fs.Parse(os.Args[1:])

	if *zipkinURL == "" {
		fmt.Println("missing zipkin url")
		os.Exit(1)
	}

	logger := log.NewLoggerFromEnv()
	reporter := http.NewReporter(*zipkinURL)
	defer reporter.Close()
	cl, err := client.New(logger, "localhost:50051", reporter)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *method == "" {
		fmt.Println("missing zipkin url")
		os.Exit(1)
	}
	switch *method {
	case "create":
		if len(fs.Args()) < 2 {
			logger.Error("create requires two parameters, title and description")
			os.Exit(1)
		}
		title := fs.Args()[0]
		desc := fs.Args()[1]
		ticket, err := cl.Create(context.Background(), title, desc)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		logger.Info("ticket created", "ticket.id", ticket.TicketID, "ticket.title", ticket.Title, "ticket.description", ticket.Description, "ticket.status", ticket.Status)
		os.Exit(0)
	default:
		fmt.Println("unknown method or not implemented")
		os.Exit(1)
	}

}
