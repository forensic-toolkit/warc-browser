# warc-browser

a cli toolkit for work with web archives.

warc-browser uses *DevTools* protocol to automate compatible web browsers, captures all content for given wep page (html, css, js, images, videos, pdfs, ...) and stores the results in *[.warc][warc-doc]* file. It came out of need for quickly archiving web pages in a scriptable manner. 
  

## Installation

```bash
make build
./warc-browser --help
```

## Usage

Archive a url running browser in headless mode.

```bash
warc-browser --output-dir /tmp/archives browser --headless archive --url http://example.com
```

Attach to a running browser, list available tabs, then capture specific tab. 

```bash
# Start chromium browser with remote debugging enabled
chromium --remote-debugging-port=9222 --url https://duckduckgo.com/?q=web+archive
# List tabs of chromium
warc-browser browser -a
# Archive first tab
warc-browser browser -a archive -t 0
```

Start a web server serving simple ui, to visualize your collected archives.

```bash
warc-browser ui
```

Open your browser at [localhost:8080](http://localhost:8080).

---

software used 

1. [github.com/go-rod/rod][go-rod/rod] web automation framework for browser automation
2. [github.com/nlnwa/gowarc][nlnwa/gowarc] for composing warc records
3. [github.com/webrecorder/replayweb.page][webrecorder/replayweb.page] for visualizing records in web ui.


```
coverage: 60.8% of statements
```

[nlnwa/gowarc]: https://github.com/nlnwa/gowarc
[go-rod/rod]: https://github.com/go-rod/rod
[warc-doc]: https://iipc.github.io/warc-specifications/specifications/warc-format/warc-1.1/
[webrecorder/replayweb.page]: https://github.com/webrecorder/replayweb.page
