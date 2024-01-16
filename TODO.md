## Wishlist

- [ ] Add CDX indexing of warc files ( golang implementations [1][nlnwa/gowarcserver/index], [2][datatogether/cdxj]).
- [ ] Add the ability to search using cdx
- [ ] Add the ability to capture url from UI
- [ ] spawning (headless or not) browser and keep it open.
- [ ] Impove warc records implementation following standards
    - Add request Record 
    - Add WARC-Concurrent-To to responces
    - Add WARC-IP-Address fields
    - Add DNS Record
- [ ] been able to control multiple browsers with different profiles [rod manager]
- [ ] Capture list of urls from file or stdin
- [ ] Provive http proxy for other tools to use
- [ ] Add configuration to tweak browser behaviour
   - default or custom browser profiles
   - customize devices/evading fingerprint
   - use proxies
   - run custom scripts for specified urls (url regex)
   - consider different implementation [rod], [playwright], [chromedp]

[replayweb.page]: https://replayweb.page/docs/embedding
[rod manager]: https://pkg.go.dev/github.com/go-rod/rod@v0.114.5/lib/launcher/rod-manager
[playwright]: https://pkg.go.dev/github.com/playwright-community/playwright-go
[chromedp]: https://github.com/chromedp/cdproto
[nlnwa/gowarcserver/index]: https://github.com/nlnwa/gowarcserver/blob/8abebb7a0fb825a602bf80f4f55763e0a33bfc0b/index/writers.go#L29
[datatogether/cdxj]: https://github.com/datatogether/cdxj/blob/master/cdxj.go
