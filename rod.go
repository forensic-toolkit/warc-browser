
package warcbrowser

import (
	"fmt"
	"time"
	"net/http"
	"errors"
	"strings"
	"github.com/nlnwa/gowarc"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"	
	"github.com/go-rod/rod/lib/proto"
    "github.com/rs/zerolog/log"
)


type RodBrowser struct {
	wrc *gowarc.WarcFileWriter
	browser *rod.Browser
	mode  string // Download mode options: MustWaitStable | MustWaitLoad | MustWaitDOMStable | MustWaitNavigation "
}

// LaunchRodBrowser returns *rod.Browser wrapped as warcbrowser.Browser
func LaunchRodBrowser(w *gowarc.WarcFileWriter, attach bool, attachTo string, headless bool) ( Browser, error ) {
	browser := rod.New().NoDefaultDevice()
	if attach {
		browser = browser.ControlURL(launcher.MustResolveURL(attachTo)).MustConnect()	
	} else {
		l := launcher.New().Headless(headless) // .Set("user-data-dir", "path").
		path, found := launcher.LookPath()
		if found {
			l = l.Bin(path)
		}
		browser = browser.ControlURL(l.MustLaunch()).MustConnect()
	}
	// __TODO__: Add flag for cert errors
	//
	// browser = browser.MustIgnoreCertErrors(true)
	// https://github.com/unstppbl/gowap/blob/master/pkg/scraper/scraper_rod.go
	return &RodBrowser{wrc: w, browser: browser, mode: "MustWaitStable",}, nil
}

// listPages will list all tabs excluding extentions or other wild content 
func (r *RodBrowser) listPages(pattern string) []*rod.Page {
	ret := []*rod.Page{}
	pages, err := r.browser.Pages()
	if err == nil {
		for _, page := range pages {
			info, err := page.Info()
			if err != nil {
				continue
			}
			// Skip unwanted types
			switch info.Type {
			case proto.TargetTargetInfoTypePage:
				// chrome://
				// chrome-extention://
				if strings.HasPrefix(info.URL, "chrome") {
					continue
				}	
			default:
				continue
			}
			ret = append(ret, page)
		}
	}
	return ret
}

func (r *RodBrowser) ListTabs(pattern string) []Tab {
	ret := []Tab{}
	for i, page := range r.listPages(pattern) {
		info, err := page.Info()
		if err != nil {
			continue
		}
		ret = append(ret, Tab{Id: i, Url: info.URL, Title: info.Title })
	}
	return ret
}

func (r *RodBrowser) ArchiveUrl(url string) error {
	page := r.browser.MustPage("")
	defer page.MustClose()
	return ArchivePage(r.wrc, page, url, false, r.mode)
}

func (r *RodBrowser) ArchiveTab(tab int) error {
	pages := r.listPages("")
	for i, page := range pages {
		if i == tab {

			// reload tab and archive it
			return ArchivePage(r.wrc, page, "", true, r.mode)		
			// or create a new page and achive it
			// info, err := page.Info()
			// if err != nil {
			// 	continue
			// }
			// _page := r.browser.MustPage("")
			// defer _page.MustClose()
			// return ArchivePage(r.wrc, _page, page.MustInfo().URL, false,  r.mode)
		}
	}
	return fmt.Errorf("Invalid Tab id provided. Number of tabs: %d", len(pages))
}

// warcWriterRoutine is a helpfull goroutine that writes warc records sequencially from given channel 
func warcWriterRoutine(wrc *gowarc.WarcFileWriter, r chan gowarc.WarcRecord){
	defer wrc.Close()
	for {
		record := <- r
		if record == nil {
			break
		}
		fl := wrc.Write(record)
		log.Debug().Str("Warc.RecordId", record.RecordId()).Str("Output", fl[0].FileName).Msg("Written successfully")
	}
}

// ArchivePage will hijack all requests on provided page, navigate to givcen url or just reload it, and 
// write all responces to provided gowarc.WarcFileWriter. The function is public, so it may be used from 
// external libraries.
func ArchivePage(warcwriter *gowarc.WarcFileWriter, page *rod.Page, url string, reload bool, mode string) error {
	
	// Set download behaviour to allow downloading files
	//
	// if downloadfiles {
	// 	browser := page.Browser()
	// 	_ = proto.BrowserSetDownloadBehavior{
	// 		Behavior:         proto.BrowserSetDownloadBehaviorBehaviorAllowAndName,
	// 		BrowserContextID: browser.BrowserContextID,
	// 		DownloadPath: filepath.Join(os.TempDir(), "rod", "downloads"),
	// 	}.Call(browser)
	//
	// 	go browser.EachEvent(
	// 		func(e *proto.PageDownloadWillBegin) {
	// 			logger.Printf("PageDownloadWillBegin [%v]\n", e.GUID)
	// 		},
	// 		func(e *proto.PageDownloadProgress) bool {
	// 			completed := "(unknown)"
	// 			if e.TotalBytes != 0 {
	// 				completed = fmt.Sprintf("%0.2f%%", e.ReceivedBytes/e.TotalBytes*100.0)
	// 			}
	// 			logger.Printf("PageDownloadProgress  state: %.10s, completed: %s\n", e.State, completed)
	// 			return e.State == proto.PageDownloadProgressStateCompleted
	// 		})()
	// }

	// Build router
	router := page.HijackRequests() //  browser.HijackRequests()	
	defer router.MustStop()

	records, callback := ArchiveHijack()
	defer close(records)

	// Intercept all page requests 
	router.MustAdd("*", callback)

	// pass requests to ArchiveHijack function that channels each request as gowarc.WarcRecord to
	// a callback, that reads them sequencially and write them to gowarc.WarcFileWriter	 
	go warcWriterRoutine(warcwriter, records)

	// Start intercepting router
	go router.Run()

	// Add event handler to auto accept popup dialogs
	go page.EachEvent(
			func(e *proto.PageJavascriptDialogOpening) {
				_ = proto.PageHandleJavaScriptDialog{
					Accept: true,
				}.Call(page)
			})

	// https://github.com/go-rod/rod/issues/226
	page = page.Timeout(30 * time.Second)

	err := rod.Try(func() {
		if reload {	
			page = page.MustReload()
		} else {
			page = page.MustNavigate(url)
		}
		
		// Wait for page to download all its content
		switch mode {
		case "MustWaitLoad":
			page.MustWaitLoad()
		case "MustWaitDOMStable":
			page.MustWaitDOMStable()	
		case "MustWaitNavigation":
			page.MustWaitNavigation()
		case "MustWaitStable":
			page.MustWaitStable()
		default:
			page.MustWaitStable()
		}
	})

	switch {
	// case err.Reason == "net::ERR_ABORTED": //
	case errors.Is(err, &rod.ErrNavigation{}):
		log.Debug().Str("Url", url).Str("Browser", "rod").Str("mode", mode).Msg(err.Error())
	case err != nil:
		log.Warn().Str("Url", url).Str("Browser", "rod").Str("mode", mode).Msg(err.Error())
		return err	
	}
	log.Info().Str("Url", url).Str("Browser", "rod").Str("mode", mode).Msg("Archived")
	return nil
}

// ArchiveHijack returns a channel and a callback function. The callback is a standard rod 
// hijacking callback that can be attached to a rod router. It filters normal out unwanted 
// content, loads the http responce and builds warc responce records. Warc records are then 
// piped to a channel for further processing.
func ArchiveHijack() (chan gowarc.WarcRecord, func(ctx *rod.Hijack)) {

	records := make(chan gowarc.WarcRecord)
	
	return records, func(ctx *rod.Hijack) {
		
		store := false

		if strings.HasPrefix(ctx.Request.URL().String(), "chrome") {
			ctx.Response.Fail(proto.NetworkErrorReasonBlockedByClient)
		} 
		// else {
		// 	ctx.ContinueRequest(&proto.FetchContinueRequest{})	
		// }

		switch ctx.Request.Type() {
		// Blocked
		case proto.NetworkResourceTypeWebSocket, // "WebSocket"
			proto.NetworkResourceTypeSignedExchange, // "SignedExchange"
			proto.NetworkResourceTypeCSPViolationReport, // "CSPViolationReport"
			// proto.NetworkResourceTypeXHR, // "XHR"
			proto.NetworkResourceTypeOther: // "Other"
			ctx.ContinueRequest(&proto.FetchContinueRequest{})
		// Allowed
		// case proto.NetworkResourceTypeDocument: // "Document"
		// case proto.NetworkResourceTypeStylesheet: // "Stylesheet"
		// case proto.NetworkResourceTypeImage: // "Image"
		// case proto.NetworkResourceTypeMedia: // "Media"
		// case proto.NetworkResourceTypeFont: // "Font"
		// case proto.NetworkResourceTypeScript: // "Script"
		// case proto.NetworkResourceTypeTextTrack: // "TextTrack"
		// case proto.NetworkResourceTypeFetch: // "Fetch"
		// case proto.NetworkResourceTypePrefetch: // "Prefetch"		
		// case proto.NetworkResourceTypeEventSource: // "EventSource"		
		// case proto.NetworkResourceTypeManifest: // "Manifest"
		// case proto.NetworkResourceTypePreflight: // "Preflight"
		// case proto.NetworkResourceTypePing: // "Ping"
		default:
			store = true
		}
	
		if ! store {
			return
		}
		
		if err := rod.Try(func() {
			ctx.MustLoadResponse()
		}); err != nil {
			log.Debug().
				Str("Browser", "rod").
				Str("Url", ctx.Request.URL().String()).
				Str("Stage", "MustLoadResponse").
				Msg(err.Error())
			return 
		}

		// Create warc record 
		builder := gowarc.NewRecordBuilder(
			gowarc.Response,
			gowarc.WithAddMissingRecordId(true),
			gowarc.WithAddMissingDigest(true),
			gowarc.WithAddMissingContentLength(true),
			gowarc.WithFixContentLength(true),
			gowarc.WithFixDigest(true),
			gowarc.WithFixSyntaxErrors(true),
			gowarc.WithFixWarcFieldsBlockErrors(true),
			gowarc.WithStrictValidation())

		// Add Warc Record headers
		builder.AddWarcHeader(gowarc.WarcTargetURI, ctx.Request.URL().String())
		builder.AddWarcHeader(gowarc.ContentType,   ctx.Response.Headers().Get("Content-Type") + ";msgtype=response" )
		builder.AddWarcHeaderTime(gowarc.WarcDate,  time.Now())

		// Write Status Line
		statusCode := ctx.Response.Payload().ResponseCode
		statusText := http.StatusText(int(statusCode))
		builder.WriteString(fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, statusText))

		// Write Headers
		for header := range ctx.Response.Headers() {
			builder.WriteString(fmt.Sprintf("%s: %s\r\n", header, ctx.Response.Headers().Get(header) ))
		}
		builder.WriteString("\r\n")

		// Write body
		builder.WriteString(ctx.Response.Body())

		// Build warc record
		wr, _, err := builder.Build()
		if err != nil {
			log.Warn().Str("Browser", "rod").Str("Url", ctx.Request.URL().String()).Msg("Failed to build record")
			log.Debug().Str("Browser", "rod").Str("Url", ctx.Request.URL().String()).Msg(err.Error())
		} else {
			records <- wr
			log.Info().Str("Browser", "rod").Str("Url", ctx.Request.URL().String()).Msg("Downloaded")
		}
	}
}
