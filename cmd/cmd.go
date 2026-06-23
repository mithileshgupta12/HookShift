package cmd

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/mithileshgupta12/hook-shift/api"
	"github.com/mithileshgupta12/hook-shift/queue"
	"github.com/mithileshgupta12/hook-shift/worker"
)

func Execute() {
	portAddr := flag.Int("port", 9000, "port defines the port on which the api runs")
	workersAddr := flag.Int("workers", 5, "workers defines the number of worker nodes for processing the jobs")
	flag.Parse()

	inMemoryQueue := queue.NewInMemoryQueue()

	worker.StartPool(inMemoryQueue, *workersAddr)

	mux := http.NewServeMux()

	api.InitializeRoutes(mux, inMemoryQueue)

	log.Printf("listening on port %d", *portAddr)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *portAddr), mux))
}
