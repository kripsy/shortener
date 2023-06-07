package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	app, err := newApp()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	defer app.MyLogger.Sync() // flushes buffer, if any

	if app.Config.DatabaseDsn != "" {
		defer app.Server.MyDB.Close()
	}

	fmt.Printf("SERVER_ADDRESS: %s\n", app.Config.URLServer)
	fmt.Printf("BASE_URL: %s\n", app.Config.URLPrefixRepo)
	err = http.ListenAndServe(app.Config.URLServer, app.Server.Router)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
