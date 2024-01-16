package main

import (
	"os"
	"fmt"
	"net/http"
	"github.com/google/uuid"
	"github.com/nlnwa/gowarc"
	"github.com/urfave/cli/v2"
	"github.com/gorilla/handlers"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"warcbrowser"
	"warcbrowser/web"
)

var warcwriter *gowarc.WarcFileWriter
var browser    warcbrowser.Browser

func ensureDirectory(dir string) error {
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(dir, os.ModePerm)
		}
	}
	return err
}


func main() {
    app := &cli.App{
		Name:  "warc-browser",
		Usage: "preparing for apocalypse",
        EnableBashCompletion: true,
		Flags: []cli.Flag{
            &cli.StringFlag{Name:"output-dir",Aliases:[]string{"d"},Value:"archives", Usage: "directory to store warc archives", },
			&cli.StringFlag{Name:"output",    Aliases:[]string{"o"},Value:"<autogenerated>", Usage: "warc output filename, can be a pattern",},
		},
        Commands: []*cli.Command{
            {
                Name:    "warc",
                Usage:   "Work with stored .warc files",
				Action: func(cCtx *cli.Context) error {
                    return fmt.Errorf("Not yet implemented\n")
                },
            },
			{
                Name: "browser",
                Aliases: []string{"br"},
                Usage: "work with browsers",
				Flags: []cli.Flag{
					&cli.StringFlag{Name:"type", Value:"rod",Usage:"type of browser wrapper to use [rod|playwright]",},
					&cli.StringFlag{Name:"attach-to", Value:"http://localhost:9222",Usage:"address of running browser to connect to",},
					&cli.BoolFlag{Name:"attach", Aliases:[]string{"a"},Value:false,Usage:"attach to existing browser",},
					&cli.BoolFlag{Name:"headless", Value:false,Usage:"Run browser in headless mode. Only for new browsers",},
				},
				// Before hook to setup gowarc.NewWarcFileWriter and warcbrowser.Browser 
				// interfaces for `browser` command and all its subcommands.
				Before: func(ctx *cli.Context) error {

                    err := ensureDirectory(ctx.String("output-dir")) 
					if err != nil {
						return err
					}
					generator := &gowarc.PatternNameGenerator{Directory: ctx.String("output-dir"),}
					if ctx.String("output") != "<autogenerated>" {
						generator.Pattern = ctx.String("output")
					}
					warcwriter = gowarc.NewWarcFileWriter(
									gowarc.WithCompression(false),
									gowarc.WithMaxConcurrentWriters(1),
									gowarc.WithWarcInfoFunc(func(b gowarc.WarcRecordBuilder) error {
										b.AddWarcHeader(gowarc.WarcRecordID, fmt.Sprintf("<%s>", uuid.New().URN()))
										return nil
									}),
									gowarc.WithFileNameGenerator(generator))
                    
					switch ctx.String("type") {
					case "rod":
						browser, err = warcbrowser.LaunchRodBrowser(
											warcwriter,
											ctx.Bool("attach"),
											ctx.String("attach-to"),
											ctx.Bool("headless"),)
						return err
					}
					return fmt.Errorf("Unknown browser") 
                },
				Action: func(ctx *cli.Context) error {
					if ctx.Bool("attach") {
						for _, t := range browser.ListTabs("") {
							fmt.Fprintln(ctx.App.Writer, t.String())
						}
					} else {
						cli.ShowSubcommandHelp(ctx)
					}
					return nil
				},
                Subcommands: []*cli.Command{

					{
						Name: "list-tabs",
						Aliases: []string{"tabs"},
						Usage:   "List available tabs",
						Action: func(ctx *cli.Context) error {
							if ! ctx.Bool("attach") {
								return fmt.Errorf("Listing tabs is only available for attached browsers. Specify --attach/-a flag")
							}
							for _, t := range browser.ListTabs("") {
								fmt.Fprintln(ctx.App.Writer, t.String())
							}
							return nil
						},	
					},

					{
						Name: "archive",
						// Aliases: []string{"a"},
						Usage:   "Capture web content and store it in disk",
						Flags: []cli.Flag{
							&cli.IntFlag{
								Name:  "tab",
								Aliases: []string{"t"},
								Value: -1,
								Usage: "specify tab to archive",
							},
							&cli.StringFlag{
								Name:  "url",
								Value: "",
								Usage: "specify url to archive",
							},
						},
						Action: func(ctx *cli.Context) error {

							if ctx.String("url") != "" {
								return browser.ArchiveUrl(ctx.String("url"))
							} 
							
							if ctx.Int("tab") >= 0 {
								if ! ctx.Bool("attach") {
									return fmt.Errorf("Archiving tabs is only available for attached browsers. Specify --attach/-a flag")
								} else {
									return browser.ArchiveTab(ctx.Int("tab"))
								}
							}

							return fmt.Errorf("Please specify either --url/--tab flag")
						},
					},
				},
			},
			{
                Name:    "daemon",
                Usage:   "Start daemon",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "address",
						Value: ":8080",
						Usage: "specify address to bind to",
					},
				},
                Action: func(ctx *cli.Context) error {
					fmt.Fprintf(ctx.App.ErrWriter, "Listening on %s\n", ctx.String("address"))
                    return http.ListenAndServe(
								ctx.String("address"),
								handlers.LoggingHandler(os.Stdout, 
									web.App(
										ctx.String("output-dir"),
									)))

                },
            },
        },
    }

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
    if err := app.Run(os.Args); err != nil {
        fmt.Printf(" ! %v", err)
		os.Exit(1)
    }
}