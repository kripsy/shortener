package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"text/template"
	"time"

	//nolint:gosec
	_ "net/http/pprof"

	//nolint:depguard
	"github.com/kripsy/shortener/internal/app/application"
	"github.com/kripsy/shortener/internal/app/utils"
)

var (
	//nolint:gochecknoglobals
	buildVersion string
	//nolint:gochecknoglobals
	buildDate string
	//nolint:gochecknoglobals
	buildCommit string
)

type BuildData struct {
	BuildVersion string
	BuildDate    string
	BuildCommit  string
}

const Template = `	Build version: {{if .BuildVersion}} {{.BuildVersion}} {{else}} N/A {{end}}
	Build version: {{if .BuildDate}} {{.BuildDate}} {{else}} N/A {{end}}
	Build version: {{if .BuildCommit}} {{.BuildCommit}} {{else}} N/A {{end}}
`

func main() {
	const idleTimeoutSeconds = 30
	const readHeaderTimeoutSeconds = 2
	ctx := context.Background()

	application, err := application.NewApp(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)

		return
	}

	l := application.GetAppLogger()

	d := &BuildData{
		BuildVersion: buildVersion,
		BuildDate:    buildDate,
		BuildCommit:  buildCommit,
	}

	t := template.Must(template.New("buildTags").Parse(Template))

	err = t.Execute(os.Stdout, *d)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)

		return
	}
	defer func() { // flushes buffer, if any
		if err = l.Sync(); err != nil {
			fmt.Printf("error: %v\n", err)

			return
		}
	}()

	defer application.GetAppRepo().Close() // close repo

	fmt.Printf("SERVER_ADDRESS: %s\n", application.GetAppConfig().URLServer)
	fmt.Printf("BASE_URL: %s\n", application.GetAppConfig().URLPrefixRepo)

	connsClosed := make(chan struct{})
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	srv := &http.Server{
		Addr:              application.GetAppConfig().URLServer,
		ReadTimeout:       time.Second,
		WriteTimeout:      time.Second,
		IdleTimeout:       idleTimeoutSeconds * time.Second,
		ReadHeaderTimeout: readHeaderTimeoutSeconds * time.Second,
		Handler:           application.GetAppServer().Router,
	}

	go func() {
		<-sigint
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "error shotdown server: %v\n", err)
		}
		close(connsClosed)
	}()

	if application.GetAppConfig().EnableHTTPS != "" {
		l.Debug("creating cert")
		err = utils.CreateCertificate()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)

			return
		}
		l.Debug("cert has been created")
		err = srv.ListenAndServeTLS(utils.ServerCertPath, utils.PrivateKeyPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)

			return
		}

		return
	}
	err = srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)

		return
	}
	<-connsClosed
	l.Debug("Server Shutdown successfully")
}
