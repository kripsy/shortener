package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"text/template"

	_ "net/http/pprof"

	"github.com/kripsy/shortener/internal/app/application"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
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
	ctx := context.Background()

	application, err := application.NewApp(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return
	}

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
		if err = application.GetAppLogger().Sync(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return
		}
	}()

	defer application.GetAppRepo().Close() // close repo

	fmt.Printf("SERVER_ADDRESS: %s\n", application.GetAppConfig().URLServer)
	fmt.Printf("BASE_URL: %s\n", application.GetAppConfig().URLPrefixRepo)

	err = http.ListenAndServe(application.GetAppConfig().URLServer, application.GetAppServer().Router)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return
	}
}
