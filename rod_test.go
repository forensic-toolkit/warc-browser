package warcbrowser_test

import (
	"os"
	"io"
	"fmt"
	"bufio"
	"testing"
	"net/http"
    "github.com/gorilla/mux"
	"net/http/httptest"
    "crypto/md5"

	"github.com/stretchr/testify/require"
	"github.com/nlnwa/gowarc"
	"warcbrowser"

)

func md5Content(t *testing.T, path string) string {

	hasher := md5.New()

	f, err := os.Open(path)
	require.Nil(t, err)
	defer f.Close()

	_, err = io.Copy(hasher, f);
	require.Nil(t, err)

	return fmt.Sprintf("%x", hasher.Sum(nil))
}


func retrieveRecordedContent(t *testing.T, warcpath, outpath, mime string)  {
	t.Log("Reading Warc file:" + warcpath)

	reader, err := gowarc.NewWarcFileReader(warcpath, 0, gowarc.WithStrictValidation())
	require.Nil(t, err )
	
	recorded, err := os.Create(outpath)
	require.Nil(t, err)
	defer recorded.Close()

	index := 0
	for {
		record, _, _, err := reader.Next()
		if err == io.EOF {
			break
		}
		require.Nil(t, err)

		t.Log( record.Type().String() )
		t.Log( record.Version() )
		
		switch index {
		// Warc header
		case 0:
			require.True(t, record.Type().String() == "warcinfo" )
		// Warc Content					
		default:
			require.True(t, record.Type().String() == "response" )
			t.Log("WARC.Record.Content-Type:" + record.WarcHeader().Get(gowarc.ContentType))
			require.True(t, record.WarcHeader().Get(gowarc.ContentType) == mime + ";msgtype=response" )

			r, err := record.Block().RawBytes()
			require.Nil(t, err)
			
			resp, err := http.ReadResponse(bufio.NewReader(r), nil)
			// require.True(t, resp.StatusCode == 200)
			t.Log("HTTP.Responce.Content-Type:" + resp.Header.Get("Content-Type"))
			require.True(t, resp.Header.Get("Content-Type") == mime)
			require.Nil(t, err)
			chunc, err := io.ReadAll(resp.Body)
			require.Nil(t, err)

			n, err := recorded.Write(chunc)
			require.Nil(t, err)
			t.Log(fmt.Sprintf("-> wrote %d bytes\n", n))

		}
		index = index + 1
	}
}

func createWarcRecord(t *testing.T, dir, name, url string) error {
	
	br, err := warcbrowser.LaunchRodBrowser(
		gowarc.NewWarcFileWriter(
			gowarc.WithCompression(false),
			gowarc.WithMaxConcurrentWriters(1),
			gowarc.WithFileNameGenerator(&gowarc.PatternNameGenerator{
				Directory: dir, 
				Pattern: name + ".warc",
			}),
			gowarc.WithWarcInfoFunc(func(rbld gowarc.WarcRecordBuilder) error {
				rbld.AddWarcHeader(gowarc.WarcRecordID, "<urn:uuid:667721cd-7619-485a-80c0-2d486b3dedf2>")
				return nil
			}),
		),
		false, // attach
		"",
		false, // headless
	)
	require.Nil(t, err)

	return br.ArchiveUrl(url)
}


func TestFiles(t *testing.T) {

	dir := t.TempDir()

	r := mux.NewRouter()
	staticFileDirectory := http.Dir("./assets/")
	staticFileHandler   := http.FileServer(staticFileDirectory)
	r.PathPrefix("/").Handler(staticFileHandler).Methods("GET")
	s := httptest.NewServer(r)
	defer s.Close()


	t.Run("setup", func(t *testing.T) {

		t.Run("test image", func(t *testing.T) {
			err := createWarcRecord(t, dir, "example.jpeg", s.URL + "/example.jpeg")
			require.Nil(t, err)
			retrieveRecordedContent(t,  dir + "/example.jpeg.warc", dir + "/example.jpeg", "image/jpeg")
		})

		t.Run("test video", func(t *testing.T) {
			err := createWarcRecord(t, dir, "example.mp4", s.URL + "/example.mp4")
			require.Nil(t, err)
			retrieveRecordedContent(t,  dir + "/example.mp4.warc", dir + "/example.mp4", "video/mp4")
		})

		t.Run("test pdf", func(t *testing.T) {
			err := createWarcRecord(t, dir, "example.pdf", s.URL + "/example.pdf")
			require.Nil(t, err)
			retrieveRecordedContent(t,  dir + "/example.pdf.warc", dir + "/example.pdf", "application/pdf")
		})


		t.Run("test html", func(t *testing.T) {
			err := createWarcRecord(t, dir, "example.html", s.URL + "/example.html")
			require.Nil(t, err)
			retrieveRecordedContent(t,  dir + "/example.html.warc", dir + "/example.html", "text/html; charset=utf-8")
		})

	})

}
