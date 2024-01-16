// __TODO__
// ADD CDX Indexing for warc records
// Add api endpoint for searching warc records
// Add api endpoint to create warc records
package web

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"fmt"
	"strings"
	"log"

	"github.com/gorilla/mux"
	// "github.com/jasonwbarnett/fileserver"
)

//go:embed dist
var app embed.FS

// https://github.com/gorilla/mux#serving-single-page-applications
type archiveHandler struct {
	staticPath string	
}

func (h archiveHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Access-Control-Allow-Origin
	// development
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	
	// Join internally call path.Clean to prevent directory traversal
	path := filepath.Join(h.staticPath, strings.TrimPrefix(r.URL.Path,"/archives/") )
	
	// check whether a file exists or is a directory at the given path
	fi, err := os.Stat(path)
	
	if os.IsNotExist(err) || fi.IsDir() {
		log.Printf("Listing [%s] %s ", h.staticPath, path)

		// file does not exist or path is a directory, so we list it and serve it as
		// json list to be used by vue app
		files, err := os.ReadDir(path)
		if err != nil {
			w.Write([]byte( fmt.Sprintf(`{"error":"%s","records":[]}`, err.Error() ) ))
			return
		}
		w.Write([]byte(`{"error":null, "records":[`))
		for i, fl := range files {
			w.Write([]byte(
				fmt.Sprintf(`{"label":"%s","branch":%v,"url":"/%s","path":"%s"}`,
								fl.Name(),
								fl.IsDir(),
								filepath.Join(path, fl.Name()), 
								filepath.Join(path, fl.Name()), )))
			if i < len(files) - 1 {
				w.Write([]byte(","))
			}
		}
		w.Write([]byte(`]}`))
		return
	}

	if err != nil {
		w.Write([]byte( fmt.Sprintf(`{"error":"%s","records":[]}`, err.Error() ) ))
        return
	}

	// otherwise, serve warc file
	http.ServeFile(w, r, path)
}

func App(archivesdir string ) http.Handler {
	router := mux.NewRouter()
	// Serve Webarchives directory listed as json
	router.PathPrefix("/archives/").Handler(archiveHandler{staticPath: archivesdir,})
	// Serve Vue application embedded in binary
	dist, _ := fs.Sub(app, "dist")
	router.PathPrefix("/").Handler(http.FileServer(http.FS(dist)))
	return router
}
