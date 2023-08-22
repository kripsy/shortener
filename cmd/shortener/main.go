package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	_ "net/http/pprof"

	"github.com/kripsy/shortener/internal/app/application"
)

func main() {
	ctx := context.Background()

	application, err := application.NewApp(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	defer application.GetAppLogger().Sync() // flushes buffer, if any
	defer application.GetAppRepo().Close()  // close repo

	fmt.Printf("SERVER_ADDRESS: %s\n", application.GetAppConfig().URLServer)
	fmt.Printf("BASE_URL: %s\n", application.GetAppConfig().URLPrefixRepo)

	err = http.ListenAndServe(application.GetAppConfig().URLServer, application.GetAppServer().Router)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
