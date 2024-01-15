//
// # WarcBrowser
//
// WarcBrowser is an experinmental package with focus on archiving web content as warc files. 
//
// The area of interest is to provide the ability to capture web content on personal computers.
// Web content is considered anything renders on a regular browsers, like webpages, pdfs, images,
// videos. The content may be public or private, where public considered any website's content
// being available without any user login, and private otherwise (e.g social media content, private 
// forums, ...). The web content, will be recorded as ".warc" files on disk. Warc files could them 
// be used to browse the content offline, and optionally publish it using [replayweb.page]. 
//
// 
// Features
//    - capture tabs (by id, url regex) on running browsers
//    - capture url from cli
//
// Wishlist
//    - Provide ui for simple archiving and browing archives (bind on localhost, or other network address or UNIX socket)
//    - Add CDX indexing of warc files. Implementations [nlnwa/gowarcserver/index], [datatogether/cdxj].
//    - Provide http api to control browser 
//    - any assets (js, ...) should be embeed into the binary
//    - be able to control existing browser
//    - spawning (headless or not) browser and keep it open.
//    - been able to control multiple browsers with different profiles [rod manager]
//    - Capture multiple urls (from file, stdin, cli args)
//    - Provive http proxy for other tools to use
//
// Customize from cli flags/configuration to tweak browser behaviour
//    - default or custom browser profiles
//    - customize devices/evading fingerprint
//	  - use proxies
//    - run custom scripts for specified urls (url regex)
//    - consider different implementation [rod], [playwright], [chromedp]
//
// [replayweb.page]: https://replayweb.page/docs/embedding
// [rod manager]: https://pkg.go.dev/github.com/go-rod/rod@v0.114.5/lib/launcher/rod-manager
// [playwright]: https://pkg.go.dev/github.com/playwright-community/playwright-go
// [chromedp]: https://github.com/chromedp/cdproto
// [nlnwa/gowarcserver/index]: https://github.com/nlnwa/gowarcserver/blob/8abebb7a0fb825a602bf80f4f55763e0a33bfc0b/index/writers.go#L29
// [datatogether/cdxj]: https://github.com/datatogether/cdxj/blob/master/cdxj.go
package warcbrowser