package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"

	"github.com/mithileshgupta12/hook-shift/api"
	"github.com/mithileshgupta12/hook-shift/queue"
	"github.com/mithileshgupta12/hook-shift/worker"
)

func Execute() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	portAddr := flag.Int("port", 9000, "port defines the port on which the api runs")
	workersAddr := flag.Int("workers", 5, "workers defines the number of worker nodes for processing the jobs")
	flag.Parse()

	inMemoryQueue := queue.NewInMemoryQueue()

	wg := &sync.WaitGroup{}

	worker.StartPool(inMemoryQueue, *workersAddr, ctx, wg)

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", *portAddr),
		Handler: mux,
	}

	go func() {
		<-ctx.Done()

		server.Shutdown(context.Background())
	}()

	api.InitializeRoutes(mux, inMemoryQueue)

	log.Printf("listening on port %d", *portAddr)
	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

	wg.Wait()
}
